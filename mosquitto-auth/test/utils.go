package test

import (
    "context"
    "fmt"
    "path/filepath"
    "net/http/httptest"
    "os"
    "runtime"
    "reflect"
    "testing"
    piotcontext "mosquitto-auth/context"
)

func CreateTestContext() context.Context {
    contextOptions := piotcontext.NewContextOptions()
    contextOptions.DbUri = os.Getenv("MONGODB_URI")
    contextOptions.DbName = "piot-test"
    contextOptions.LogLevel = "DEBUG"
    ctx := piotcontext.NewContext(contextOptions)

    //mqtt := &MqttMock{}
    //ctx := piotcontext.NewContext(os.Getenv("MONGODB_URI"), "piot-test", mqtt, "DEBUG")
    callerEmail := "caller@test.com"
    ctx = context.WithValue(ctx, "user_email", &callerEmail)
    ctx = context.WithValue(ctx, "is_authorized", true)
    return ctx
}

// assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
    if !condition {
        _, file, line, _ := runtime.Caller(1)
        fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
        tb.FailNow()
    }
}

// ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
    if err != nil {
        _, file, line, _ := runtime.Caller(1)
        fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
        tb.FailNow()
    }
}

// equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
    if !reflect.DeepEqual(exp, act) {
        _, file, line, _ := runtime.Caller(1)
        fmt.Printf("\033[31m%s:%d:\n\texp: %#v\n\tgot: %#v\033[39m\n", filepath.Base(file), line, exp, act)
        tb.FailNow()
    }
}

// helper function for checking and logging respone status
func CheckStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
    if status := rr.Code; status != expected {
        t.Errorf("\033[31mWrong response status code: got %v want %v, body:\n%s\033[39m",
            status, expected, rr.Body.String())
    }
}
