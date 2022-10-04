package processor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/omarahm3/video-processor/downloader"
	"github.com/omarahm3/video-processor/helpers"
)

type ProcessedVideo struct {
	Path  string
	Video downloader.Video
}

const (
	ffmpeg_merge_command = "ffmpeg -f concat -safe 0 -i %s %s"
	ffmpeg_quality       = "1080k"
	ffmpeg_command       = `ffmpeg -i %s -lavfi %s -vb %s -c:v libx264 -crf 20 %s.mp4 -n`
	ffmpeg_filters       = `[0:v]scale=ih*16/9:-1,boxblur=luma_radius=min(h\,w)/20:luma_power=1:chroma_radius=min(cw\,ch)/20:chroma_power=1[bg];[bg][0:v]overlay=(W-w)/2:(H-h)/2,crop=h=iw*9/16`
)

func ProcessAll(videos []downloader.Video, base string) []ProcessedVideo {
	var (
		wg  sync.WaitGroup
		all []ProcessedVideo
	)

	wg.Add(len(videos))

	for _, v := range videos {
		go func(v downloader.Video) {
			path, out, err := process(v, base)
			if err != nil {
				log.Fatalf("error processing video: %q of this post: %q::: %q\ncommand output: %s", v.Path, v.Post.Hash, err.Error(), out)
			}

			log.Printf("processed video: %q with hash %q on %q\n", v.Post.Title, v.Post.Hash, path)

			all = append(all, ProcessedVideo{
				Path:  path,
				Video: v,
			})
			wg.Done()
		}(v)
	}

	wg.Wait()

	return all
}

func MergeAll(videos []ProcessedVideo, base, output string) error {
	t, err := createVideosFile(videos, base)
	if err != nil {
		return err
	}
	log.Printf("created temp file: %q", t)

	err = runMerge(output, t)
	if err != nil {
		return err
	}

	return os.Remove(t)
}

func runMerge(o, t string) error {
	cmdString := fmt.Sprintf(ffmpeg_merge_command, t, o)
	args := strings.Split(cmdString, " ")
	log.Printf("running command: %q", cmdString)
	cmd := exec.Command(args[0], args[1:]...)

	_, err := helpers.RunCmd(cmd)
	return err
}

func createVideosFile(videos []ProcessedVideo, base string) (string, error) {
	n := path.Join(base, "all_videos.txt")

	f, err := os.Create(n)
	if err != nil {
		return "", err
	}
	defer f.Close()

	for _, v := range videos {
		f.WriteString(fmt.Sprintf("file %s\n", v.Path))
	}

	return n, nil
}

func process(video downloader.Video, base string) (string, string, error) {
	o := path.Join(base, "blurry", video.Post.Hash)
	cmdString := fmt.Sprintf(ffmpeg_command, video.Path, ffmpeg_filters, ffmpeg_quality, o)
	args := strings.Split(cmdString, " ")

	log.Printf("running command: %q", cmdString)
	cmd := exec.Command(args[0], args[1:]...)
	out, err := helpers.RunCmd(cmd)
	if strings.Contains(out, "already exists. Exiting") {
		err = nil
	}

	return args[len(args)-2], out, err
}
