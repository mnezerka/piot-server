package utils

import (
    "golang.org/x/crypto/bcrypt"
)

func GetPasswordHash(pwd string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 5)
    return string(hash), err
}

