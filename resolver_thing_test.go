package main_test

import (
    "context"
    "fmt"
    "testing"
    graphql "github.com/graph-gophers/graphql-go"
    "github.com/graph-gophers/graphql-go/gqltesting"
    "piot-server/schema"
)

func TestThingCreate(t *testing.T) {
    db := GetDb(t)
    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema: schema,
        Query: `
            mutation {
                createThing(name: "NewThing", type: "sensor") { name, type }
            }
        `,
        ExpectedResult: `
            {
                "createThing": {
                    "name": "NewThing",
                    "type": "sensor"
                }
            }
        `,
    })
}

func TestThingsGet(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    thingId := CreateThing(t, db, "thing1")
    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTests(t, []*gqltesting.Test{
        {
            Context: context.TODO(),
            Schema: schema,
            Query: `
                {
                    things { id, name }
                }
            `,
            ExpectedResult: fmt.Sprintf(`
                {
                    "things": [
                        {
                            "id": "%s",
                            "name": "thing1"
                        }
                    ]
                }
            `, thingId.Hex()),
        },
    })
}

func TestThingGet(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    thingId := CreateThing(t, db, "thing1")
    orgId := CreateOrg(t, db, "org1")
    AddOrgThing(t, db, orgId, "thing1")
    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema:  schema,
        Query: fmt.Sprintf(`
            {
                thing(id: "%s") {name, sensor {class, measurement_topic}}
            }
        `, thingId.Hex()),
        ExpectedResult: `
            {
                "thing": {
                    "name": "thing1", "sensor": {"class": "temperature", "measurement_topic": "value"}
                }
            }
        `,
    })
}

func TestThingUpdate(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    id := CreateThing(t, db, "thing1")
    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    t.Logf("Thing to be updated %s", id)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema: schema,
        Query: fmt.Sprintf(`
            mutation {
                updateThing(
                    thing: {
                        id: "%s",
                        name: "thing1new",
                        enabled: true,
                        last_seen_interval: 345,
                        availability_topic: "at",
                        telemetry_topic: "tt",
                        store_mysqldb: true,
                        store_mysqldb_interval: 60,
                        location_mqtt_topic: "ltopic",
                        location_mqtt_lat_value: "llatval",
                        location_mqtt_lng_value: "llngval",
                        location_mqtt_sat_value: "lsatval",
                        location_mqtt_ts_value: "ltsval",
                        location_lat: 34.555,
                        location_lng: 10.121212,
                        location_tracking: true
                    }
                ) {
                    name,
                    enabled,
                    last_seen_interval,
                    availability_topic,
                    telemetry_topic
                    store_mysqldb,
                    store_mysqldb_interval,
                    location_mqtt_topic,
                    location_mqtt_lat_value,
                    location_mqtt_lng_value,
                    location_mqtt_sat_value,
                    location_mqtt_ts_value,
                    location_lat,
                    location_lng,
                    location_sat,
                    location_ts,
                    location_tracking
                }
            }
        `, id.Hex()),
        ExpectedResult: `
            {
                "updateThing": {
                    "name": "thing1new",
                    "enabled": true,
                    "last_seen_interval": 345,
                    "availability_topic": "at",
                    "telemetry_topic": "tt",
                    "store_mysqldb": true,
                    "store_mysqldb_interval": 60,
                    "location_mqtt_topic": "ltopic",
                    "location_mqtt_lat_value": "llatval",
                    "location_mqtt_lng_value": "llngval",
                    "location_mqtt_sat_value": "lsatval",
                    "location_mqtt_ts_value": "ltsval",
                    "location_lat": 34.555,
                    "location_lng": 10.121212,
                    "location_sat": 0,
                    "location_ts": 0,
                    "location_tracking": true
                }
            }
        `,
    })
}

func TestThingSensorDataUpdate(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    id := CreateThing(t, db, "thing1")
    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    t.Logf("Thing to be updated %s", id)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema: schema,
        Query: fmt.Sprintf(`
            mutation {
                updateThingSensorData(data: {id: "%s", measurement_topic: "xyz"}) {sensor {measurement_topic}}
            }
        `, id.Hex()),
        ExpectedResult: `
            {
                "updateThingSensorData": {
                    "sensor": {"measurement_topic": "xyz"}
                }
            }
        `,
    })
}

func TestThingSwitchDataUpdate(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    id := CreateSwitch(t, db, "thing1")
    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    t.Logf("Thing to be updated %s", id)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema: schema,
        Query: fmt.Sprintf(`
            mutation {
                updateThingSwitchData(data: {id: "%s", state_topic: "statetopic"}) {switch {state_topic}}
            }
        `, id.Hex()),
        ExpectedResult: `
            {
                "updateThingSwitchData": {
                    "switch": {"state_topic": "statetopic"}
                }
            }
        `,
    })
}

func TestThingSetAlarm(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    id := CreateThing(t, db, "thing1")
    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    t.Logf("Thing to set alarm %s", id)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: context.TODO(),
        Schema: schema,
        Query: fmt.Sprintf(`
            mutation {
                setThingAlarm(id: "%s", active: true)
            }
        `, id.Hex()),
        ExpectedResult: `
            {
                "setThingAlarm": true
            }
        `,
    })
}
