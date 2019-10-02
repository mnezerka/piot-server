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

    pd := ctx.Value("piotdevices").(*service.PiotDevices)
    if err := pd.ProcessPacket(ctx, devicePacket); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
}
