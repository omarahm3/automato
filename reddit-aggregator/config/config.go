package config

import (
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURI string
	SubReddit   string
	PostsLimit  int
	PostsType   string
}

func LoadConfig() (*Config, error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		return nil, err
	}

	return &Config{
		PostsLimit:  ToInt(myEnv["POSTS_LIMIT"]),
		SubReddit:   myEnv["SUBREDDIT"],
		PostsType:   myEnv["POSTS_TYPE"],
		DatabaseURI: myEnv["DATABASE_URI"],
	}, nil
}

func ToInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		panic(err)
	}

	return int(i)
}
