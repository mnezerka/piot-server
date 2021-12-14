package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"piot-server/utils"
	"time"

	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegistrationHandler struct {
	log *logging.Logger
	db  *mongo.Database
}

func NewRegistrationHandler(log *logging.Logger, db *mongo.Database) *RegistrationHandler {
	return &RegistrationHandler{log: log, db: db}
}

func (h *RegistrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// check http method, POST is required
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, errors.New("only POST method is allowed"), http.StatusMethodNotAllowed)
		return
	}

	// decode json from request body
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		WriteErrorResponse(w, err, 400)
		return
	}

	// check required attributes
	if len(credentials.Email) == 0 {
		WriteErrorResponse(w, errors.New("email field is empty or not specified"), 400)
		return
	}
	if len(credentials.Password) == 0 {
		WriteErrorResponse(w, errors.New("password field is empty or not specified"), 400)
		return
	}
	if !ValidateEmail(credentials.Email) {
		WriteErrorResponse(w, errors.New("email field has wrong format"), 400)
		return
	}

	// try to find existing user
	var user User
	collection := h.db.Collection("users")
	err = collection.FindOne(context.TODO(), bson.M{"email": credentials.Email}).Decode(&user)
	if err == nil {
		WriteErrorResponse(w, errors.New("user identified by this email already exists"), 409)
		return
	}

	// generate hash for given password (we don't store passwords in plain form)
	hash, err := utils.GetPasswordHash(credentials.Password)
	if err != nil {
		WriteErrorResponse(w, errors.New("error while hashing password, try again"), 500)
		return
	}
	user.Email = credentials.Email
	user.Password = hash
	user.Created = int32(time.Now().Unix())

	// user does not exist -> create new one
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		WriteErrorResponse(w, errors.New("User while creating user, try again"), 500)
		return
	}

	var response ResponseResult
	response.Result = "Registration successful"

	h.log.Debugf("User is registered: %s", user.Email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
