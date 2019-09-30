package model

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

const THING_TYPE_DEVICE = "device"
const THING_TYPE_SENSOR = "sensor"

// Represents any device or app
type Thing struct {
    Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name        string `json:"name"`
    Type        string `json:"type"`
    Enabled     bool   `json:"enabled"`
    Created     int32  `json:"created"`
}
