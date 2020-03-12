package handler

import (
    "bytes"
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "github.com/op/go-logging"
    //"piot-server/model"
    "github.com/mnezerka/go-piot"
    "github.com/mnezerka/go-piot/model"
)

type Adapter struct {
    log *logging.Logger
    piotDevices *piot.PiotDevices
}

func NewAdapter(log *logging.Logger, piotDevices *piot.PiotDevices) *Adapter {
    return &Adapter{log: log, piotDevices: piotDevices}
}

func (h *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    h.log.Debugf("Incoming packet")

    // get info of debug mode directly from logger
    if h.log.IsEnabledFor(logging.DEBUG) {
        body, err := ioutil.ReadAll(r.Body)
        if err == nil {
            h.log.Errorf("Reading request body error: %s", err)
        }
        reqStr := ioutil.NopCloser(bytes.NewBuffer(body))
        h.log.Debugf("Request body: %s", reqStr)
        r.Body = reqStr
    }

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

    h.log.Debugf("Packet decoded %v", devicePacket)

    if err := h.piotDevices.ProcessPacket(devicePacket); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
}
