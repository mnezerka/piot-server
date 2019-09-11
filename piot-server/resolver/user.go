package resolver

import (
    //graphql "github.com/graph-gophers/graphql-go"
    //"strconv"
    "piot-server/model"
    "piot-server/service"
    "errors"
    "github.com/op/go-logging"
    "golang.org/x/net/context"
)

type UserResolver struct {
    u *model.User
}

func (r *UserResolver) Email() string {
    return r.u.Email
}

func (r *UserResolver) Password() *string {
    maskedPassword := "********"
    return &maskedPassword
}

// get user by email query
func (r *Resolver) User(ctx context.Context, args struct {Email string}) (*UserResolver, error) {

    userId := ctx.Value("user_id").(*int64)

    user, err := ctx.Value("userService").(*service.UserService).FindByEmail(args.Email)

    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }
    ctx.Value("log").(*logging.Logger).Debugf("Retrieved user by user_id[%d] : %v", *userId, *user)
    return &UserResolver{user}, nil
}

// get users query
func (r *Resolver) Users(ctx context.Context) ([]*UserResolver, error) {

    ctx.Value("log").(*logging.Logger).Debugf("here")

    if isAuthorized := ctx.Value("is_authorized").(bool); !isAuthorized {
        return nil, errors.New(configuration.CredentialsError)
    }

    ctx.Value("log").(*logging.Logger).Debugf("here2")

    userId := ctx.Value("user_id").(*int64)

    ctx.Value("log").(*logging.Logger).Debugf("here3")

    users, err := ctx.Value("userService").(*service.UserService).List()
    //count, err := ctx.Value("userService").(*UserService).Count()

    ctx.Value("log").(*logging.Logger).Debugf("here4")


    ctx.Value("log").(*logging.Logger).Debugf("Retrieved users by user_id[%d] :", *userId)

    config := ctx.Value("config").(*configuration.Config)

    if config.DebugMode {
        for _, user := range users {
            ctx.Value("log").(*logging.Logger).Debugf("%v", *user)
        }
    }

    ctx.Value("log").(*logging.Logger).Debugf("Retrieved total users count by user_id[%d] : %v", *userId, len(users))

    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
        return nil, err
    }

    var result []*UserResolver
    for _, v := range users {
        result = append(result, &UserResolver{v})
    }

    //return &UserResolver{users: users, totalCount: count, from: &(users[0].ID), to: &(users[len(users)-1].ID)}, nil

    return result, nil
}
