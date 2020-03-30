package main

import (
    "fmt"
    "strconv"
    "reflect"
    "golang.org/x/crypto/bcrypt"
)

func GetPasswordHash(pwd string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 5)
    return string(hash), err
}

func PrimitiveToString(value interface{}) (string, error) {
    var output string

    // Values may come in by pointer for optionals, so make sure to dereference.
    v := reflect.Indirect(reflect.ValueOf(value))
    t := v.Type()
    kind := t.Kind()

    switch kind {
    case reflect.Int8, reflect.Int32, reflect.Int64, reflect.Int:
        output = strconv.FormatInt(v.Int(), 10)
    case reflect.Float32, reflect.Float64:
        output = strconv.FormatFloat(v.Float(), 'f', -1, 64)
    case reflect.Bool:
        if v.Bool() {
            output = "true"
        } else {
            output = "false"
        }
    case reflect.String:
        output = v.String()
    default:
        return "", fmt.Errorf("unsupported primitive type %s", reflect.TypeOf(value).String())
    }
    return output, nil
}
