package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//CreateDBClient -creates the connection to the database
func CreateDBClient(ctx context.Context, connexionstring string, retry int32) (*mongo.Client, error) {
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(connexionstring))
	if err := conn.Ping(ctx, nil); err != nil {
		if retry >= 3 {
			return nil, err
		}
		retry++
		time.Sleep(time.Second * 2)
		return CreateDBClient(ctx, connexionstring, retry)
	}
	return conn, err
}
