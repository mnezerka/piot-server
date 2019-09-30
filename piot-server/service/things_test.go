package service_test

import (
    "os"
    "testing"
    "piot-server/service"
    "piot-server/test"
    piotcontext "piot-server/context"
)

func TestFindUnknownThing(t *testing.T) {
    ctx := piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test")
    things := service.Things{}
    _, err := things.Find(ctx, "xx")
    test.Assert(t, err != nil, "Thing shall not be found")
}

func TestFindExistingThing(t *testing.T) {
    ctx := piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test")
    test.CleanDb(t, ctx)
    test.CreateThing(t, ctx, "thing1")
    things := service.Things{}
    _, err := things.Find(ctx, "thing1")
    test.Ok(t, err)
}

func TestRegisterThing(t *testing.T) {
    ctx := piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test")
    test.CleanDb(t, ctx)
    things := service.Things{}
    thing, err := things.Register(ctx, "thing1", "sensor")
    test.Ok(t, err)
    test.Assert(t, thing.Name == "thing1", "Wrong thing name")
    test.Assert(t, thing.Type == "sensor", "Wrong thing type")
}
