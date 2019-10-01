package handler_test

import (
    "net/http"
    "testing"
    "piot-server/test"
)

func TestLoginSuccessful(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    test.CreateUser(t, ctx, ADMIN_EMAIL, ADMIN_PASSWORD)
    Login(t, &ctx, ADMIN_EMAIL, ADMIN_PASSWORD, http.StatusOK)
}

func TestLoginWrongPassword(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    test.CreateUser(t, ctx, ADMIN_EMAIL, ADMIN_PASSWORD)
    Login(t, &ctx, ADMIN_EMAIL, "xxx", 401)
}

func TestLoginWrongEmail(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    test.CreateUser(t, ctx, ADMIN_EMAIL, ADMIN_PASSWORD)
    Login(t, &ctx, "xxx", ADMIN_PASSWORD, 401)
}

func TestLoginWrongEmailAndPassword(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    test.CreateUser(t, ctx, ADMIN_EMAIL, ADMIN_PASSWORD)
    Login(t, &ctx, "xxx", "yyy", 401)
}

func TestLoginEmptyEmailAndPassword(t *testing.T) {
    ctx := test.CreateTestContext()
    test.CleanDb(t, ctx)
    test.CreateUser(t, ctx, ADMIN_EMAIL, ADMIN_PASSWORD)
    Login(t, &ctx, "", "", 401)
}

