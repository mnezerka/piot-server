package model

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents any device or app
type Thing struct {
    Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name        string `json:"name"`
    Type        string `json:"type"`
    Available   bool `json:"available"`
    Created     int32  `json:"created"`
}
