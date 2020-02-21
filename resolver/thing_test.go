package resolver

import (
    "fmt"
    "testing"
    graphql "github.com/graph-gophers/graphql-go"
    "github.com/graph-gophers/graphql-go/gqltesting"
    "piot-server/schema"
    "piot-server/test"
)

func TestThingCreate(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
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
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    thingId := test.CreateThing(t, ctx, "thing1")

    gqltesting.RunTests(t, []*gqltesting.Test{
        {
            Context: ctx,
            Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
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
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    thingId := test.CreateThing(t, ctx, "thing1")
    orgId := test.CreateOrg(t, ctx, "org1")
    test.AddOrgThing(t, ctx, orgId, "thing1")

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
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
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    id := test.CreateThing(t, ctx, "thing1")

    t.Logf("Thing to be updated %s", id)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
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
                        store_mysqldb_interval: 60
                    }
                ) {
                    name,
                    enabled,
                    last_seen_interval,
                    availability_topic,
                    telemetry_topic
                    store_mysqldb,
                    store_mysqldb_interval
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
                    "store_mysqldb_interval": 60
                }
            }
        `,
    })
}

func TestThingSensorDataUpdate(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    id := test.CreateThing(t, ctx, "thing1")

    t.Logf("Thing to be updated %s", id)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
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
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    id := test.CreateSwitch(t, ctx, "thing1")

    t.Logf("Thing to be updated %s", id)

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
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
