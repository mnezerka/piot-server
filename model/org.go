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

// Represents assignment of user to org
type OrgUser struct {
    OrgId       primitive.ObjectID `json:"org_id" bson:"org_id,omitempty"`
    UserId      primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
    Created     int32  `json:"created"`
}
