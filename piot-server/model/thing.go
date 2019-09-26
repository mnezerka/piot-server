package model

// Represents any device or app
type Thing struct {
    Name        string `json:"name"`
    Type        string `json:"type"`
    Available   bool `json:"available"`
    Created     int32  `json:"created"`
}
