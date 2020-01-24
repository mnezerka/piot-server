package resolver

import (
    //"context"
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
                thing(id: "%s") {name, sensor {class, measurement_topic, store_influxdb}}
            }
        `, thingId.Hex()),
        ExpectedResult: `
            {
                "thing": {
                    "name": "thing1", "sensor": {"class": "temperature", "measurement_topic": "org1/thing1/value", "store_influxdb": true}
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
                updateThing(thing: {id: "%s", name: "thing1new", enabled: true}) {name, enabled}
            }
        `, id.Hex()),
        ExpectedResult: `
            {
                "updateThing": {
                    "name": "thing1new",
                    "enabled": true
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
                updateThingSensorData(data: {id: "%s", store_influxdb: true}) {sensor {store_influxdb}}
            }
        `, id.Hex()),
        ExpectedResult: `
            {
                "updateThingSensorData": {
                    "sensor": {"store_influxdb": true}
                }
            }
        `,
    })
}
