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
    userId := CreateUser(t, db, "test@test.com", "passwd")
    thing2Id := CreateThing(t, db, "thing2")
    thing1Id := CreateThing(t, db, "thing1")
    thing3Id := CreateThing(t, db, "thing3")
    CreateThing(t, db, "thingX")  // this thing is created to test filtering based on active org
    orgId := CreateOrg(t, db, "org1")
    AddOrgThing(t, db, orgId, "thing3")
    AddOrgThing(t, db, orgId, "thing1")
    AddOrgThing(t, db, orgId, "thing2")

    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTests(t, []*gqltesting.Test{
        {
            Context: AuthContext(t, userId, orgId),
            Schema: schema,
            Query: `
                {
                    things(sort: {field: name, order: asc}) { id, name }
                }
            `,
            ExpectedResult: fmt.Sprintf(`
                {
                    "things": [
                        {
                            "id": "%s",
                            "name": "thing1"
                        },
                        {
                            "id": "%s",
                            "name": "thing2"
                        },
                        {
                            "id": "%s",
                            "name": "thing3"
                        }
                    ]
                }
            `, thing1Id.Hex(), thing2Id.Hex(), thing3Id.Hex()),
        },
    })
}

func TestThingsGetFiltered(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    userId := CreateUser(t, db, "test@test.com", "passwd")
    thing1Id := CreateThing(t, db, "thing1")
    CreateThing(t, db, "thing2")
    CreateThing(t, db, "somedevice")
    CreateThing(t, db, "thingX")  // this thing is created to test filtering based on active org
    orgId := CreateOrg(t, db, "org1")
    AddOrgThing(t, db, orgId, "thing1")
    AddOrgThing(t, db, orgId, "thing2")
    AddOrgThing(t, db, orgId, "somedevice")

    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    // filter by name
    gqltesting.RunTests(t, []*gqltesting.Test{
        {
            Context: AuthContext(t, userId, orgId),
            Schema: schema,
            Query: `
                {
                    things(filter: {name: "thing1"}) { id, name }
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
            `, thing1Id.Hex()),
        },
    })

    // filter by name 
    gqltesting.RunTests(t, []*gqltesting.Test{
        {
            Context: AuthContext(t, userId, orgId),
            Schema: schema,
            Query: `
                {
                    things(filter: {nameContains: "thing"}) { name }
                }
            `,
            ExpectedResult: `
                {
                    "things": [
                        {
                            "name": "thing1"
                        },
                        {
                            "name": "thing2"
                        }

                    ]
                }
            `,
        },
    })

}



func TestThingsGetAll(t *testing.T) {
    db := GetDb(t)
    CleanDb(t, db)
    userId := CreateAdmin(t, db, "admin@test.com", "passwd")
    thing2Id := CreateThing(t, db, "thing2")
    thing1Id := CreateThing(t, db, "thing1")
    thing3Id := CreateThing(t, db, "thing3")
    thingXId := CreateThing(t, db, "thingX")  // this thing is created to test filtering based on active org
    orgId := CreateOrg(t, db, "org1")
    AddOrgThing(t, db, orgId, "thing3")
    AddOrgThing(t, db, orgId, "thing1")
    AddOrgThing(t, db, orgId, "thing2")

    schema := graphql.MustParseSchema(schema.GetRootSchema(), getResolver(t, db))

    gqltesting.RunTests(t, []*gqltesting.Test{
        {
            Context: AuthContext(t, userId, orgId),
            Schema: schema,
            Query: `
                {
                    things(sort: {field: name, order: asc}, all: true) { id, name }
                }
            `,
            ExpectedResult: fmt.Sprintf(`
                {
                    "things": [
                        {
                            "id": "%s",
                            "name": "thing1"
                        },
                        {
                            "id": "%s",
                            "name": "thing2"
                        },
                        {
                            "id": "%s",
                            "name": "thing3"
                        },
                        {
                            "id": "%s",
                            "name": "thingX"
                        }

                    ]
                }
            `, thing1Id.Hex(), thing2Id.Hex(), thing3Id.Hex(), thingXId.Hex()),
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
