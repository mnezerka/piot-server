package main_test

import (
    "bytes"
    "crypto/aes"
    "encoding/hex"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "piot-server"
)

func getAdapter(t *testing.T) *main.Adapter {
    log := GetLogger(t)
    db := GetDb(t)
    things := GetThings(t, log, db)
    mqtt := GetMqtt(t, log)
    pdevices := GetPiotDevices(t, log, things, mqtt)

    return main.NewAdapter(log, pdevices, "1234567890123456")
}

/* GET method is not supported */
func TestForbiddenGet(t *testing.T) {
    req, err := http.NewRequest("GET", "/", strings.NewReader(""))
    Ok(t, err)

    rr := httptest.NewRecorder()

    adapter := getAdapter(t)
    adapter.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 405)
}

/* Post data for device that is not registered  */
func TestPacketForUnknownThing(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)

    deviceData := `
    {
        "device": "Device123",
        "readings": [
            {
                "address": "SensorXYZ",
                "t": 23
            }
        ]
    }`

    req, err := http.NewRequest("POST", "/", strings.NewReader(deviceData))
    Ok(t, err)

    rr := httptest.NewRecorder()

    adapter := getAdapter(t)
    adapter.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 200)

    // TODO: Check if device is registered
}

/* Post data in short notation for device that is not registered */
func TestPacketShortNotationForUnknownThing(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)

    deviceData := `
    {
        "d": "Device123",
        "r": [
            {
                "a": "SensorXYZ",
                "t": 23
            }
        ]
    }`

    req, err := http.NewRequest("POST", "/", strings.NewReader(deviceData))
    Ok(t, err)

    rr := httptest.NewRecorder()

    adapter := getAdapter(t)
    adapter.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 200)

    // TODO: Check if device is registered
}

/* Post encrypted data  */
func TestPacketEncrypted(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)

    raw := `
    {
        "device": "Device123",
        "readings": [
            {
                "address": "SensorXYZ",
                "t": 23
            }
        ]
    }`

    /*
    * key The key used for encryption. The key length can be any of 128bit, 192bit, and 256bit.
    * 16-bit key corresponds to 128bit
    */
    key := "1234567890123456"

    size := 16

    padding := size - len(raw) % size
    t.Logf("raw text size: %d, right padding with %d PKCS#7 bytes", len(raw), padding)
    // padding of block for pkcs#7 padding
    if padding > 0 {
        raw = raw + string(bytes.Repeat([]byte{byte(padding)}, padding))
    // add empty block for pkcs#7 padding
    } else {
        raw = raw + string(bytes.Repeat([]byte{byte(size)}, size))
    }
    t.Logf("%s", hex.Dump([]byte(raw)))

    cipher, err := aes.NewCipher([]byte(key))
    Ok(t, err)

    encrypted := make([]byte, len(raw))
    open := []byte(raw)
    for bs, be := 0, size; bs < len(open); bs, be = bs + size, be + size {
        cipher.Encrypt(encrypted[bs:be], open[bs:be])
    }

    req, err := http.NewRequest("POST", "/", strings.NewReader(string(encrypted)))
    Ok(t, err)

    rr := httptest.NewRecorder()

    adapter := getAdapter(t)
    adapter.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, 200)

    // TODO: Check if device is registered
}

