package main

import (
	"context"
	"log"
	"time"

	"github.com/omarahm3/reddit-aggregator/api"
	"github.com/omarahm3/reddit-aggregator/config"
	"github.com/omarahm3/reddit-aggregator/database"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	PostsLimit int
	PostsType  string
}

func main() {
	c, err := config.LoadConfig()
	check(err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	db, err := database.Connect(c.DatabaseURI, ctx)
	check(err)
	defer cancel()

	err = db.EnsureIndex("hash")
	check(err)

	posts, err := api.Fetch(c)
	check(err)

	var ps []interface{}
	for _, p := range posts {
		ps = append(ps, p)
	}

	err = db.InsertMany(ps)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		check(err)
	}
}

func check(err error) {
	if err == nil {
		return
	}

	log.Fatalln("error occurred", err)
}
