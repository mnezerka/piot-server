package handler

import (
    "encoding/json"
    "errors"
    "net/http"
    "github.com/op/go-logging"
    "piot-server/model"
    "piot-server/service"
)

type Adapter struct { }

func (h *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    ctx := r.Context()
    ctx.Value("log").(*logging.Logger).Debugf("Incoming packet")

    // check http method, POST is required
    if r.Method != http.MethodPost {
        WriteErrorResponse(w, errors.New("Only POST method is allowed"), http.StatusMethodNotAllowed)
        return
    }

    // try to decode packet
    var devicePacket model.PiotDevicePacket
    if err := json.NewDecoder(r.Body).Decode(&devicePacket); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("Packet decoded %v", devicePacket)

    // DOS Protection
    // TODO
    // store device to cache together with date
    // allow to process incoming data only if new packet comes later than time windown


    // look for device (chip) and register it if it doesn't exist
    things := ctx.Value("things").(*service.Things)
    device, err := things.Find(ctx, devicePacket.Device)
    if err != nil {
        // register register device
        device, err = things.Register(ctx, devicePacket.Device, model.THING_TYPE_DEVICE)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
    }

    // post data to MQTT if device is enabled
    if (device.Enabled) {
        ctx.Value("log").(*logging.Logger).Debugf("TODO - write data to mqtt for enabled device %v", devicePacket.Device)
    }

    // look for sensors and register those that doesn't exist
    for _, sensor := range devicePacket.Readings {
        // look for (device
        device, err = things.Find(ctx, sensor.Address)
        if err != nil {
            // register register device
            device, err = things.Register(ctx, sensor.Address, model.THING_TYPE_SENSOR)
            if err != nil {
                http.Error(w, err.Error(), 500)
                return
            }
        }

        // post data to MQTT if device is enabled
        if (device.Enabled) {
            ctx.Value("log").(*logging.Logger).Debugf("TODO - write data to mqtt for enabled sensor %v", sensor.Address)
        }
    }
}
