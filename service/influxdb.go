package service

import (
    //"bytes"
    "context"
    //"net/http"
    //"net/url"
    "piot-server/model"
    "github.com/op/go-logging"
)

type IInfluxDb interface {
    PostMeasurement(ctx context.Context, thing *model.Thing, value string)
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

    /*

    //db.Uri influxurl + "/write?db=demo" + org.InfluxDB

    body := "home,room=livingroom temp=23,humidity=99"

    url, err := url.Parse(db.Uri)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Cannot decode InfluxDB url from %s (%s)", db.Uri, err.Error())
        return
    }
    q := url.Query()
    q.Add("db", org.InfluxDb)

    client := &http.Client{}

    req, err := http.NewRequest("POST", url.String(), bytes.NewReader(body))
    req.SetBasicAuth(db.Username, db.Password)
    res, err := client.Do(req)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Push to InfluxDB failed (%s)", err.Error())
        return
    }
    //robots, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    //if err != nil {
    //    log.Fatal(err)
    //}
    //fmt.Printf("%s", robots)
    */
}
