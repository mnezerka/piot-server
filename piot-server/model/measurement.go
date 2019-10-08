package model

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents a single sensor measurement
type Measurement struct {
    Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Created     int32  `json:"created"`
}
