package main

import (
	"context"
	"errors"

	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Orgs struct {
	log *logging.Logger
	db  *mongo.Database
}

func NewOrgs(log *logging.Logger, db *mongo.Database) *Orgs {
	return &Orgs{log: log, db: db}
}

func (t *Orgs) Get(id primitive.ObjectID) (*Org, error) {
	t.log.Debugf("Get org: %s", id.Hex())

	var org Org

	collection := t.db.Collection("orgs")
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&org)
	if err != nil {
		t.log.Errorf("Org service error : %v", err)
		return nil, err
	}

	return &org, nil
}

func (t *Orgs) GetByName(name string) (*Org, error) {
	t.log.Debugf("Finding org by name <%s>", name)

	var org Org

	// try to find thing in DB by its name
	err := t.db.Collection("orgs").FindOne(context.TODO(), bson.M{"name": name}).Decode(&org)
	if err != nil {
		return nil, errors.New("Org not found")
	}

	return &org, nil
}

func (t *Orgs) GetAll() ([]*Org, error) {
	ctx := context.TODO()

	var result []*Org

	// try to find thing in DB by its name
	cur, err := t.db.Collection("orgs").Find(ctx, bson.M{})
	if err != nil {
		t.log.Errorf("Orgs service error: %v", err)
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		// To decode into a struct, use cursor.Decode()
		org := Org{}
		err := cur.Decode(&org)
		if err != nil {
			t.log.Errorf("Orgs service error: %v", err)
			return nil, err
		}
		result = append(result, &org)
	}

	return result, nil
}
