package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/omarahm3/automato/database"
	"github.com/omarahm3/automato/types"
	"github.com/omarahm3/video-processor/config"
	"github.com/omarahm3/video-processor/downloader"
	"github.com/omarahm3/video-processor/processor"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	c, err := config.LoadConfig()
	check(err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	db, err := database.Connect(c.DatabaseURI, ctx)
	check(err)
	defer cancel()

	var posts []types.Post
	err = db.FindUnpublished(&posts, options.Find().SetLimit(1))
	check(err)

	log.Printf("retrieved [%d] posts", len(posts))

	videos := downloader.DownloadAll(posts)

	// TODO Instead of removing videos, mark all as published
	// err = db.RemoveAll()
	// check(err)

	log.Printf("downloaded [%d] videos", len(videos))

	processedVideos := processor.ProcessAll(videos)
	fmt.Println(processedVideos)

	var ids []primitive.ObjectID
	for _, p := range processedVideos {
		ids = append(ids, p.Video.Post.ID)
	}
	err = db.MarkPublished(ids)
	check(err)
}

func check(err error) {
	if err == nil {
		return
	}

	log.Fatalln("error occurred", err)
}
