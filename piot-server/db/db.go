package db

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    //"github.com/mongodb/mongo-go-driver/mongo"
    //"github.com/mongodb/mongo-go-driver/mongo/options"
)

func GetDB(uri string) (*mongo.Database, error) {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }

    // Check the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        return nil, err
    }

    return client.Database("piot"), nil
}


func GetDBCollection() (*mongo.Collection, error) {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        return nil, err
    }
    // Check the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        return nil, err
    }
    collection := client.Database("piot").Collection("users")
    return collection, nil
}
