package main_test

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "testing"
    "time"
    "piot-server"
)

func setLastSeenAttributes(t *testing.T, id primitive.ObjectID, lastSeen, lastSeenInterval int32) {

    _, err := db.Collection("things").UpdateOne(context.TODO(), bson.M{
        "_id": id},
        bson.M{
            "$set": bson.M{"last_seen": lastSeen, "last_seen_interval": lastSeenInterval},
        })
    Ok(t, err)
}

func TestMonitorCheckClear(t *testing.T) {
    const THING = "device1"
    const THING2 = "device2"
    const THING3 = "device3"
    const THING4 = "device4"
    const ORG = "org1"

    log := GetLogger(t)
    db := GetDb(t)
    things := GetThings(t, log, db)
    cfg := GetConfig()
    users := GetUsers(t, log, db)

    now := int32(time.Now().Unix())

    CleanDb(t, db)

    // we need several users to verify email recepients
    CreateUser(t, db, "test2@com", "pass2")
    CreateAdmin(t, db, "admin1@com", "admpass1")
    CreateAdmin(t, db, "admin2@com", "admpass2")

    // new org
    orgId := CreateOrg(t, db, ORG)

    // create valid thing (no last seen interval set)
    thingId := CreateDevice(t, db, THING)
    setLastSeenAttributes(t, thingId, 456, 0)
    AddOrgThing(t, db, orgId, THING)

    // create valid thing (last seen fits to last seen interval)
    thing2Id := CreateDevice(t, db, THING2)
    setLastSeenAttributes(t, thing2Id, now - 50, 56)
    AddOrgThing(t, db, orgId, THING2)

    // create not valid thing
    thing3Id := CreateDevice(t, db, THING3)
    setLastSeenAttributes(t, thing3Id, now - 57, 56)
    AddOrgThing(t, db, orgId, THING3)

    // create second not valid thing
    thing4Id := CreateDevice(t, db, THING4)
    setLastSeenAttributes(t, thing4Id, 23, 56)
    AddOrgThing(t, db, orgId, THING4)

    mockMail := GetMockMailClient(t, log);

    monitor := main.NewMonitor(log, db, mockMail, things, cfg, users)

    monitor.Check()

    Equals(t, 1, len(mockMail.Calls))
    Equals(t, mockMail.Calls[0].Subject, "[piot][alarm] Not Available Devices")
    Equals(t, mockMail.Calls[0].From, cfg.MailFrom)
    Equals(t, mockMail.Calls[0].To, []string{"admin1@com", "admin2@com"})
    Contains(t, mockMail.Calls[0].Message, "Following")
    Contains(t, mockMail.Calls[0].Message, "device3")
    Contains(t, mockMail.Calls[0].Message, "device4")
}
