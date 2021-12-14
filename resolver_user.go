package main

import (
	"errors"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type userUpdateInput struct {
	Id       graphql.ID
	Email    *string
	Password *string
	IsAdmin  *bool
	OrgId    *graphql.ID
}

/////////// User Resolver

type UserResolver struct {
	log   *logging.Logger
	users *Users
	db    *mongo.Database
	u     *User
}

func (r *UserResolver) Id() graphql.ID {
	return graphql.ID(r.u.Id.Hex())
}

func (r *UserResolver) Email() string {
	return r.u.Email
}

func (r *UserResolver) Password() string {
	//maskedPassword := "********"
	//return maskedPassword
	return r.u.Password
}

func (r *UserResolver) Created() int32 {
	return r.u.Created
}

func (r *UserResolver) Orgs() []*OrgResolver {
	var result []*OrgResolver

	// get all orgs assigned to user
	orgs, err := r.users.FindUserOrgs(r.u.Id)
	if err != nil {
		r.log.Errorf("GQL: error : %v", err)
		return result
	}

	// convert orgs to org resolvers
	for i := 0; i < len(orgs); i++ {
		result = append(result, &OrgResolver{r.log, r.db, r.users, &orgs[i]})
	}

	return result
}

func (r *UserResolver) IsAdmin() bool {
	return r.u.IsAdmin
}

/////////// Resolver

// get user by email query
func (r *Resolver) User(args struct{ Id graphql.ID }) (*UserResolver, error) {

	r.log.Debugf("GQL: Fetch user: %v", args.Id)

	// create ObjectID from string
	id, err := primitive.ObjectIDFromHex(string(args.Id))
	if err != nil {
		r.log.Errorf("Graphql error : %v", err)
		return nil, errors.New("cannot decode ID")
	}

	user := User{}

	collection := r.db.Collection("users")
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		r.log.Errorf("Graphql error : %v", err)
		return nil, err
	}

	r.log.Debugf("GQL: Retrieved user %v", user)
	return &UserResolver{r.log, r.users, r.db, &user}, nil
}

// get users query
func (r *Resolver) Users(ctx context.Context) ([]*UserResolver, error) {

	collection := r.db.Collection("users")

	count, _ := collection.EstimatedDocumentCount(context.TODO())
	r.log.Debugf("GQL: Estimated users count %d", count)

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		r.log.Errorf("GQL: error : %v", err)
		return nil, err
	}
	defer cur.Close(context.TODO())

	var result []*UserResolver

	for cur.Next(context.TODO()) {
		// To decode into a struct, use cursor.Decode()
		user := User{}
		err := cur.Decode(&user)
		if err != nil {
			r.log.Errorf("GQL: error : %v", err)
			return nil, err
		}
		result = append(result, &UserResolver{r.log, r.users, r.db, &user})
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Resolver) CreateUser(ctx context.Context, args struct{ Email, Password string }) (*UserResolver, error) {

	r.log.Infof("Creating user %s", args.Email)

	user, err := r.users.Create(args.Email, args.Password)
	if err != nil {
		return nil, err
	}

	return &UserResolver{r.log, r.users, r.db, user}, nil
}

func (r *Resolver) UpdateUser(args struct{ User userUpdateInput }) (*UserResolver, error) {

	r.log.Debugf("Updating user %s", args.User.Id)

	// create ObjectID from string
	id, err := primitive.ObjectIDFromHex(string(args.User.Id))
	if err != nil {
		return nil, err
	}

	// try to find user to be updated
	var user User
	collection := r.db.Collection("users")
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, errors.New("User does not exist")
	}

	// try to find similar user matching new email
	if args.User.Email != nil {
		var similarUser User
		err := collection.FindOne(context.TODO(), bson.M{"$and": []bson.M{{"email": args.User.Email}, {"_id": bson.M{"$ne": id}}}}).Decode(&similarUser)
		if err == nil {
			return nil, errors.New("User of such name already exists")
		}
	}

	// user exists -> update it
	updateFields := bson.M{}
	if args.User.Email != nil {
		updateFields["email"] = args.User.Email
	}

	// if password was specified
	if args.User.Password != nil {
		if len(*args.User.Password) > 0 {
			// generate hash for given password (we don't store passwords in plain form)
			hash, err := GetPasswordHash(*args.User.Password)
			if err != nil {
				return nil, errors.New("error while hashing password, try again")
			}
			updateFields["password"] = hash
		}
	}

	if args.User.IsAdmin != nil {
		updateFields["is_admin"] = args.User.IsAdmin
	}
	update := bson.M{"$set": updateFields}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		r.log.Errorf("Updating user failed %v", err)
		return nil, errors.New("error while updating user")
	}

	// read user
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, errors.New("cannot fetch user data")
	}

	r.log.Debugf("User updated %v", user)
	return &UserResolver{r.log, r.users, r.db, &user}, nil
}
