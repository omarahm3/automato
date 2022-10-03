package api

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/omarahm3/reddit-aggregator/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const api_url = "https://www.reddit.com/r"

type Post struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Hash      string             `json:"hash" bson:"hash"`
	Title     string             `json:"title" bson:"title"`
	Video     string             `json:"video" bson:"video"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func Fetch(c *config.Config) ([]Post, error) {
	u := buildUrl(c)
	log.Printf("sending request: %s", u)

	body, err := request(u)
	if err != nil {
		return nil, err
	}

	posts, err := parsePosts(body)
	if err != nil {
		return nil, err
	}

	log.Printf("parsed [%d] posts", len(posts))

	return posts, nil
}

func parsePosts(b []byte) ([]Post, error) {
	var raw map[string]interface{}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return nil, err
	}

	data := raw["data"].(map[string]interface{})

	children := data["children"].([]interface{})

	var posts []Post
	for _, v := range children {
		rawv := v.(map[string]interface{})
		post := rawv["data"].(map[string]interface{})
		t := post["title"].(string)

		if post["media"] == nil {
			continue
		}

		media := post["media"].(map[string]interface{})
		redditVideo := media["reddit_video"].(map[string]interface{})
		videoUrl := redditVideo["fallback_url"].(string)

		posts = append(posts, Post{
			ID:        primitive.NewObjectID(),
			Hash:      hashit(videoUrl),
			Title:     t,
			Video:     videoUrl,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return posts, nil
}

func request(u string) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-agent", "automato")
	req.Header.Add("accept", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func buildUrl(c *config.Config) string {
	return fmt.Sprintf("%s/%s/%s.json?limit=%d", api_url, c.SubReddit, c.PostsType, c.PostsLimit)
}

func hashit(s string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(s))
	return hex.EncodeToString(algorithm.Sum(nil))
}
