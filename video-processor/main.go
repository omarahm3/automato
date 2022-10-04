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

// Time output: --- WTF!
// ________________________________________________________
// Executed in   24.86 mins    fish           external
//    usr time  131.57 mins    0.00 micros  131.57 mins
//    sys time    1.89 mins  498.00 micros    1.89 mins

// TODO script is taking so much time to download & process all videos that makes it fail at the end with context deadline exceeded
// Maybe split 3 functionalities into 3 different packages? (downloader, processor, merger)
const database_timeout = 1000 * 60 * time.Second

func main() {
	c, err := config.LoadConfig()
	check(err)

	err = ensureBaseDir(c.BaseDir)
	check(err)

	ctx, cancel := context.WithTimeout(context.Background(), database_timeout)
	db, err := database.Connect(c.DatabaseURI, ctx)
	check(err)
	defer cancel()

	posts := getPosts(db, c.PostsLimit)
	log.Printf("retrieved [%d] posts", len(posts))

	videos := downloader.DownloadAll(posts, c.BaseDir, c.Threads)
	log.Printf("downloaded [%d] videos", len(videos))

	processedVideos := processor.ProcessAll(videos, c.BaseDir)
	log.Printf("processed [%d] videos", len(processedVideos))

	err = processor.MergeAll(processedVideos, c.BaseDir, c.OutputPath)
	check(err)
	log.Println("merged videos into a single video")

	markPostsPublished(db, processedVideos)

	err = clean(c)
	check(err)
}

func clean(c *config.Config) error {
	// remove downloads dir
	d := path.Join(c.BaseDir, "downloads")
	err := os.RemoveAll(d)
	if err != nil {
		return err
	}

	b := path.Join(c.BaseDir, "blurry")
	return os.RemoveAll(b)
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

	log.Println("error occurred")
	panic(err)
}
