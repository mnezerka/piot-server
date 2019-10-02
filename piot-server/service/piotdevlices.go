package service

import (
    "context"
    "errors"
    "fmt"
    "strconv"
    "time"
    "github.com/op/go-logging"
    "piot-server/model"
)

// the minimal allowed time interval between two packets from
// the same device
//const DOS_TRESHOLD = 30
const DOS_TRESHOLD = 1

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

    // DOS Protection
    // allow to process data from this packet only if it didn't come too close
    // to previous packet from the same device (treshold in seconds is defined)
    if lastSeen, ok := p.cache[packet.Device]; ok {
        delta := time.Now().Sub(lastSeen)
        ctx.Value("log").(*logging.Logger).Debugf("Time cache holds entry for device: %s, time diff is %f seconds", packet.Device, delta.Seconds())

        if delta.Seconds() <= DOS_TRESHOLD {
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
    }

    if err = p.processDevice(ctx, thing, packet); err != nil {
        return err
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
        }

        if err = p.processReading(ctx, thing, reading); err != nil {
            return err
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

    if reading.Temperature != nil {
        if err := mqtt.PushThingData(ctx, thing, TOPIC_TEMP, strconv.FormatFloat(float64(*reading.Temperature), 'f', -1, 32)); err != nil {
            return err
        }
        if err := mqtt.PushThingData(ctx, thing, fmt.Sprintf("%s/%s", TOPIC_TEMP, TOPIC_UNIT), "C"); err != nil {
            return err
        }
    }

    if reading.Pressure!= nil {
        if err := mqtt.PushThingData(ctx, thing, TOPIC_PRESSURE, strconv.FormatFloat(float64(*reading.Pressure), 'f', -1, 32)); err != nil {
            return err
        }
        if err := mqtt.PushThingData(ctx, thing, fmt.Sprintf("%s/%s", TOPIC_PRESSURE, TOPIC_UNIT), "Pa"); err != nil {
            return err
        }
    }

    if reading.Humidity!= nil {
        if err := mqtt.PushThingData(ctx, thing, TOPIC_HUMIDITY, strconv.FormatFloat(float64(*reading.Humidity), 'f', -1, 32)); err != nil {
            return err
        }
        if err := mqtt.PushThingData(ctx, thing, fmt.Sprintf("%s/%s", TOPIC_HUMIDITY, TOPIC_UNIT), "%"); err != nil {
            return err
        }
    }

    return nil
}
