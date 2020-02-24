package handler

import (
    "bytes"
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "github.com/op/go-logging"
    "piot-server/model"
    "piot-server/service"
)

type Adapter struct { }

func (h *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    ctx := r.Context()
    ctx.Value("log").(*logging.Logger).Debugf("Incoming packet")

    // get info of debug mode directly from logger
    if ctx.Value("log").(*logging.Logger).IsEnabledFor(logging.DEBUG) {
        body, err := ioutil.ReadAll(r.Body)
        if err == nil {
            ctx.Value("log").(*logging.Logger).Errorf("Reading request body error: %s", err)
        }
        reqStr := ioutil.NopCloser(bytes.NewBuffer(body))
        ctx.Value("log").(*logging.Logger).Debugf("Request body: %s", reqStr)
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

    ctx.Value("log").(*logging.Logger).Debugf("Packet decoded %v", devicePacket)

    // get instance of piot devices service
    pd := ctx.Value("piotdevices").(*service.PiotDevices)
    if err := pd.ProcessPacket(ctx, devicePacket); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
}
