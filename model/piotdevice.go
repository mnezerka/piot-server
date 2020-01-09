package model

type PiotSensorReading struct {
    Address        string   `json:"address"`
    AddressShort   string   `json:"a"`
    Temperature    *float32 `json:"t,omitempty"`
    Humidity       *float32 `json:"h,omitempty"`
    Pressure       *float32 `json:"p,omitempty"`
}

type PiotDevicePacket struct {
    Device          string   `json:"device"`
    DeviceShort     string   `json:"d"`
    Ip              *string  `json:"ip,omitempty"`
    WifiSSID        *string  `json:"wifi-ssid,omitempty"`
    WifiStrength    *float32 `json:"wifi-strength,omitempty"`
    Readings        []PiotSensorReading `json:"readings"`
    ReadingsShort   []PiotSensorReading `json:"r"`
}
