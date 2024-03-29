package main

import (
	"errors"

	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"

	//"go.mongodb.org/mongo-driver/bson"
	graphql "github.com/graph-gophers/graphql-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userProfileUpdateInput struct {
	OrgId *graphql.ID
}

/////////// User Profile Resolver

type UserProfileResolver struct {
	log   *logging.Logger
	db    *mongo.Database
	users *Users
	up    *UserProfile
}

func (r *UserProfileResolver) Email() string {
	return r.up.Email
}

func (r *UserProfileResolver) IsAdmin() bool {
	return r.up.IsAdmin
}

func (r *UserProfileResolver) OrgId() graphql.ID {
	return graphql.ID(r.up.OrgId.Hex())
}

func (r *UserProfileResolver) Orgs() []*OrgResolver {
	var result []*OrgResolver

	// get all orgs assigned to user
	orgs, err := r.users.FindUserOrgs(r.up.Id)
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

/////////// Resolver

// get active user profile
func (r *Resolver) UserProfile(ctx context.Context) (*UserProfileResolver, error) {

	// authorization checks
	profileValue := ctx.Value("profile")
	if profileValue == nil {
		r.log.Errorf("GQL: Missing user profile")
		return nil, errors.New("missing user profile")
	}
	profile := profileValue.(*UserProfile)
	r.log.Debugf("ctx %v", profile)

	return &UserProfileResolver{r.log, r.db, r.users, profile}, nil
}

func (r *Resolver) UpdateUserProfile(ctx context.Context, args struct{ Profile userProfileUpdateInput }) (*UserProfileResolver, error) {

	r.log.Debugf("Updating user profile")

	// get profile
	profileValue := ctx.Value("profile")
	if profileValue == nil {
		r.log.Errorf("GQL: Missing user profile")
		return nil, errors.New("missing user profile")
	}

	profile := profileValue.(*UserProfile)

	// try to find similar user matching new email
	if args.Profile.OrgId != nil {

		r.log.Debugf("Updating user profile active org to <%s>", string(*args.Profile.OrgId))

		// create ObjectID from string
		orgId, err := primitive.ObjectIDFromHex(string(*args.Profile.OrgId))
		if err != nil {
			return nil, err
		}

		// set active org
		err = r.users.SetActiveOrg(profile.Id, orgId)
		if err != nil {
			return nil, err
		}

		profile.OrgId = orgId
	}

	r.log.Debugf("User profile updated")
	return &UserProfileResolver{r.log, r.db, r.users, profile}, nil
}
