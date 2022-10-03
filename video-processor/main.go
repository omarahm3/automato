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
)

func main() {
	c, err := config.LoadConfig()
	check(err)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	db, err := database.Connect(c.DatabaseURI, ctx)
	check(err)
	defer cancel()

	var posts []types.Post
	err = db.FindAll(&posts)
	check(err)

	// err = db.RemoveAll()
	// check(err)

	videos := downloader.DownloadAll(posts)
	check(err)

	fmt.Println(videos)
}

func check(err error) {
	if err == nil {
		return
	}

	log.Fatalln("error occurred", err)
}
