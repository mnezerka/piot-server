package service

import (
    "context"
    "fmt"
    "time"
    "github.com/op/go-logging"
    "piot-server/model"
    "go.mongodb.org/mongo-driver/bson/primitive"
    mqtt "github.com/eclipse/paho.mqtt.golang"
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


type IMqtt interface {
    PushThingData(ctx context.Context, thing *model.Thing, topic, value string) error
    ProcessMessage(ctx context.Context, topic, payload string)
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    SetUsername(username string)
    SetPassword(password string)
}

type Mqtt struct {
    Uri string
    Username *string
    Password *string
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

func (t *Mqtt) Connect(ctx context.Context) error {
    ctx.Value("log").(*logging.Logger).Infof("Connecting to MQTT broker %s", t.Uri)

    // create a ClientOptions struct setting the broker address, clientid, turn
    // off trace output and set the default message handler
    opts := mqtt.NewClientOptions().AddBroker(t.Uri)
    opts.SetClientID("piot-server")
    if t.Username != nil {
        opts.SetUsername(*t.Username)
    }
    if t.Password != nil {
        opts.SetPassword(*t.Password)
    }

    opts.OnConnect = func(client mqtt.Client) {
        const topic string = "org/#"

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

func (t *Mqtt) PushThingData(ctx context.Context, thing *model.Thing, topic, value string) (error) {
    ctx.Value("log").(*logging.Logger).Debugf("Push thing data to mqtt broker: %s", thing.Name)

    // post data to MQTT if device is enabled
    if thing.OrgId == primitive.NilObjectID {
        err := fmt.Errorf("Rejecting push to mqtt due to missing organization assignment of thing \"%s\"", thing.Name)
        ctx.Value("log").(*logging.Logger).Infof(err.Error())
        return err
    }

    // get thing org
    orgs := ctx.Value("orgs").(*Orgs)
    org, err := orgs.Get(ctx, thing.OrgId)
    if err != nil {
        return err
    }

    mqttTopic := fmt.Sprintf("org/%s/%s/%s", org.Name, thing.Name, topic)

    ctx.Value("log").(*logging.Logger).Debugf("MQTT Publish, topic: \"%s\", value: \"%s\"", mqttTopic, value)

    token := t.client.Publish(mqttTopic, 0, false, value)
    token.Wait()
    return nil
}

func (t *Mqtt) ProcessMessage(ctx context.Context, topic, payload string) {
    ctx.Value("log").(*logging.Logger).Debugf("Recieved MQTT message (topic: %s, val: %s)", topic, payload)
}
