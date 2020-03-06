package service

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "strconv"
    "time"
    "github.com/op/go-logging"
    "piot-server/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/tidwall/gjson"
)

const VALUE_YES  = "yes"
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
    PushThingData(ctx context.Context, thing *model.Thing, topic, value string) error
    ProcessMessage(ctx context.Context, topic, payload string)
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    SetUsername(username string)
    SetPassword(password string)
    SetClient(id string)
}

type Mqtt struct {
    Uri string
    Username *string
    Password *string
    Client *string
    client mqtt.Client
}

func NewMqtt(uri string) IMqtt {
    m := &Mqtt{}
    m.Uri = uri

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

func (t *Mqtt) Connect(ctx context.Context) error {
    ctx.Value("log").(*logging.Logger).Infof("Connecting to MQTT broker %s", t.Uri)

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
        topic := fmt.Sprintf("%s/#", TOPIC_ROOT)

        ctx.Value("log").(*logging.Logger).Infof("Connectedt to MQTT broker %s", t.Uri)

        // subscribe for all topcis
        ctx.Value("log").(*logging.Logger).Infof("Subscribing to topic #")
        token := client.Subscribe(topic, 0, func(_ mqtt.Client, msg mqtt.Message) {
            //processUpdate(msg.Topic(), string(msg.Payload()))
            t.ProcessMessage(ctx, msg.Topic(), string(msg.Payload()))
        })
        if !token.WaitTimeout(10 * time.Second) {
            ctx.Value("log").(*logging.Logger).Errorf("Timeout subscribing to topic %s (%s)", topic, token.Error())
        }
        if err := token.Error(); err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("Failed to subscribe to topic %s (%s)", topic, err)
        }

        ctx.Value("log").(*logging.Logger).Infof("Subscribed to topic %s", topic)
    }

    opts.OnConnectionLost = func(client mqtt.Client, err error) {
        ctx.Value("log").(*logging.Logger).Infof("Error: Connection to MQTT broker %s lost (%s)", t.Uri, err.Error())
    }

    // create and start a client using the above ClientOptions
    t.client = mqtt.NewClient(opts)
    if token := t.client.Connect(); token.Wait() && token.Error() != nil {
        ctx.Value("log").(*logging.Logger).Infof("Connection failed (%s)", token.Error())
        return token.Error()
    }

    ctx.Value("log").(*logging.Logger).Infof("Connected to MQTT broker")
    return nil
}

func (t *Mqtt) Disconnect(ctx context.Context) error {
    ctx.Value("log").(*logging.Logger).Infof("Disconnecting from MQTT broker")
    t.client.Disconnect(250)
    return nil
}

func (t *Mqtt) GetThingTopic(ctx context.Context, thing *model.Thing, topic string) (string, error) {
    // get thing org
    orgs := ctx.Value("orgs").(*Orgs)
    org, err := orgs.Get(ctx, thing.OrgId)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("%s/%s/%s/%s", TOPIC_ROOT, org.Name, thing.Name, topic), nil
}

func (t *Mqtt) PushThingData(ctx context.Context, thing *model.Thing, topic, value string) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Push thing data to mqtt broker: %s", thing.Name)

    // post data to MQTT if device is enabled
    if thing.OrgId == primitive.NilObjectID {
        err := fmt.Errorf("Rejecting push to mqtt due to missing organization assignment of thing \"%s\"", thing.Name)
        ctx.Value("log").(*logging.Logger).Infof(err.Error())
        return err
    }

    mqttTopic, err := t.GetThingTopic(ctx, thing, topic)
    if err != nil {
        return err
    }

    ctx.Value("log").(*logging.Logger).Debugf("MQTT Publish, topic: \"%s\", value: \"%s\"", mqttTopic, value)

    token := t.client.Publish(mqttTopic, 0, false, value)
    token.Wait()
    return nil
}

func (t *Mqtt) ProcessDevices(ctx context.Context, org *model.Org, topic, payload string) {
    ctx.Value("log").(*logging.Logger).Debugf("Processing MQTT message with topic \"%s\" for devices in org \"%s\"", topic, org.Name)

    things := ctx.Value("things").(*Things)

    // update availability
    devices, err := things.GetFiltered(ctx, bson.M{"org_id": org.Id, "type": model.THING_TYPE_DEVICE, "availability_topic": topic})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error, falied fetching of org \"%s\" devices: %s", org.Name, err.Error())
        return
    }

    for i := 0; i < len(devices); i++ {

        thing := devices[i]

        // update sensor last seen status
        err = things.TouchThing(ctx, thing.Id)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error: %s", err.Error())
        }
    }

    // update telemetry
    devices, err = things.GetFiltered(ctx, bson.M{"org_id": org.Id, "type": model.THING_TYPE_DEVICE, "telemetry_topic": topic})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error, falied fetching of org \"%s\" devices: %s", org.Name, err.Error())
        return
    }

    for i := 0; i < len(devices); i++ {

        thing := devices[i]

        // update sensor last seen status
        err = things.TouchThing(ctx, thing.Id)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error: %s", err.Error())
        }

        err = things.SetTelemetry(ctx, thing.Id, payload)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error: %s", err.Error())
        }
    }

    // update location
    devices, err = things.GetFiltered(ctx, bson.M{"org_id": org.Id, "type": model.THING_TYPE_DEVICE, "location_topic": topic})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error, falied fetching of org \"%s\" devices: %s", org.Name, err.Error())
        return
    }

    for i := 0; i < len(devices); i++ {


        thing := devices[i]

        // update sensor last seen status
        err = things.TouchThing(ctx, thing.Id)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error: %s", err.Error())
        }

        // decode lat and lng value from json - templates are mandatory
        if thing.LocationLatValue == "" || thing.LocationLngValue == "" {
            ctx.Value("log").(*logging.Logger).Warningf("Ignoring MQTT location message for device %s (\"%s\") due to missing location templates: %s", thing.Id.Hex(), org.Name)
            continue
        }

        var loc model.LocationData
        var haveLat = false;
        var haveLng = false;

        parsedValue := gjson.Get(payload, thing.LocationLatValue)
        if parsedValue.Exists() {
            loc.Latitude, err = strconv.ParseFloat(parsedValue.String(), 8)
            if (err == nil) {
                haveLat = true;
            }
        }
        if !haveLat {
            ctx.Value("log").(*logging.Logger).Warningf("Ignoring MQTT location message for device %s (\"%s\") due to failed parsing of lat value", thing.Id.Hex(), org.Name)
            continue
        }

        parsedValue = gjson.Get(payload, thing.LocationLngValue)
        if parsedValue.Exists() {
            loc.Longitude, err = strconv.ParseFloat(parsedValue.String(), 8)
            if (err == nil) {
                haveLng = true;
            }
        }
        if !haveLng {
            ctx.Value("log").(*logging.Logger).Warningf("Ignoring MQTT location message for device %s (\"%s\") due to failed parsing of lng value", thing.Id.Hex(), org.Name)
            continue
        }

        // if both latitude and longitude were parsed, update thing
        if haveLat && haveLng {
            err = things.SetLocation(ctx, thing.Id, loc)
            if err != nil {
                ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error: %s", err.Error())
            }
        }
    }

}

func (t *Mqtt) ProcessSensors(ctx context.Context, org *model.Org, topic, payload string) {
    ctx.Value("log").(*logging.Logger).Debugf("Processing MQTT message with topic \"%s\" for sensors in org \"%s\"", topic, org.Name)

    // look for sensors attached to this topic from active org
    things := ctx.Value("things").(*Things)

    sensors, err := things.GetFiltered(ctx, bson.M{"org_id": org.Id, "type": model.THING_TYPE_SENSOR, "sensor.measurement_topic": topic})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error, falied fetching of org \"%s\" sensors: %s", org.Name, err.Error())
        return
    }

    // convert orgs to org resolvers
    for i := 0; i < len(sensors); i++ {
        thing := sensors[i]

        value := payload

        ctx.Value("log").(*logging.Logger).Debugf("MQTT sensor measurement value template: \"%s\"", thing.Sensor.MeasurementValue)

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
        err = things.TouchThing(ctx, thing.Id)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error: %s", err.Error())
        }

        // set value to one from incoming payload
        err = things.SetSensorValue(ctx, thing.Id, value)
        if err != nil {
            // report error, but don't interrupt processing
            ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error: %s", err.Error())
        }

        // store it to influx db if configured
        if thing.StoreInfluxDb  {
            influxDb := ctx.Value("influxdb").(IInfluxDb)
            influxDb.PostMeasurement(ctx, thing, value)
        }

        // store it to mysql db if configured
        if thing.StoreMysqlDb {
            mysqlDb := ctx.Value("mysqldb").(IMysqlDb)
            mysqlDb.StoreMeasurement(ctx, thing, value)
        }

    }
}

func (t *Mqtt) ProcessSwitches(ctx context.Context, org *model.Org, topic, payload string) {
    ctx.Value("log").(*logging.Logger).Debugf("Processing MQTT message with topic \"%s\" for switches in org \"%s\"", topic, org.Name)

    // look for sensors attached to this topic from active org
    things := ctx.Value("things").(*Things)

    switches, err := things.GetFiltered(ctx, bson.M{"org_id": org.Id, "type": model.THING_TYPE_SWITCH, "switch.state_topic": topic})
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error, falied fetching of org \"%s\" switches: %s", org.Name, err.Error())
        return
    }

    // convert orgs to org resolvers
    for i := 0; i < len(switches); i++ {

        thing := switches[i]


        // update sensor last seen status
        err = things.TouchThing(ctx, thing.Id)
        if err != nil {
            ctx.Value("log").(*logging.Logger).Errorf("MQTT processing error: %s", err.Error())
        }

        dbValue := ""
        switch(payload) {
        case thing.Switch.StateOn:
            err = things.SetSwitchState(ctx, thing.Id, true)
            dbValue = "1"
        case thing.Switch.StateOff:
            err = things.SetSwitchState(ctx, thing.Id, false)
            dbValue = "0"
        default:
            err = errors.New("Unknown switch state")
        }
        if err != nil {
            ctx.Value("log").(*logging.Logger).Warningf("Issue with processing of switch %s MQTT state messsage: %s", thing.Name, err.Error())
        }

        // store it to influx db if configured
        if thing.StoreInfluxDb {
            influxDb := ctx.Value("influxdb").(IInfluxDb)
            influxDb.PostSwitchState(ctx, thing, dbValue)
        }

        // store it to mysql db if configured
        if thing.StoreMysqlDb {
            mysqlDb := ctx.Value("mysqldb").(IMysqlDb)
            mysqlDb.StoreSwitchState(ctx, thing, dbValue)
        }
    }
}

// Process message received from MQTT broker for org subscription
func (t *Mqtt) ProcessMessage(ctx context.Context, topic, payload string) {
    ctx.Value("log").(*logging.Logger).Debugf("Recieved MQTT message (topic: %s, val: %s)", topic, payload)

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
    orgs := ctx.Value("orgs").(*Orgs)
    org, err := orgs.GetByName(ctx, topicParts[1])
    if err != nil {
        // unknown organization
        ctx.Value("log").(*logging.Logger).Warningf("MQTT processing error, unknown org: %s (%s)", topicParts[1], err.Error())
        return
    }


    t.ProcessDevices(ctx, org, topicThing, payload);

    t.ProcessSensors(ctx, org, topicThing, payload);

    t.ProcessSwitches(ctx, org, topicThing, payload);

    // look for switches attached to this topic
}
