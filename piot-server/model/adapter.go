package model

type PiotSensorReading struct {
    Address     string   `json:"address"`
    Temperature *float32 `json:"t,omitempty"`
    Humidity    *float32 `json:"h,omitempty"`
    Pressure    *float32 `json:"p,omitempty"`
}

type PiotDevicePacket struct {
    Device          string   `json:"device"`
    Ip              *string  `json:"ip,omitempty"`
    WifiSSID        *string  `json:"wifi-ssid,omitempty"`
    WifiStrength    *float32 `json:"wifi-strength,omitempty"`
    Time            *int32   `json:"time,omitempty"`
    Readings        []PiotSensorReading `json:"readings"`
}
