package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const VALUE_YES = "yes"
const VALUE_NO = "no"

const TOPIC_UNIT = "unit"

const TOPIC_AVAILABLE = "available"

const TOPIC_NET = "net"
const TOPIC_IP = "net/ip"
const TOPIC_WIFI_SSID = "net/wifi/ssid"
const TOPIC_WIFI_STRENGTH = "net/wifi/strength"

const TOPIC_TEMP = "temperature"
const TOPIC_PRESSURE = "pressure"
const TOPIC_HUMIDITY = "humidity"

const TOPIC_ROOT = "org"

type IMqtt interface {
	PushThingData(thing *Thing, topic, value string) error
	ProcessMessage(topic, payload string)
	Connect(subscribe bool) error
	Disconnect() error
	SetUsername(username string)
	SetPassword(password string)
	SetClient(id string)
}

type Mqtt struct {
	log      *logging.Logger
	things   *Things
	orgs     *Orgs
	influxDb IInfluxDb
	mysqlDb  IMysqlDb

	Uri      string
	Username *string
	Password *string
	Client   *string
	client   mqtt.Client
}

func NewMqtt(uri string, log *logging.Logger, things *Things, orgs *Orgs, influxDb IInfluxDb, mysqlDb IMysqlDb) IMqtt {
	m := &Mqtt{log: log, Uri: uri, things: things, orgs: orgs, influxDb: influxDb, mysqlDb: mysqlDb}

	return m
}

func (t *Mqtt) SetUsername(username string) {
	t.Username = &username
}

func (t *Mqtt) SetPassword(password string) {
	t.Password = &password
}

func (t *Mqtt) SetClient(id string) {
	t.Client = &id
}

func (t *Mqtt) Connect(subscribe bool) error {
	t.log.Infof("Connecting to MQTT broker %s", t.Uri)

	// create a ClientOptions struct setting the broker address, clientid, turn
	// off trace output and set the default message handler
	opts := mqtt.NewClientOptions().AddBroker(t.Uri)
	opts.SetClientID(*t.Client)
	if t.Username != nil {
		opts.SetUsername(*t.Username)
	}
	if t.Password != nil {
		opts.SetPassword(*t.Password)
	}

	opts.OnConnect = func(client mqtt.Client) {

		t.log.Infof("Connectedt to MQTT broker %s", t.Uri)
		if subscribe {

			topic := fmt.Sprintf("%s/#", TOPIC_ROOT)

			// subscribe for all topcis
			t.log.Infof("Subscribing to topic #")
			token := client.Subscribe(topic, 0, func(_ mqtt.Client, msg mqtt.Message) {
				//processUpdate(msg.Topic(), string(msg.Payload()))
				t.ProcessMessage(msg.Topic(), string(msg.Payload()))
			})
			if !token.WaitTimeout(10 * time.Second) {
				t.log.Errorf("Timeout subscribing to topic %s (%s)", topic, token.Error())
			}
			if err := token.Error(); err != nil {
				t.log.Errorf("Failed to subscribe to topic %s (%s)", topic, err)
			}

			t.log.Infof("Subscribed to topic %s", topic)
		}
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		t.log.Infof("Error: Connection to MQTT broker %s lost (%s)", t.Uri, err.Error())
	}

	// create and start a client using the above ClientOptions
	t.client = mqtt.NewClient(opts)
	if token := t.client.Connect(); token.Wait() && token.Error() != nil {
		t.log.Infof("Connection failed (%s)", token.Error())
		return token.Error()
	}

	t.log.Infof("Connected to MQTT broker")
	return nil
}

func (t *Mqtt) Disconnect() error {
	t.log.Infof("Disconnecting from MQTT broker")
	t.client.Disconnect(250)
	return nil
}

func (t *Mqtt) GetThingTopic(thing *Thing, topic string) (string, error) {
	// get thing org
	org, err := t.orgs.Get(thing.OrgId)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s/%s", TOPIC_ROOT, org.Name, thing.Name, topic), nil
}

func (t *Mqtt) PushThingData(thing *Thing, topic, value string) error {
	t.log.Debugf("Push thing data to mqtt broker: %s", thing.Name)

	// post data to MQTT if device is enabled
	if thing.OrgId == primitive.NilObjectID {
		err := fmt.Errorf("Rejecting push to mqtt due to missing organization assignment of thing \"%s\"", thing.Name)
		t.log.Infof(err.Error())
		return err
	}

	mqttTopic, err := t.GetThingTopic(thing, topic)
	if err != nil {
		return err
	}

	t.log.Debugf("MQTT Publish, topic: \"%s\", value: \"%s\"", mqttTopic, value)

	token := t.client.Publish(mqttTopic, 0, false, value)
	token.Wait()
	return nil
}

func (t *Mqtt) ProcessAll(org *Org, topic, payload string) {
	t.log.Debugf("Processing MQTT message with topic \"%s\" for all things in org \"%s\"", topic, org.Name)

	// update battery level
	things, err := t.things.GetFiltered(bson.M{"org_id": org.Id, "battery_mqtt_topic": topic})
	if err != nil {
		t.log.Errorf("MQTT processing error, falied fetching of org \"%s\" devices: %s", org.Name, err.Error())
		return
	}
	for i := 0; i < len(things); i++ {

		thing := things[i]

		// update sensor last seen status
		err = t.things.TouchThing(thing.Id)
		if err != nil {
			t.log.Errorf("MQTT processing error: %s", err.Error())
		}

		// battery level value json key (optional value)
		var level int32
		level_string := payload

		if thing.BatteryMqttLevelValue != "" {
			parsedValue := gjson.Get(payload, thing.BatteryMqttLevelValue)
			if parsedValue.Exists() {
				level_string = parsedValue.String()
			}
		}

		parsedLevel, err := strconv.ParseInt(level_string, 10, 32)
		if err == nil {
			level = int32(parsedLevel)
		} else {
			t.log.Warningf("Ignoring MQTT battery message level for device %s (\"%s\") due to failed parsing of value (must be int)", thing.Id.Hex(), org.Name)
		}

		// store -> persistent storage
		err = t.things.SetBatteryLevel(thing.Id, level)
		if err != nil {
			t.log.Errorf("MQTT processing error: %s", err.Error())
		}

		// store -> time series in influxdb
		if thing.BatteryLevelTracking {
			t.influxDb.PostBatteryLevel(thing, level)
		}

	}
}

func (t *Mqtt) ProcessDevices(org *Org, topic, payload string) {
	t.log.Debugf("Processing MQTT message with topic \"%s\" for devices in org \"%s\"", topic, org.Name)

	// update availability
	devices, err := t.things.GetFiltered(bson.M{"org_id": org.Id, "type": THING_TYPE_DEVICE, "availability_topic": topic})
	if err != nil {
		t.log.Errorf("MQTT processing error, falied fetching of org \"%s\" devices: %s", org.Name, err.Error())
		return
	}
	for i := 0; i < len(devices); i++ {

		thing := devices[i]

		// update sensor last seen status
		err = t.things.TouchThing(thing.Id)
		if err != nil {
			t.log.Errorf("MQTT processing error: %s", err.Error())
		}
	}

	// update telemetry
	devices, err = t.things.GetFiltered(bson.M{"org_id": org.Id, "type": THING_TYPE_DEVICE, "telemetry_topic": topic})
	if err != nil {
		t.log.Errorf("MQTT processing error, falied fetching of org \"%s\" devices: %s", org.Name, err.Error())
		return
	}
	for i := 0; i < len(devices); i++ {

		thing := devices[i]

		// update sensor last seen status
		err = t.things.TouchThing(thing.Id)
		if err != nil {
			t.log.Errorf("MQTT processing error: %s", err.Error())
		}

		err = t.things.SetTelemetry(thing.Id, payload)
		if err != nil {
			t.log.Errorf("MQTT processing error: %s", err.Error())
		}
	}

	// update location
	devices, err = t.things.GetFiltered(bson.M{"org_id": org.Id, "type": THING_TYPE_DEVICE, "loc_mqtt_topic": topic})
	if err != nil {
		t.log.Errorf("MQTT processing error, falied fetching of org \"%s\" devices: %s", org.Name, err.Error())
		return
	}
	for i := 0; i < len(devices); i++ {

		thing := devices[i]

		// update sensor last seen status
		err = t.things.TouchThing(thing.Id)
		if err != nil {
			t.log.Errorf("MQTT processing error: %s", err.Error())
		}

		// decode lat and lng value from json - templates are mandatory
		if thing.LocationMqttLatValue == "" || thing.LocationMqttLngValue == "" {
			t.log.Warningf("Ignoring MQTT location message for device %s (\"%s\") due to missing location templates: %s", thing.Id.Hex(), org.Name)
			continue
		}

		//var loc LocationData
		var lat, lng float64
		var sat, ts int32
		var haveLat = false
		var haveLng = false

		// parse LAT
		parsedValue := gjson.Get(payload, thing.LocationMqttLatValue)
		if parsedValue.Exists() {
			lat, err = strconv.ParseFloat(parsedValue.String(), 8)
			if err == nil {
				haveLat = true
			}
		}
		if !haveLat {
			t.log.Warningf("Ignoring MQTT location message for device %s (\"%s\") due to failed parsing of lat value", thing.Id.Hex(), org.Name)
			continue
		}

		// parse LNG
		parsedValue = gjson.Get(payload, thing.LocationMqttLngValue)
		if parsedValue.Exists() {
			lng, err = strconv.ParseFloat(parsedValue.String(), 8)
			if err == nil {
				haveLng = true
			}
		}
		if !haveLng {
			t.log.Warningf("Ignoring MQTT location message for device %s (\"%s\") due to failed parsing of lng value", thing.Id.Hex(), org.Name)
			continue
		}

		// parse timestamp (optional value)
		if thing.LocationMqttTsValue != "" {
			parsedValue := gjson.Get(payload, thing.LocationMqttTsValue)
			if parsedValue.Exists() {
				parsedTs, err := strconv.ParseInt(parsedValue.String(), 10, 32)

				if err == nil {
					ts = int32(parsedTs)
				} else {
					t.log.Warningf("Ignoring MQTT location message timestamp for device %s (\"%s\") due to failed parsing of value", thing.Id.Hex(), org.Name)
				}
			}
		}

		// parse satelites (optional value)
		if thing.LocationMqttSatValue != "" {
			parsedValue := gjson.Get(payload, thing.LocationMqttSatValue)
			if parsedValue.Exists() {
				parsedSat, err := strconv.ParseInt(parsedValue.String(), 10, 32)

				if err == nil {
					sat = int32(parsedSat)
				} else {
					t.log.Warningf("Ignoring MQTT location message satelites for device %s (\"%s\") due to failed parsing of value", thing.Id.Hex(), org.Name)
				}
			}
		}

		// if both latitude and longitude were parsed, update thing
		if haveLat && haveLng {
			// if date is not set (didn't come in payload), use current date
			if ts == 0 {
				ts = int32(time.Now().Unix())
			}

			err = t.things.SetLocation(thing.Id, lat, lng, sat, ts)
			if err != nil {
				t.log.Errorf("MQTT processing error: %s", err.Error())
			}

			if thing.LocationTracking {
				t.influxDb.PostLocation(thing, lat, lng, sat, ts)
			}
		}
	}
}

func (t *Mqtt) ProcessSensors(org *Org, topic, payload string) {
	t.log.Debugf("Processing MQTT message with topic \"%s\" for sensors in org \"%s\"", topic, org.Name)

	// look for sensors attached to this topic from active org
	sensors, err := t.things.GetFiltered(bson.M{"org_id": org.Id, "type": THING_TYPE_SENSOR, "sensor.measurement_topic": topic})
	if err != nil {
		t.log.Errorf("MQTT processing error, falied fetching of org \"%s\" sensors: %s", org.Name, err.Error())
		return
	}

	// convert orgs to org resolvers
	for i := 0; i < len(sensors); i++ {
		thing := sensors[i]

		value := payload

		t.log.Debugf("MQTT sensor measurement value template: \"%s\"", thing.Sensor.MeasurementValue)

		// decode value from json in case value has template
		if thing.Sensor.MeasurementValue != "" {
			parsedValue := gjson.Get(payload, thing.Sensor.MeasurementValue)
			if !parsedValue.Exists() {
				value = ""
			} else {
				value = parsedValue.String()
			}
		}

		// update sensor last seen status
		err = t.things.TouchThing(thing.Id)
		if err != nil {
			t.log.Errorf("MQTT processing error: %s", err.Error())
		}

		// set value to one from incoming payload
		err = t.things.SetSensorValue(thing.Id, value)
		if err != nil {
			// report error, but don't interrupt processing
			t.log.Errorf("MQTT processing error: %s", err.Error())
		}

		// store it to influx db if configured
		if thing.StoreInfluxDb {
			t.influxDb.PostMeasurement(thing, value)
		}

		// store it to mysql db if configured
		if thing.StoreMysqlDb {
			t.mysqlDb.StoreMeasurement(thing, value)
		}

	}
}

func (t *Mqtt) ProcessSwitches(org *Org, topic, payload string) {
	t.log.Debugf("Processing MQTT message with topic \"%s\" for switches in org \"%s\"", topic, org.Name)

	// look for sensors attached to this topic from active org
	switches, err := t.things.GetFiltered(bson.M{"org_id": org.Id, "type": THING_TYPE_SWITCH, "switch.state_topic": topic})
	if err != nil {
		t.log.Errorf("MQTT processing error, falied fetching of org \"%s\" switches: %s", org.Name, err.Error())
		return
	}

	// convert orgs to org resolvers
	for i := 0; i < len(switches); i++ {

		thing := switches[i]

		// update sensor last seen status
		err = t.things.TouchThing(thing.Id)
		if err != nil {
			t.log.Errorf("MQTT processing error: %s", err.Error())
		}

		dbValue := ""
		switch payload {
		case thing.Switch.StateOn:
			err = t.things.SetSwitchState(thing.Id, true)
			dbValue = "1"
		case thing.Switch.StateOff:
			err = t.things.SetSwitchState(thing.Id, false)
			dbValue = "0"
		default:
			err = errors.New("Unknown switch state")
		}
		if err != nil {
			t.log.Warningf("Issue with processing of switch %s MQTT state messsage: %s", thing.Name, err.Error())
		}

		// store it to influx db if configured
		if thing.StoreInfluxDb {
			t.influxDb.PostSwitchState(thing, dbValue)
		}

		// store it to mysql db if configured
		if thing.StoreMysqlDb {
			t.mysqlDb.StoreSwitchState(thing, dbValue)
		}
	}
}

// Process message received from MQTT broker for org subscription
func (t *Mqtt) ProcessMessage(topic, payload string) {
	t.log.Debugf("Recieved MQTT message (topic: %s, val: %s)", topic, payload)

	topicParts := strings.Split(topic, "/")

	// skip topics that don't contain org section (first two
	// parts
	if len(topicParts) < 3 {
		return
	}

	// skip topics that doesn't belong to organization root
	if topicParts[0] != TOPIC_ROOT {
		return
	}

	topicThing := strings.Join(topicParts[2:], "/")

	// get org ID
	org, err := t.orgs.GetByName(topicParts[1])
	if err != nil {
		// unknown organization
		t.log.Warningf("MQTT processing error, unknown org: %s (%s)", topicParts[1], err.Error())
		return
	}

	t.ProcessAll(org, topicThing, payload)
	t.ProcessDevices(org, topicThing, payload)
	t.ProcessSensors(org, topicThing, payload)
	t.ProcessSwitches(org, topicThing, payload)
}
