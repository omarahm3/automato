package main

import (
	"fmt"
	"os"

	"github.com/omarahm3/automato/uploader/config"
	"google.golang.org/api/youtube/v3"
)

func main() {
	c, err := config.LoadConfig()
	check(err)

	client, err := GetClient(youtube.YoutubeUploadScope, c)
	check(err)

	service, err := youtube.New(client)
	check(err)

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       "Top 10 best/worst Tiktoks today",
			Description: "Prepare yourself for the definitely not automated dose of Tiktoks",
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: c.PrivacyStatus,
		},
	}

	call := service.Videos.Insert([]string{"snippet", "status"}, upload)

	file, err := os.Open(c.OutputFile)
	check(err)
	defer file.Close()

	r, err := call.Media(file).Do()
	check(err)

	fmt.Printf("upload successful, video ID: %q\n", r.Id)
}

func check(err error) {
	if err == nil {
		return
	}

	panic(err)
}
