package test

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "path/filepath"
    "bytes"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "runtime"
    "reflect"
    "testing"
    "piot-server/resolver"
    "piot-server/handler"
    "piot-server/model"
    "piot-server/schema"
    graphql "github.com/graph-gophers/graphql-go"
)

type GqlResponseMessage struct {
    Message string `json:message`
}
type GqlResponse struct {
    Errors  []GqlResponseMessage `json:errors`
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

// helper function for checking and logging respone status
func CheckGqlResult(t *testing.T, rr *httptest.ResponseRecorder) {
    CheckStatusCode(t, rr, 200);
    //fmt.Print(rr.Body.String())

    var response GqlResponse
    if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
        t.Error(err)
    }

    if len(response.Errors) > 0 {
        fmt.Printf("%v", response)
        t.Errorf("\033[31mNot empty list of errors: %v\033[39m", response.Errors)
    }
}


func Body2Bytes(body *bytes.Buffer) ([]byte) {
    var result []byte
    result, _ = ioutil.ReadAll(body)
    return result
}

func Login(t *testing.T, ctx *context.Context, email string, password string, statusCode int) (string) {
    req, err := http.NewRequest("POST", "/login", strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)))
    Ok(t, err)

    rr := httptest.NewRecorder()

    handler := handler.AddContext(*ctx, handler.LoginHandler())
    handler.ServeHTTP(rr, req)

    CheckStatusCode(t, rr, statusCode)

    var response model.Token
    Ok(t, json.Unmarshal(Body2Bytes(rr.Body), &response))

    return response.Token
}


func GetAuthGqlRequest(t *testing.T, ctx *context.Context, email, password, body string) (*http.Request) {
    token := Login(t, ctx, email, password, 200)

    req, err := http.NewRequest("POST", "/any-path", strings.NewReader(body))
    Ok(t, err)
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

    return req
}

func GetGqlResponseRecorder(t *testing.T, ctx *context.Context, email, password, request string) (*httptest.ResponseRecorder) {

    req := GetAuthGqlRequest(t, ctx, email, password, request)

    rr := httptest.NewRecorder()

    graphqlSchema := graphql.MustParseSchema(schema.GetRootSchema(), &resolver.Resolver{})

    handler := handler.AddContext(*ctx, handler.Authorize(&handler.GraphQL{Schema: graphqlSchema}))

    handler.ServeHTTP(rr, req)

    return rr
}
