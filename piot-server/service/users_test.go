package service_test

import (
    "testing"
    "piot-server/service"
    "piot-server/test"
)

func TestFindUserByNotExistingEmail(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    users := service.Users{}
    _, err := users.FindByEmail(ctx, "xx")
    test.Assert(t, err != nil, "User shall not be found")
}

func TestFindUserByExistingEmail(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    userId := test.CreateUser(t, ctx, "test1@test.com", "pass")
    orgId := test.CreateOrg(t, ctx, "testorg")
    test.AddOrgUser(t, ctx, orgId, userId)

    users := service.Users{}
    user, err := users.FindByEmail(ctx, "test1@test.com")
    test.Ok(t, err)
    test.Equals(t, "test1@test.com", user.Email)
    test.Equals(t, 1, len(user.Orgs))
    test.Equals(t, "testorg", user.Orgs[0].Name)
}
