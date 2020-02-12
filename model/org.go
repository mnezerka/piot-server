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
    InfluxDb    string `json:"influxdb"`
    InfluxDbUsername   string `json:"influxdb_username" bson:"influxdb_username"`
    InfluxDbPassword   string `json:"influxdb_password" bson:"influxdb_password"`
    MqttUsername   string `json:"mqtt_username" bson:"mqtt_username"`
    MqttPassword   string `json:"mqtt_password" bson:"mqtt_password"`
    MysqlDb    string `json:"mysqldb"`
    MysqlDbUsername   string `json:"mysqldb_username" bson:"mysqldb_username"`
    MysqlDbPassword   string `json:"mysqldb_password" bson:"mysqldb_password"`
}

// Represents assignment of user to org
type OrgUser struct {
    OrgId       primitive.ObjectID `json:"org_id" bson:"org_id,omitempty"`
    UserId      primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
    Created     int32  `json:"created"`
}
