package model

// Represents device
type Device struct {
    Name        string `json:"name"`
    Type        string `json:"type"`
    Available   bool `json:"available"`
    Created     int32  `json:"created"`
}
