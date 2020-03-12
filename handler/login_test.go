package handler_test

import (
    "net/http"
    "testing"
    "github.com/mnezerka/go-piot/test"
)

func TestLoginSuccessful(t *testing.T) {
    db := test.GetDb(t)
    log := test.GetLogger(t)
    test.CleanDb(t, db)
    test.CreateUser(t, db, ADMIN_EMAIL, ADMIN_PASSWORD)
    LoginUser(t, log, db, ADMIN_EMAIL, ADMIN_PASSWORD, http.StatusOK)
}

func TestLoginWrongPassword(t *testing.T) {
    db := test.GetDb(t)
    log := test.GetLogger(t)
    test.CleanDb(t, db)
    test.CreateUser(t, db, ADMIN_EMAIL, ADMIN_PASSWORD)
    LoginUser(t, log, db, ADMIN_EMAIL, "xxx", 401)
}

func TestLoginWrongEmail(t *testing.T) {
    db := test.GetDb(t)
    log := test.GetLogger(t)
    test.CleanDb(t, db)
    test.CreateUser(t, db, ADMIN_EMAIL, ADMIN_PASSWORD)
    LoginUser(t, log, db, "xxx", ADMIN_PASSWORD, 401)
}

func TestLoginWrongEmailAndPassword(t *testing.T) {
    db := test.GetDb(t)
    log := test.GetLogger(t)
    test.CleanDb(t, db)
    test.CreateUser(t, db, ADMIN_EMAIL, ADMIN_PASSWORD)
    LoginUser(t, log, db, "xxx", "yyy", 401)
}

func TestLoginEmptyEmailAndPassword(t *testing.T) {
    db := test.GetDb(t)
    log := test.GetLogger(t)
    test.CleanDb(t, db)
    test.CreateUser(t, db, ADMIN_EMAIL, ADMIN_PASSWORD)
    LoginUser(t, log, db, "", "", 401)
}

