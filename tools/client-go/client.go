package main

import (
    "fmt"
    "log"
    "os"
    "time"
    "github.com/urfave/cli"
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connect(clientId, uri, username, password string) mqtt.Client {

    opts := mqtt.NewClientOptions()
    opts.AddBroker(uri)
    opts.SetClientID(clientId)
    if username != "" {
        opts.SetUsername(username)
        if password != "" {
            opts.SetPassword(password)
        }
    }

    client := mqtt.NewClient(opts)
    token := client.Connect()
    for !token.WaitTimeout(3 * time.Second) {
    }
    if err := token.Error(); err != nil {
        log.Fatal(err)
    }
    return client
}


/*
func listen(uri *url.URL, topic string) {
    client := connect("sub", uri)
    client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
        fmt.Printf("<- [%s] %s\n", msg.Topic(), string(msg.Payload()))
    })
}
*/

func runClient(c *cli.Context) {

    uri := c.GlobalString("mqtt-uri")
    topic := c.GlobalString("topic")
    value := c.GlobalString("value")

    fmt.Printf("MQTT address - %s\n", uri)
    fmt.Printf("MQTT user: %s\n", c.GlobalString("mqtt-user"))
    fmt.Printf("MQTT topic: %s\n", topic)
    fmt.Printf("MQTT value: %s\n", value)

    //topicdev1 := "test/dev1"
    //topicdev2 := "test/dev2"

    // subscribe for messages from my clients
    //go listen(uri, "test/#")

    client := connect("mqtt-client", uri, c.GlobalString("mqtt-user"), c.GlobalString("mqtt-password"))

    defer client.Disconnect(250)

    token := client.Publish(topic, 0, false, value)
    for !token.WaitTimeout(3 * time.Second) { }
    if err := token.Error(); err != nil {
        log.Fatal(err)
    }


    /*
    timer := time.NewTicker(1 * time.Second)
    for t := range timer.C {
        fmt.Printf("-> [%s], message: %s\n", topic, t.String())
        client.Publish(topic, 0, false, t.String())

        fmt.Printf("-> [%s]\n", topicdev1)
        client.Publish(topicdev1, 0, false, "dev1 tick")

        fmt.Printf("-> [%s]\n", topicdev2)
        client.Publish(topicdev2, 0, false, "dev2 tick")
    }
    */
}

func main() {
    app := cli.NewApp()

    app.Name = "MQTT Client"
    app.Action = runClient
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:   "mqtt-uri,q",
            Usage:  "Endpoint of the MQTT message broker",
            EnvVar: "MQTT_URI",
            Value:  "tcp://localhost:1883",
        },
        cli.StringFlag{
            Name:   "mqtt-user",
            Usage:  "Username for mqtt authentication",
            EnvVar: "MQTT_USER",
        },
        cli.StringFlag{
            Name:   "mqtt-password",
            Usage:  "Password for mqtt authentication",
            EnvVar: "MQTT_PASSWORD",
        },
        cli.StringFlag{
            Name:   "topic",
            Usage:  "MQTT Topic",
            EnvVar: "MQTT_TOPIC",
        },
        cli.StringFlag{
            Name:   "value",
            Usage:  "MQTT Value",
            EnvVar: "MQTT_VALUE",
        },
    }

    app.Run(os.Args)

}
