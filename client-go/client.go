package main

import (
    "fmt"
    "log"
    "net/url"
    "os"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connect(clientId string, uri *url.URL) mqtt.Client {
    opts := createClientOptions(clientId, uri)
    client := mqtt.NewClient(opts)
    token := client.Connect()
    for !token.WaitTimeout(3 * time.Second) {
    }
    if err := token.Error(); err != nil {
        log.Fatal(err)
    }
    return client
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
    opts.SetUsername(uri.User.Username())
    password, _ := uri.User.Password()
    opts.SetPassword(password)
    opts.SetClientID(clientId)
    return opts
}

func listen(uri *url.URL, topic string) {
    client := connect("sub", uri)
    client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
        fmt.Printf("<- [%s] %s\n", msg.Topic(), string(msg.Payload()))
    })
}

func main() {
    // in format tcp://localhost:1883
    uri, err := url.Parse(os.Getenv("MQTT"))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Broker address - tcp://%s\n", uri.Host)

    topic := "test"
    topicdev1 := "test/dev1"
    topicdev2 := "test/dev2"

    // subscribe for messages from my clients
    go listen(uri, "test/#")

    client := connect("pub", uri)
    timer := time.NewTicker(1 * time.Second)
    for t := range timer.C {
        fmt.Printf("-> [%s], message: %s\n", topic, t.String())
        client.Publish(topic, 0, false, t.String())

        fmt.Printf("-> [%s]\n", topicdev1)
        client.Publish(topicdev1, 0, false, "dev1 tick")

        fmt.Printf("-> [%s]\n", topicdev2)
        client.Publish(topicdev2, 0, false, "dev2 tick")
    }
}
