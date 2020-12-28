package main_test

import (
    "testing"
    "piot-server"
)

func TestFindUserByNotExistingEmail(t *testing.T) {
    db := GetDb(t)
    log := GetLogger(t)
    users := main.NewUsers(log, db)

    CleanDb(t, db)
    _, err := users.FindByEmail("xx")
    Assert(t, err != nil, "User shall not be found")
}

func TestFindUserByExistingEmail(t *testing.T) {
    db := GetDb(t)
    log := GetLogger(t)
    users := main.NewUsers(log, db)

    CleanDb(t, db)
    userId := CreateUser(t, db, "test1@com", "pass")
    orgId := CreateOrg(t, db, "testorg")
    AddOrgUser(t, db, orgId, userId)

    user, err := users.FindByEmail("test1@com")
    Ok(t, err)
    Equals(t, "test1@com", user.Email)
    Equals(t, 1, len(user.Orgs))
    Equals(t, "testorg", user.Orgs[0].Name)
}

func TestGetAdmins(t *testing.T) {
    db := GetDb(t)
    log := GetLogger(t)
    users := main.NewUsers(log, db)

    CleanDb(t, db)
    CreateUser(t, db, "test1@com", "pass")
    CreateUser(t, db, "test2@com", "pass2")
    CreateAdmin(t, db, "admin1@com", "admpass1")
    CreateAdmin(t, db, "admin2@com", "admpass2")

    admins, err := users.GetAdmins()
    Ok(t, err)
    Equals(t, 2, len(admins))
    Equals(t, "admin1@com", admins[0].Email)
    Equals(t, "admin2@com", admins[1].Email)
}

func TestCreateUser(t *testing.T) {
    db := GetDb(t)
    log := GetLogger(t)
    users := main.NewUsers(log, db)

    CleanDb(t, db)
    user, err := users.Create("test1@test.com", "pass")
    Ok(t, err)

    //compare ids

    userRead, err := users.FindByEmail("test1@test.com")
    Ok(t, err)
    Equals(t, user.Id, userRead.Id)
    Equals(t, "test1@test.com", userRead.Email)
    Equals(t, 0, len(userRead.Orgs))
}
