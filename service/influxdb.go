package service

import (
    //"bytes"
    "context"
    //"net/http"
    //"net/url"
    "path"
    "fmt"
    "net/url"
    "piot-server/model"
    "github.com/op/go-logging"
)

type IInfluxDb interface {
    PostMeasurement(ctx context.Context, thing *model.Thing, value string)
    PostSwitchState(ctx context.Context, thing *model.Thing, value string)
}

type InfluxDb struct {
    Uri string
    Username string
    Password string
}

func NewInfluxDb(uri, username, password string) IInfluxDb {
    db := &InfluxDb{}
    db.Uri = uri
    db.Username = username
    db.Password = password

    return db
}

func (db *InfluxDb) PostMeasurement(ctx context.Context, thing *model.Thing, value string) {
    ctx.Value("log").(*logging.Logger).Debugf("Posting measurement to InfluxDB, thing: %s, val: %s", thing.Name, value)

    // get thing org -> get influxdb assigned to org
    orgs := ctx.Value("orgs").(*Orgs)
    org, err := orgs.Get(ctx, thing.OrgId)
    if err != nil {
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("Going to post to InfluxDB %s as %s", org.InfluxDb, org.InfluxDbUsername)

    httpClient := ctx.Value("httpclient").(IHttpClient)

    if thing.Type != model.THING_TYPE_SENSOR {
        // ignore things which don't represent sensor
        return
    }

    // get thing name, use alias if set
    name := thing.Name
    if thing.Alias != "" {
        name = thing.Alias
    }

    body := fmt.Sprintf("sensor,id=%s,name=%s,class=%s value=%s", thing.Id.Hex(), name, thing.Sensor.Class, value)

    url, err := url.Parse(db.Uri)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Cannot decode InfluxDB url from %s (%s)", db.Uri, err.Error())
        return
    }

    url.Path = path.Join(url.Path, "write")

    params := url.Query()
    params.Add("db", org.InfluxDb)
    url.RawQuery = params.Encode()

    httpClient.PostString(ctx, url.String(), body, &db.Username, &db.Password)
}

func (db *InfluxDb) PostSwitchState(ctx context.Context, thing *model.Thing, value string) {
    ctx.Value("log").(*logging.Logger).Debugf("Posting switch state to InfluxDB, thing: %s, val: %s", thing.Name, value)

    // get thing org -> get influxdb assigned to org
    orgs := ctx.Value("orgs").(*Orgs)
    org, err := orgs.Get(ctx, thing.OrgId)
    if err != nil {
        return
    }

    ctx.Value("log").(*logging.Logger).Debugf("Going to post to InfluxDB %s as %s", org.InfluxDb, org.InfluxDbUsername)

    httpClient := ctx.Value("httpclient").(IHttpClient)

    if thing.Type != model.THING_TYPE_SWITCH {
        // ignore things which don't represent switch
        return
    }

    // get thing name, use alias if set
    name := thing.Name
    if thing.Alias != "" {
        name = thing.Alias
    }

    body := fmt.Sprintf("switch,id=%s,name=%s value=%s", thing.Id.Hex(), name, value)

    url, err := url.Parse(db.Uri)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Cannot decode InfluxDB url from %s (%s)", db.Uri, err.Error())
        return
    }

    url.Path = path.Join(url.Path, "write")

    params := url.Query()
    params.Add("db", org.InfluxDb)
    url.RawQuery = params.Encode()

    httpClient.PostString(ctx, url.String(), body, &db.Username, &db.Password)
}
