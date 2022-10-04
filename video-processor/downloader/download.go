package downloader

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/omarahm3/automato/types"
	"github.com/omarahm3/video-processor/helpers"
)

type Video struct {
	Path string
	Post types.Post
}

const ydl_command = "youtube-dl %s -q -o %s.%%(ext)s"

func DownloadAll(posts []types.Post, base string) []Video {
	var (
		videos []Video
		wg     sync.WaitGroup
	)

	wg.Add(len(posts))

	for _, p := range posts {
		log.Printf("downloading video of post %q", p.Hash)
		go func(p types.Post) {
			log.Printf("calling download %q", p.Hash)
			downloadedPath, out, err := download(p.Video, base, p.Hash)
			if err != nil {
				log.Fatalf("error downloading video: %q of this post: %q::: %q\ncommand output: %s", p.Video, p.Hash, err.Error(), out)
			}

			log.Printf("Downloaded video: %q with hash %q on %q\n", p.Title, p.Hash, downloadedPath)
			videos = append(videos, Video{
				Post: p,
				Path: strings.ReplaceAll(downloadedPath, "%(ext)s", "mp4"),
			})
			wg.Done()
		}(p)
	}

	wg.Wait()

	return videos
}

func download(u, base, output string) (string, string, error) {
	o := path.Join(base, "downloads", output)
	cmdString := fmt.Sprintf(ydl_command, u, o)
	args := strings.Split(cmdString, " ")
	log.Printf("running command: %q", strings.Join(args, ", "))
	cmd := exec.Command(args[0], args[1:]...)
	out, err := helpers.RunCmd(cmd)

	return args[len(args)-1], out, err
}
