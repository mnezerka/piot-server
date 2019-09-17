package model

// Represents single reading (one sensor)
type Reading struct {
    Address         string `json:"address"`
    Temperature     string `json:"t"`
    Humidity        string `json:"h"`
}

// Represents incoming metering data
type Request struct {
    Device          string `json:"device"`
    IP              string `json:"ip"`
    WifiSSID        string `json:"wifi-ssid"`
    WifiStrength    string `json:"wifi-strength"`
    Created         int32  `json:"created"`
    Readings        []Reading `json:"readings"`
}
