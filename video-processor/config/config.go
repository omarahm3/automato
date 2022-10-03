package config

import (
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURI string
	BaseDir     string
	OutputPath  string
	PostsLimit  int
}

func LoadConfig() (*Config, error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURI: myEnv["DATABASE_URI"],
		BaseDir:     myEnv["BASE_DIR"],
		OutputPath:  myEnv["OUTPUT_PATH"],
		PostsLimit:  ToInt(myEnv["POSTS_LIMIT"]),
	}, nil
}

func ToInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		panic(err)
	}

	return int(i)
}
