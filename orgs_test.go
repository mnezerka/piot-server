package main_test

import (
    "testing"
    "piot-server"
)

func TestOrgsGetAll(t *testing.T) {
    db := GetDb(t)
    log := GetLogger(t)
    orgs:= main.NewOrgs(log, db)

    CleanDb(t, db)

    CreateOrg(t, db, "org1")
    CreateOrg(t, db, "org2")

    allOrgs, err := orgs.GetAll()
    Ok(t, err)
    Equals(t, 2, len(allOrgs))
    Equals(t, "org1", allOrgs[0].Name)
    Equals(t, "org2", allOrgs[1].Name)
}
