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

    // handle short notation of attributes (assign short to long attributes)
    if len(packet.DeviceShort) > 0 { packet.Device = packet.DeviceShort }
    if len(packet.ReadingsShort) > 0 { packet.Readings = packet.ReadingsShort }

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

    // name of the device cannot be empty
    if packet.Device == "" {
        return errors.New("Device name cannot be empty")
    }

    // store device name to cache together with date it was seen
    p.cache[packet.Device] = time.Now()

    // get instance of Things service and look for the device (chip),
    // register it if it doesn't exist
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

        // handle short notation of address attribute
        if len(reading.AddressShort) > 0 { reading.Address = reading.AddressShort }

        if reading.Temperature != nil {
            class := model.THING_CLASS_TEMPERATURE
            if err = p.processReading(ctx, class, thing, reading); err != nil {
                ctx.Value("log").(*logging.Logger).Debugf("Failed to process reading data for thing <%s>", thing.Name)
            }
        }

        if reading.Humidity != nil {
            class := model.THING_CLASS_HUMIDITY
            if err = p.processReading(ctx, class, thing, reading); err != nil {
                ctx.Value("log").(*logging.Logger).Debugf("Failed to process humidity reading data for thing <%s>", thing.Name)
            }
        }

        if reading.Pressure != nil {
            class := model.THING_CLASS_PRESSURE
            if err = p.processReading(ctx, class, thing, reading); err != nil {
                ctx.Value("log").(*logging.Logger).Debugf("Failed to process pressure reading data for thing <%s>", thing.Name)
            }
        }
    }

    return nil
}

func (p *PiotDevices) processDevice(ctx context.Context, thing *model.Thing, packet model.PiotDevicePacket) error {

    ctx.Value("log").(*logging.Logger).Debugf("Process PIOT device data: %v", packet)
    mqtt := ctx.Value("mqtt").(IMqtt)

    // dont' push anything if device is disabled
    if !thing.Enabled {
        return nil
    }

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

func (p *PiotDevices) processReading(ctx context.Context, class string, thing *model.Thing, reading model.PiotSensorReading) error {
    ctx.Value("log").(*logging.Logger).Debugf("Process PIOT device reading data of class \"%s\": %v", class, reading)
    mqtt := ctx.Value("mqtt").(IMqtt)

    var address string = reading.Address
    var value string
    var unit string

    // determine address from class
    // this is necessary to have separate things for all sensor measurements
    switch class {
    case model.THING_CLASS_TEMPERATURE:
        address = "T" + address
        unit = "C"
        if reading.Temperature != nil {
            value = strconv.FormatFloat(float64(*reading.Temperature), 'f', -1, 32)
        }
    case model.THING_CLASS_HUMIDITY:
        address = "H" + address
        unit = "%"
        if reading.Humidity!= nil {
            value = strconv.FormatFloat(float64(*reading.Humidity), 'f', -1, 32)
        }
    case model.THING_CLASS_PRESSURE:
        address = "P" + address
        unit = "mPa"
        if reading.Pressure!= nil {
            value = strconv.FormatFloat(float64(*reading.Pressure), 'f', -1, 32)
        }
    }

    // look for thing representing sensor
    things := ctx.Value("things").(*Things)
    sensor_thing, err := things.Find(ctx, address)

    // if thing not found
    if err != nil {

        // register register device
        sensor_thing, err = things.Register(ctx, address, model.THING_TYPE_SENSOR)
        if err != nil {
            return err
        }

        // register topics for measurements (if presetn)
        if things.SetSensorMeasurementTopic(ctx, address, PIOT_MEASUREMENT_TOPIC); err != nil {
            return err
        }

        // set proper device class according to received measurement type
        if err := things.SetSensorClass(ctx, address, class); err != nil {
            return err
        }
    }

    // update parent thing (this can happen any time since sensor can be
    // re-connected to another device
    if (sensor_thing.ParentId != thing.Id) {
        err = things.SetParent(ctx, sensor_thing.Id, thing.Id);
        if err != nil {
            return err
        }
    }

    // if thing is not assigned to org
    if sensor_thing.OrgId == primitive.NilObjectID {
        ctx.Value("log").(*logging.Logger).Debugf("Ignoring processing of data for thing <%s> that is not assigned to any organization", sensor_thing.Name)

        // stop processing here
        return nil
    }

    // dont' push anything if device is disabled
    if !thing.Enabled {
        return nil
    }

    // update avalibility channel
    err = mqtt.PushThingData(ctx, sensor_thing, TOPIC_AVAILABLE, VALUE_YES)
    if err != nil {
        return err
    }

    if value != "" {

        if err := mqtt.PushThingData(ctx, sensor_thing, PIOT_MEASUREMENT_TOPIC, value); err != nil {
            return err
        }
        if err := mqtt.PushThingData(ctx, sensor_thing, fmt.Sprintf("%s/%s", PIOT_MEASUREMENT_TOPIC, TOPIC_UNIT), unit); err != nil {
            return err
        }
    } else {
        ctx.Value("log").(*logging.Logger).Warningf("Processing unkonwn sensor reading data <%v>", reading)
    }

    return nil
}
