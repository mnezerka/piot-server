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
    Orgs        []Org  `json:"orgs"`
    ActiveOrgId primitive.ObjectID `json:"active_org_id" bson:"active_org_id"`
}
