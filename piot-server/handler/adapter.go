package handler

import (
    "encoding/json"
    "errors"
    "net/http"
    "github.com/op/go-logging"
    "piot-server/model"
    "piot-server/service"
)

type Adapter struct {
    things *service.Things
}

func (h *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    ctx := r.Context()
    ctx.Value("log").(*logging.Logger).Debugf("[AD] Incoming packet")

    //db := ctx.Value("db").(*mongo.Database)

    // check http method, POST is required
    if r.Method != http.MethodPost {
        WriteErrorResponse(w, errors.New("Only POST method is allowed"), http.StatusMethodNotAllowed)
        return
    }

    var devicePacket model.PiotDevicePacket

    if err := json.NewDecoder(r.Body).Decode(&devicePacket); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("[AD] Packet decoded %v", devicePacket)

    // look for things (device + sensors)

    // register if necessary

    // post data to MQTT if device (sensor) is registered and enabled
}
