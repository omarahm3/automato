package main

import (
	"context"
	"log"
	"os"
	"path"
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

	err = ensureBaseDir(c.BaseDir)
	check(err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	db, err := database.Connect(c.DatabaseURI, ctx)
	check(err)
	defer cancel()

	posts := getPosts(db, c.PostsLimit)
	log.Printf("retrieved [%d] posts", len(posts))

	videos := downloader.DownloadAll(posts, c.BaseDir)
	log.Printf("downloaded [%d] videos", len(videos))

	processedVideos := processor.ProcessAll(videos, c.BaseDir)
	log.Printf("processed [%d] videos", len(processedVideos))

	markPostsPublished(db, processedVideos)
}

func ensureBaseDir(b string) error {
	downloadsDir := path.Join(b, "downloads")
	err := os.MkdirAll(downloadsDir, os.ModePerm)
	if err != nil {
		return err
	}

	blurryDir := path.Join(b, "blurry")
	return os.MkdirAll(blurryDir, os.ModePerm)
}

func markPostsPublished(db *database.Database, videos []processor.ProcessedVideo) {
	var ids []primitive.ObjectID
	for _, p := range videos {
		ids = append(ids, p.Video.Post.ID)
	}
	err := db.MarkPublished(ids)
	check(err)
}

func getPosts(db *database.Database, limit int) []types.Post {
	var posts []types.Post
	err := db.FindUnpublished(&posts, options.Find().SetLimit(int64(limit)))
	check(err)
	return posts
}

func check(err error) {
	if err == nil {
		return
	}

	log.Fatalln("error occurred", err)
}
