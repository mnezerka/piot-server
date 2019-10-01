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

    gqltesting.RunTest(t, &gqltesting.Test{
        Context: ctx,
        Schema:  graphql.MustParseSchema(schema.GetRootSchema(), &Resolver{}),
        Query: fmt.Sprintf(`
            {
                thing(id: "%s") {name}
            }
        `, thingId.Hex()),
        ExpectedResult: `
            {
                "thing": {
                    "name": "thing1"
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
                updateThing(thing: {id: "%s", name: "thing1new"}) { name}
            }
        `, id.Hex()),
        ExpectedResult: `
            {
                "updateThing": {
                    "name": "thing1new"
                }
            }
        `,
    })
}
