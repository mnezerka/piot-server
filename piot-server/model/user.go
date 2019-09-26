package model

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// Represents user as stored in database
type User struct {
    Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Email       string `json:"email"`
    Password    string `json:"password"`
    Created     int32  `json:"created"`
    OrgId  primitive.ObjectID `json:"org_id" bson:"org_id,omitempty"`
}

type UserProfile struct {
    Email     string `json:"email"`
}
