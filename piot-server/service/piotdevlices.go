package service

import (
    "context"
    "errors"
    "fmt"
    "strconv"
    "time"
    "github.com/op/go-logging"
    "piot-server/model"
    "piot-server/config"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// topic name used for publishing sensor readings
const PIOT_MEASUREMENT_TOPIC = "value"

type PiotDevices struct {
    cache map[string]time.Time
}

// constructor
func NewPiotDevices() (*PiotDevices) {
    p := PiotDevices{}
    p.cache = make(map[string]time.Time)
    return &p
}

func (p *PiotDevices) ProcessPacket(ctx context.Context, packet model.PiotDevicePacket) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Process PIOT device packet: %v", packet)
    params := ctx.Value("params").(*config.Parameters)

    // DOS Protection
    // allow to process data from this packet only if it didn't come too close
    // to previous packet from the same device (treshold in seconds is defined)
    if lastSeen, ok := p.cache[packet.Device]; ok {
        delta := time.Now().Sub(lastSeen)
        ctx.Value("log").(*logging.Logger).Debugf("Time cache holds entry for device: %s, time diff is %f seconds", packet.Device, delta.Seconds())

        if delta <= params.DOSInterval {
            return errors.New("Exceeded dos protection treshold")
        }
    }

    // store device name to cache together with date it was seen
    p.cache[packet.Device] = time.Now()

    // name of the device cannot be empty
    if packet.Device == "" {
        return errors.New("Device name cannot be empty")
    }

    // look for device (chip) and register it if it doesn't exist
    things := ctx.Value("things").(*Things)

    thing, err := things.Find(ctx, packet.Device)
    if err != nil {
        // register device
        thing, err = things.Register(ctx, packet.Device, model.THING_TYPE_DEVICE)
        if err != nil {
            return err
        }

        // configure availability topic
        if err := things.SetAvailabilityTopic(ctx, packet.Device, "available"); err != nil {
            return err
        }
        if err := things.SetAvailabilityYesNo(ctx, packet.Device, "yes", "no"); err != nil {
            return err
        }
    }

    // if thing is assigned to org
    if thing.OrgId != primitive.NilObjectID {
        // try to push data to mqtt 
        if err = p.processDevice(ctx, thing, packet); err != nil {
            return err
        }
    } else {
        ctx.Value("log").(*logging.Logger).Debugf("Ignoring processing of data for thing <%s> that is not assigned to any organization", thing.Name)
    }

    // look for sensors and register those that doesn't exist
    for _, reading := range packet.Readings {
        // look for (device
        thing, err = things.Find(ctx, reading.Address)
        if err != nil {
            // register register device
            thing, err = things.Register(ctx, reading.Address, model.THING_TYPE_SENSOR)
            if err != nil {
                return err
            }

            // register topics for measurements (if presetn)
            if things.SetSensorMeasurementTopic(ctx, reading.Address, PIOT_MEASUREMENT_TOPIC); err != nil {
                return err
            }

            // set proper device class according to received measurement type
            var class string
            if reading.Temperature != nil {
                class = model.THING_CLASS_TEMPERATURE
            } else if reading.Humidity != nil {
                class = model.THING_CLASS_HUMIDITY
            } else if reading.Pressure != nil {
                class = model.THING_CLASS_HUMIDITY
            } else {
                ctx.Value("log").(*logging.Logger).Warningf("Registering sensor for reading of unknown type <%v>", reading)
            }

            if class != "" {
                if err := things.SetSensorClass(ctx, reading.Address, class); err != nil {
                    return err
                }
            }
        }

        // if thing is assigned to org
        if thing.OrgId != primitive.NilObjectID {
            if err = p.processReading(ctx, thing, reading); err != nil {
                return err
            }
        } else {
            ctx.Value("log").(*logging.Logger).Debugf("Ignoring processing of data for thing <%s> that is not assigned to any organization", thing.Name)
        }
    }

    return nil
}

func (p *PiotDevices) processDevice(ctx context.Context, thing *model.Thing, packet model.PiotDevicePacket) error {

    ctx.Value("log").(*logging.Logger).Debugf("Process PIOT device data: %v", packet)
    mqtt := ctx.Value("mqtt").(IMqtt)

    // update avalibility channel
    err := mqtt.PushThingData(ctx, thing, TOPIC_AVAILABLE, VALUE_YES)
    if err != nil {
        return err
    }

    if packet.Ip != nil {
        err := mqtt.PushThingData(ctx, thing, TOPIC_IP, *packet.Ip)
        if err != nil {
            return err
        }
    }

    if packet.WifiSSID != nil {
        err := mqtt.PushThingData(ctx, thing, TOPIC_WIFI_SSID, *packet.WifiSSID)
        if err != nil {
            return err
        }
    }

    if packet.WifiStrength != nil {
        if err := mqtt.PushThingData(ctx, thing, TOPIC_WIFI_STRENGTH, fmt.Sprintf("%f", *packet.WifiStrength)); err != nil {
            return err
        }
    }

    return nil
}

func (p *PiotDevices) processReading(ctx context.Context, thing *model.Thing, reading model.PiotSensorReading) error {
    ctx.Value("log").(*logging.Logger).Debugf("Process PIOT device reading data: %v", reading)
    mqtt := ctx.Value("mqtt").(IMqtt)

    // update avalibility channel
    err := mqtt.PushThingData(ctx, thing, TOPIC_AVAILABLE, VALUE_YES)
    if err != nil {
        return err
    }

    var value string
    var unit string

    if reading.Temperature != nil {

        value = strconv.FormatFloat(float64(*reading.Temperature), 'f', -1, 32)
        unit = "C"
    }

    if reading.Pressure!= nil {
        value = strconv.FormatFloat(float64(*reading.Pressure), 'f', -1, 32)
        unit = "mPa"
    }

    if reading.Humidity!= nil {
        value = strconv.FormatFloat(float64(*reading.Humidity), 'f', -1, 32)
        unit = "%"
    }

    if value != "" {

        if err := mqtt.PushThingData(ctx, thing, PIOT_MEASUREMENT_TOPIC, value); err != nil {
            return err
        }
        if err := mqtt.PushThingData(ctx, thing, fmt.Sprintf("%s/%s", PIOT_MEASUREMENT_TOPIC, TOPIC_UNIT), unit); err != nil {
            return err
        }
    } else {
        ctx.Value("log").(*logging.Logger).Warningf("Processing unkonw sensor reading data <%v>", reading)
    }

    return nil
}
