package model

type PiotSensorReading struct {
    Address     string  `json:"address"`
    Temperature float32 `json:"t"`
    Humidity    float32 `json:"h"`
    Pressure    float32 `json:"p"`
}

type PiotDevicePacket struct {
    Device          string  `json:"device"`
    Ip              string  `json:"ip"`
    WifiSSID        string  `json:"wifi-ssid"`
    WifiStrength    float32 `json:"wifi-strength"`
    Time            int32   `json:"time"`
    Readings        []PiotSensorReading `json:"readings"`
}
