package main

import (
    "bytes"
    "crypto/aes"
    "encoding/json"
    "encoding/hex"
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

    bodyLen := len(body)

    h.log.Debugf("Packet body length: %d", bodyLen)

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
        if bodyLen % size != 0 {
            h.log.Errorf("Invalid length of body for decryption %d", bodyLen)
            WriteErrorResponse(w, errors.New("Invalid length of body for decryption"), http.StatusBadRequest)
            return
        }

        // json decode from raw data failed => try to decrypt first
        cipher, _ := aes.NewCipher([]byte(h.password))

        decrypted := make([]byte, bodyLen)

        // decrypt by individual blocks
        decryptedLen := 0
        for bs, be := 0, size; bs < bodyLen; bs, be = bs + size, be + size {
            cipher.Decrypt(decrypted[bs:be], body[bs:be])

            // last block needs special attention due to padding that must be removed
            // before json parsing
            if bs + size == bodyLen {
                // strip pkcs7 padding
                stripped, err := pkcs7strip(decrypted[bs:be], size)
                if err != nil {
                    h.log.Errorf("PKCS#7 padding stripping failed (%e)", err.Error())
                    WriteErrorResponse(w, errors.New("Wrong PKCS#7 padding of encrypted content"), http.StatusBadRequest)
                }

                decryptedLen += len(stripped)

            } else {
                decryptedLen += size
            }
        }

        h.log.Debugf("Decrypted message <%s>", decrypted[:decryptedLen])
        h.log.Debugf("%s", hex.Dump(decrypted))

        // try to decode decrypted data
        if err := json.Unmarshal(decrypted[:decryptedLen], &devicePacket); err != nil {
            h.log.Debugf("Decrypted data json decode failed (%s)", err.Error())

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


// pkcs7strip remove pkcs7 padding
func pkcs7strip(data []byte, blockSize int) ([]byte, error) {

    length := len(data)

    // no empty blocks can exist
    if length == 0 {
        return nil, errors.New("pkcs7: data is empty")
    }

    // all bytes are always filled (padding values are always non zero)
    if length % blockSize != 0 {
        return nil, errors.New("pkcs7: data is not block-aligned")
    }

    // get number of bytes used for padding from last block byte
    padLen := int(data[length-1])

    // generate sequence of bytes that should match end of the block
    ref := bytes.Repeat([]byte{byte(padLen)}, padLen)

    // check if padding is encoded correctly - it must be smaller than block size,
    // non zero and match generated sequence
    if padLen > blockSize || padLen == 0 || !bytes.HasSuffix(data, ref) {
        return nil, errors.New("pkcs7: invalid padding")
    }

    return data[:length-padLen], nil
}
