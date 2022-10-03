package downloader

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/omarahm3/automato/types"
)

type Video struct {
	Path string
	Post types.Post
}

const ydl_command = "youtube-dl %s -q -o ./downloads/%s.%%(ext)s"

func DownloadAll(posts []types.Post) []Video {
	var (
		videos []Video
		wg     sync.WaitGroup
	)

	for _, p := range posts {
		go func(p types.Post) {
			wg.Add(1)
			downloadedPath, err := download(p.Video, p.Hash)
			if err != nil {
				log.Fatalf("error downloading video: %q of this post: %q::: %q", p.Video, p.Hash, err.Error())
			}

			log.Printf("Downloaded video: %q with hash %q on %q\n", p.Title, p.Hash, downloadedPath)
			videos = append(videos, Video{
				Post: p,
				Path: downloadedPath,
			})
			wg.Done()
		}(p)
	}

	wg.Wait()

	return videos
}

func download(u, output string) (string, error) {
	cmdString := fmt.Sprintf(ydl_command, u, output)
	args := strings.Split(cmdString, " ")
	log.Printf("running command: %q", strings.Join(args, ", "))
	cmd := exec.Command(args[0], args[1:]...)

	return args[len(args)-1], runCmd(cmd)
}

func runCmd(cmd *exec.Cmd) error {
	err := cmd.Start()
	if err != nil {
		log.Printf("run:: error occurred running command: %q", err.Error())
		return err
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("run:: error occurred running command: %q", err.Error())
		return err
	}

	return nil
}
