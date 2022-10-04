package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (d *Database) FindRandomOfToday(size int, results interface{}) error {
	t := time.Now()
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, t.Location())

	matchStage := bson.D{{
		Key: "$match",
		Value: bson.D{{
			Key: "created_at",
			Value: bson.M{
				"$gte": start,
				"$lt":  end,
			},
		}},
	}}
	sampleStage := bson.D{{
		Key: "$sample",
		Value: bson.M{
			"size": size,
		},
	}}
	cur, err := d.Collection.Aggregate(d.ctx, mongo.Pipeline{matchStage, sampleStage})
	if err != nil {
		return err
	}

	return cur.All(d.ctx, results)
}

func (d *Database) FindAll(results interface{}) error {
	cur, err := d.Collection.Find(d.ctx, bson.D{})
	if err != nil {
		return err
	}

	return cur.All(d.ctx, results)
}

func (d *Database) FindUnpublished(results interface{}, opts ...*options.FindOptions) error {
	cur, err := d.Collection.Find(d.ctx, bson.D{{
		Key:   "published",
		Value: false,
	}}, opts...)
	if err != nil {
		return err
	}

	return cur.All(d.ctx, results)
}

func (d *Database) MarkPublished(ids []primitive.ObjectID) error {
	_, err := d.Collection.UpdateMany(d.ctx, bson.D{{
		Key: "_id",
		Value: bson.D{{
			Key:   "$in",
			Value: ids,
		}},
	}}, bson.D{{
		Key: "$set",
		Value: bson.D{{
			Key:   "published",
			Value: true,
		}},
	}})

	if err != nil {
		return err
	}

	return err
}

func (d *Database) Find(filters, results interface{}, opts ...*options.FindOptions) error {
	cur, err := d.Collection.Find(d.ctx, filters, opts...)
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
