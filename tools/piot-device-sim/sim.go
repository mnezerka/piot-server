package main

import (
    "bytes"
    "log"
    "fmt"
    "encoding/json"
    "net/http"
    "os"
    "time"
    "github.com/urfave/cli"
)

type MsgParams struct {
    Device string
    Sensor string
}

////////////// PeekGenerator ////////////////////

type PeakGenerator struct {
    Params MsgParams
    Value int
    Start int
    Step int
    End int
}

func NewPeakGenerator(params MsgParams, start, step, end int) *PeakGenerator {
    pg := &PeakGenerator{}
    pg.Params = params
    pg.Value = start
    pg.Start = start
    pg.Step = step
    pg.End = end
    return pg
}

func (pg *PeakGenerator) GetValue() int {
    result := pg.Value
    pg.Value += pg.Step
    if pg.Value > pg.End {
        pg.Value = pg.Start
    }
    return result
}

func (pg *PeakGenerator) GetMsg() map[string] interface{} {

    message := map[string] interface{} {
        "device": pg.Params.Device,
        "readings": [](map[string] interface{}) {{"address": pg.Params.Sensor, "t": pg.GetValue()}},
    }

    return message
}

////////////// Main ////////////////////


func sendMsg(uri string, g *PeakGenerator) {

    /*
    message := map[string] interface{} {
        "device": device,
        "readings": [](map[string] interface{}) {{"address": sensor, "t": 23}},
    }
    */

    message := g.GetMsg();

    bytesRepresentation, err := json.Marshal(message)
    if err != nil {
        log.Fatalln(err)
    }

    fmt.Printf("-> %s\n", string(bytesRepresentation))

    resp, err := http.Post(uri, "application/json", bytes.NewBuffer(bytesRepresentation))
    if err != nil {
        log.Fatalln(err)
    }

    fmt.Printf("<- %s\n", resp.Status);
}

func sim(c *cli.Context) {

    uri := c.GlobalString("uri")
    mode := c.GlobalString("mode")
    interval := c.GlobalUint("interval")

    msg_params := MsgParams{}
    msg_params.Device = c.GlobalString("device")
    msg_params.Sensor = c.GlobalString("sensor")

    fmt.Printf("PIOT address: %s\n", uri)
    fmt.Printf("Simulation mode: %s\n", mode)
    fmt.Printf("Device ID: %s\n", msg_params.Device)
    fmt.Printf("Sensor ID: %s\n", msg_params.Sensor)
    fmt.Printf("Interval: %d\n", interval)

    g := NewPeakGenerator(msg_params, 23, 1, 50);

    sendMsg(uri, g);
    timer := time.NewTicker(time.Duration(interval) * time.Second)
    for range timer.C {
        sendMsg(uri, g);
    }
}

func main() {
    app := cli.NewApp()

    app.Name = "MQTT Client"
    app.Action = sim
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:   "uri",
            Usage:  "Endpoint for pushing data (e.g. piot adapter)",
            EnvVar: "PIOT_URI",
            Value:  "http://localhost:9096",
        },
        cli.StringFlag{
            Name:   "mode",
            Usage:  "Simulation mode",
            EnvVar: "SIM_MODE",
            Value:  "peaks",
        },
        cli.StringFlag{
            Name:   "device",
            Usage:  "Device ID",
            EnvVar: "SIM_DEVICE",
            Value:  "SIMDEVICE",
        },
        cli.StringFlag{
            Name:   "sensor",
            Usage:  "Sensor ID",
            EnvVar: "SIM_SENSOR",
            Value: "SIMSENSOR",
        },
        cli.UintFlag{
            Name:   "interval,i",
            Usage:  "Interval in seconds",
            EnvVar: "SIM_INTERVAL",
            Value: 5,
        },

    }

    app.Run(os.Args)

}
