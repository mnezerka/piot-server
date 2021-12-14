package main

import (
	"errors"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type orgUpdateInput struct {
	Id               graphql.ID
	Name             *string
	Description      *string
	InfluxDb         *string
	InfluxDbUsername *string
	InfluxDbPassword *string
	MysqlDb          *string
	MysqlDbUsername  *string
	MysqlDbPassword  *string
	MqttUsername     *string
	MqttPassword     *string
}

/////////// Org Resolver

type OrgResolver struct {
	log   *logging.Logger
	db    *mongo.Database
	users *Users
	org   *Org
}

func (r *OrgResolver) Id() graphql.ID {
	return graphql.ID(r.org.Id.Hex())
}

func (r *OrgResolver) Name() string {
	return r.org.Name
}

func (r *OrgResolver) Description() string {
	return r.org.Description
}

func (r *OrgResolver) InfluxDb() string {
	return r.org.InfluxDb
}

func (r *OrgResolver) InfluxdbUsername() string {
	return r.org.InfluxDbUsername
}

func (r *OrgResolver) InfluxdbPassword() string {
	return r.org.InfluxDbPassword
}

func (r *OrgResolver) MysqlDb() string {
	return r.org.MysqlDb
}

func (r *OrgResolver) MysqldbUsername() string {
	return r.org.MysqlDbUsername
}

func (r *OrgResolver) MysqldbPassword() string {
	return r.org.MysqlDbPassword
}

func (r *OrgResolver) MqttUsername() string {
	return r.org.MqttUsername
}

func (r *OrgResolver) MqttPassword() string {
	return r.org.MqttPassword
}

func (r *OrgResolver) Created() int32 {
	return r.org.Created
}

// select all users that are assigned to current org
func (r *OrgResolver) Users() []*UserResolver {

	var result []*UserResolver

	r.log.Debugf("GQL: Fetching users for org: %s", r.org.Id.Hex())

	collection := r.db.Collection("orgusers")

	// filter orusers to current (single) org
	stage_match := bson.M{"$match": bson.M{"org_id": r.org.Id}}

	// find assignments to orgs
	stage_lookup := bson.M{"$lookup": bson.M{"from": "users", "localField": "user_id", "foreignField": "_id", "as": "users"}}

	// unwind users
	stage_unwind := bson.M{"$unwind": "$users"}

	// replace root
	stage_new_root := bson.M{"$replaceWith": "$users"}

	pipeline := []bson.M{stage_match, stage_lookup, stage_unwind, stage_new_root}

	cur, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		r.log.Errorf("GQL: error : %v", err)
		return result
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var user User
		if err := cur.Decode(&user); err != nil {
			r.log.Errorf("GQL: error : %v", err)
			return result
		}
		result = append(result, &UserResolver{r.log, r.users, r.db, &user})
	}

	if err := cur.Err(); err != nil {
		r.log.Errorf("GQL: error during cursor processing: %v", err)
		return result
	}

	return result
}

/////////// Resolver

func (r *Resolver) Org(args struct{ Id graphql.ID }) (*OrgResolver, error) {

	org := Org{}

	r.log.Debugf("GQL: Fetching org %v", args.Id)

	// create ObjectID from string
	id, err := primitive.ObjectIDFromHex(string(args.Id))
	if err != nil {
		r.log.Errorf("Graphql error : %v", err)
		return nil, errors.New("cannot decode ID")
	}

	collection := r.db.Collection("orgs")
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&org)

	if err != nil {
		r.log.Errorf("Graphql error : %v", err)
		return nil, err
	}

	return &OrgResolver{r.log, r.db, r.users, &org}, nil
}

func (r *Resolver) Orgs() ([]*OrgResolver, error) {

	collection := r.db.Collection("orgs")

	count, _ := collection.EstimatedDocumentCount(context.TODO())
	r.log.Debugf("GQL: Estimated orgs count %d", count)

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		r.log.Errorf("GQL: error : %v", err)
		return nil, err
	}
	defer cur.Close(context.TODO())

	var result []*OrgResolver

	for cur.Next(context.TODO()) {
		// To decode into a struct, use cursor.Decode()
		var org Org
		if err := cur.Decode(&org); err != nil {
			r.log.Debugf("GQL: After decode %v", err)
			r.log.Errorf("GQL: error : %v", err)
			return nil, err
		}
		result = append(result, &OrgResolver{r.log, r.db, r.users, &org})
	}

	r.log.Debug("Have orgs")

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Resolver) CreateOrg(args *struct {
	Name        string
	Description string
}) (*OrgResolver, error) {

	org := &Org{
		Name:        args.Name,
		Description: args.Description,
		Created:     int32(time.Now().Unix()),
	}

	r.log.Infof("Creating org %s", args.Name)

	// try to find existing user
	var orgExisting Org
	collection := r.db.Collection("orgs")
	err := collection.FindOne(context.TODO(), bson.M{"name": args.Name}).Decode(&orgExisting)
	if err == nil {
		return nil, errors.New("organization of such name already exists")
	}

	// org does not exist -> create new one
	_, err = collection.InsertOne(context.TODO(), org)
	if err != nil {
		return nil, errors.New("error while creating organizaton")
	}

	r.log.Debugf("Created organization: %v", *org)

	return &OrgResolver{r.log, r.db, r.users, org}, nil
}

func (r *Resolver) UpdateOrg(args struct{ Org orgUpdateInput }) (*OrgResolver, error) {

	r.log.Debugf("Updating org %ss", args.Org.Id)

	// create ObjectID from string
	id, err := primitive.ObjectIDFromHex(string(args.Org.Id))
	if err != nil {
		return nil, err
	}

	// try to find org to be updated
	var org Org
	collection := r.db.Collection("orgs")
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&org)
	if err != nil {
		return nil, errors.New("Org does not exist")
	}

	// try to find similar org matching new name
	if args.Org.Name != nil {
		var similarOrg Org
		collection := r.db.Collection("orgs")
		err := collection.FindOne(context.TODO(), bson.M{"$and": []bson.M{{"name": args.Org.Name}, {"_id": bson.M{"$ne": id}}}}).Decode(&similarOrg)
		if err == nil {
			return nil, errors.New("Org of such name already exists")
		}
	}

	// org exists -> update it
	updateFields := bson.M{}
	if args.Org.Name != nil {
		updateFields["name"] = args.Org.Name
	}
	if args.Org.Description != nil {
		updateFields["description"] = args.Org.Description
	}
	if args.Org.InfluxDb != nil {
		updateFields["influxdb"] = args.Org.InfluxDb
	}
	if args.Org.InfluxDbUsername != nil {
		updateFields["influxdb_username"] = args.Org.InfluxDbUsername
	}
	if args.Org.InfluxDbPassword != nil {
		updateFields["influxdb_password"] = args.Org.InfluxDbPassword
	}
	if args.Org.MysqlDb != nil {
		updateFields["mysqldb"] = args.Org.MysqlDb
	}
	if args.Org.MysqlDbUsername != nil {
		updateFields["mysqldb_username"] = args.Org.MysqlDbUsername
	}
	if args.Org.MysqlDbPassword != nil {
		updateFields["mysqldb_password"] = args.Org.MysqlDbPassword
	}
	if args.Org.MqttUsername != nil {
		updateFields["mqtt_username"] = args.Org.MqttUsername
	}
	if args.Org.MqttPassword != nil {
		updateFields["mqtt_password"] = args.Org.MqttPassword
	}

	update := bson.M{"$set": updateFields}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		r.log.Errorf("Updating org failed %v", err)
		return nil, errors.New("error while updating org")
	}

	// read org
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&org)
	if err != nil {
		return nil, errors.New("cannot fetch org data")
	}

	r.log.Debugf("Org updated %v", org)
	return &OrgResolver{r.log, r.db, r.users, &org}, nil
}

func (r *Resolver) AddOrgUser(args *struct {
	OrgId  graphql.ID
	UserId graphql.ID
}) (*bool, error) {

	r.log.Debugf("Adding user %s to org %s", args.UserId, args.OrgId)

	// create ObjectIDs from string
	orgId, err := primitive.ObjectIDFromHex(string(args.OrgId))
	if err != nil {
		return nil, err
	}
	userId, err := primitive.ObjectIDFromHex(string(args.UserId))
	if err != nil {
		return nil, err
	}

	// try to find existing assignment
	var similarOrgUser OrgUser
	collection := r.db.Collection("orgusers")
	err = collection.FindOne(context.TODO(), bson.M{"$and": []bson.M{{"user_id": userId}, {"org_id": orgId}}}).Decode(&similarOrgUser)
	if err == nil {
		return nil, errors.New("User is allready assigned to given organization")
	}

	// assignment does not exist -> create new one
	orgUser := &OrgUser{
		UserId:  userId,
		OrgId:   orgId,
		Created: int32(time.Now().Unix()),
	}
	_, err = collection.InsertOne(context.TODO(), orgUser)
	if err != nil {
		return nil, errors.New("error while adding user to organization")
	}

	r.log.Debugf("User %s added to Org %s", userId, orgId)
	return nil, nil
}

func (r *Resolver) RemoveOrgUser(args *struct {
	OrgId  graphql.ID
	UserId graphql.ID
}) (*bool, error) {

	r.log.Debugf("Removing user %s from org %s", args.UserId, args.OrgId)

	// create ObjectIDs from string
	orgId, err := primitive.ObjectIDFromHex(string(args.OrgId))
	if err != nil {
		return nil, err
	}
	userId, err := primitive.ObjectIDFromHex(string(args.UserId))
	if err != nil {
		return nil, err
	}

	collection := r.db.Collection("orgusers")
	_, err = collection.DeleteOne(context.TODO(), bson.M{"$and": []bson.M{{"user_id": userId}, {"org_id": orgId}}})
	if err != nil {
		r.log.Errorf("Cannot remove user %s from org %s (%v)", args.UserId, args.OrgId, err)
		return nil, errors.New("remove user from organization failed")
	}

	r.log.Debugf("User %s removed from  org %s", args.UserId, args.OrgId)
	return nil, nil
}
