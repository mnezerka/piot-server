package service

import (
    "context"
    "fmt"
    "strings"
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
        ctx.Value("log").(*logging.Logger).Debugf("have sensor")

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
        if thing.Sensor.StoreInfluxDb  {
            influxDb := ctx.Value("influxdb").(IInfluxDb)
            influxDb.PostMeasurement(ctx, thing, value)
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



    t.ProcessSensors(ctx, org, topicThing, payload);

    // look for switches attached to this topic
}
