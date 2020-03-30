package model

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

const THING_TYPE_DEVICE = "device"
const THING_TYPE_SENSOR = "sensor"
const THING_TYPE_SWITCH = "switch"

const THING_CLASS_TEMPERATURE = "temperature"
const THING_CLASS_HUMIDITY = "humidity"
const THING_CLASS_PRESSURE = "pressure"

// Represents any device or app
type Thing struct {

    // unique id of the thing
    Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`

    // PIOT identification of the thing, this
    // arrtibute is used for processing messages coming
    // from PIOT chips via adapter
    PiotId      string `json:"piot_id" bson:"piot_id"`

    // name of the thing
    Name        string `json:"name" bson:"name"`

    // optional description
    Description string `json:"description" bson:"description"`

    // optional alias of the thing - used e.g. in graphs
    Alias       string `json:"alias" bson:"alias"`

    // basic classification of things - device vs sensor
    Type        string `json:"type" bson:"type"`

    // is thing enabled?
    Enabled     bool   `json:"enabled" bson:"enabled"`

    // date of thing creation
    Created     int32  `json:"created" bson:"created"`

    // id of the organization thing belongs to
    OrgId       primitive.ObjectID `json:"org_id" bson:"org_id"`

    // ?
    Org         *Org `json:"org" bson:"org"`

    // id of the parent thing (e.g. chip - sensor relation)
    ParentId    primitive.ObjectID `json:"parent_id" bson:"parent_id"`

    // is thing available
    Available   bool   `json:"available" bson:"available"`

    // time the thing was seen last time
    LastSeen    int32  `json:"last_seen" bson:"last_seen"`

    // maximal interval (in seconds) for which device can be unseen
    // (see LastSeen attribute). Default value is 0, which disables
    // this feature
    LastSeenInterval    int32  `json:"last_seen_interval" bson:"last_seen_interval"`

    // The MQTT topic subscribed to receive thing availability
    AvailabilityTopic   string `json:"availability_topic" bson:"availability_topic"`
    AvailabilityYes     string `json:"availability_yes" bson:"availability_yes"`
    AvailabilityNo      string `json:"availability_no" bson:"availability_no"`

    // the MQTT topic for receiving telemetry information
    TelemetryTopic      string `json:"telemetry_topic" bson:"telemetry_topic"`

    // time the thing was seen last time
    Telemetry           string `json:"telemetry" bson:"telemetry"`

    // Enable or Disable pushing values to organization assigned Influx database
    StoreInfluxDb bool `json:"store_influxdb" bson:"store_influxdb"`

    // Enable or Disable storing values to mysql db assigned to organization
    StoreMysqlDb bool `json:"store_mysqldb" bson:"store_mysqldb"`

    // minimal interval (in seconds) for storing values to mysql db,
    // more values to be stored in same inteval will be ignored, only
    // first one will be stored
    StoreMysqlDbInterval int32 `json:"store_mysqldb_interval" bson:"store_mysqldb_interval"`

     // The latitude in degrees. It must be in the range [-90.0, +90.0].
    LocationLatitude float64 `json:"loc_lat" bson:"loc_lat"`

    // The longitude in degrees. It must be in the range [-180.0, +180.0].
    LocationLongitude float64  `json:"loc_lng" bson:"loc_lng"`

    // The date when location was taken (unix timestamp)
    LocationSatelites int32  `json:"loc_sat" bson:"loc_sat"`

    // The date when location was taken (unix timestamp)
    LocationTs int32  `json:"loc_ts" bson:"loc_ts"`

    // The MQTT topic subscribed to receive thing location
    LocationMqttTopic   string `json:"loc_mqtt_topic" bson:"loc_mqtt_topic"`
    LocationMqttLatValue string `json:"loc_mqtt_lat_value" bson:"loc_mqtt_lat_value"`
    LocationMqttLngValue string `json:"loc_mqtt_lng_value" bson:"loc_mqtt_lng_value"`
    LocationMqttTsValue string `json:"loc_mqtt_ts_value" bson:"loc_mqtt_ts_value"`
    LocationMqttSatValue string `json:"loc_mqtt_sat_value" bson:"loc_mqtt_sat_value"`

    // Persistency of location changes
    LocationTracking bool `json:"loc_tracking" bson:"loc_tracking"`

    // The unit of measurement that the sensor is expressed in.
    Sensor SensorData `json:"sensor" bson:"sensor"`

    // The unit of measurement that the sensor is expressed in.
    Switch SwitchData `json:"switch" bson:"switch"`
}

// Represents measurements for things that are sensors
type SensorData struct {

    // Last value
    Value string `json:"value" bson:"value"`

    // The MQTT topic where sensor values are published
    MeasurementTopic string `json:"measurement_topic" bson:"measurement_topic"`

    // The template for parsing value from MQTT payload, empty value means use
    // payload as it is. Else the value is extraced according to
    // https://github.com/tidwall/gjson
    MeasurementValue string `json:"measurement_value" bson:"measurement_value"`

    // Time when last measurement was received
    MeasurementLast int32 `json:"measurement_last" bson:"measurement_last"`

    // Type of sensor measurement
    Class string `json:"class" bson:"class"`

    // Defines the number of seconds since last measurement for which the
    // measurement is valid
    Validity int32  `json:"validity" bson:"validity"`

    // The unit of measurement that the sensor is expressed in.
    Unit string `json:"unit bson"unit""`
}

// Represents switch (e.g. high voltage power switch)
type SwitchData struct {

    State bool `json:"state" bson:"state"`

    // Topic to send commands
    CommandTopic string `json:"command_topic" bson:"command_topic"`

    // Send command to switch on
    CommandOn string `json:"command_on" bson:"command_on"`

    // Send command to switch off
    CommandOff string `json:"command_off" bson:"command_off"`

    // Topic to receive switch state (ON or OFF)
    StateTopic string `json:"state_topic" bson:"state_topic"`

    // Value that represents ON state
    StateOn string `json:"state_on" bson:"state_on"`

    // Value that represents OFF state
    StateOff string `json:"state_off" bson:"state_off"`
}
