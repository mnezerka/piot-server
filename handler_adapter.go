package main

import (
    "bytes"
    "crypto/aes"
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "github.com/op/go-logging"
    "piot-server/model"
)

type Adapter struct {
    log *logging.Logger
    piotDevices *PiotDevices
    password string
}

func NewAdapter(log *logging.Logger, piotDevices *PiotDevices, password string) *Adapter {
    return &Adapter{log: log, piotDevices: piotDevices, password: password}
}

func (h *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    const size = 16

    h.log.Debugf("Incoming packet")

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        h.log.Errorf("Reading request body error: %s", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // log message content in DEBUG mode
    // get info of debug mode directly from logger
    if h.log.IsEnabledFor(logging.DEBUG) {
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

        h.log.Debugf("Raw data json decode failed, trying to decrypt")

        if len(h.password) != 16 {
            h.log.Error("Failed to decrypt, PIOT password not configured or doesn't have 16 chars")
            WriteErrorResponse(w, errors.New("Missing or wrong encryption configuration"), 500)
            return
        }

        // body shall have length which is multiplication of cipher size (size constant)
        if len(body) % size != 0 {
            h.log.Errorf("Invalid length of body for decryption %d", len(body))
            WriteErrorResponse(w, errors.New("Invalid length of body for decryption"), http.StatusBadRequest)
            return
        }

        // json decode from raw data failed => try to decrypt first
        cipher, _ := aes.NewCipher([]byte(h.password))

        decrypted := make([]byte, len(body))

        // decrypt by blocks
        for bs, be := 0, size; bs < len(body); bs, be = bs + size, be + size {
            cipher.Decrypt(decrypted[bs:be], body[bs:be])
        }

        h.log.Debugf("Decrypted message <%s>", decrypted)

        // try to decode decrypted data
        if err := json.Unmarshal(decrypted, &devicePacket); err != nil {
            h.log.Debugf("Decrypted data json decode failed")

            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    }

    h.log.Debugf("Packet decoded %v", devicePacket)

    if err := h.piotDevices.ProcessPacket(devicePacket); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
}
