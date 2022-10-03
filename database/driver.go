package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	ctx        context.Context
	Collection *mongo.Collection
}

func (d *Database) InsertMany(data []interface{}) error {
	_, err := d.Collection.InsertMany(d.ctx, data, options.InsertMany().SetOrdered(false))
	return err
}

func (d *Database) FindAll(results interface{}) error {
	cur, err := d.Collection.Find(d.ctx, bson.D{})
	if err != nil {
		return err
	}

	return cur.All(d.ctx, results)
}

func (d *Database) RemoveAll() error {
	_, err := d.Collection.DeleteMany(d.ctx, bson.D{})
	return err
}

func (d *Database) EnsureIndex(key string) error {
	_, err := d.Collection.Indexes().CreateOne(d.ctx, mongo.IndexModel{
		Keys: bson.M{
			key: 1,
		},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func Connect(u string, ctx context.Context) (*Database, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(u).SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database("automato").Collection("posts")

	return &Database{
		Collection: collection,
		ctx:        ctx,
	}, nil
}
