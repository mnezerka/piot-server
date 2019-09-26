package model

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents org as stored in database
type Org struct {
    Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Created     int32  `json:"created"`
}
