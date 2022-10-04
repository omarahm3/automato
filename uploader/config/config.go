package config

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/omarahm3/automato/database"
	"github.com/omarahm3/automato/types"
)

const (
	database_timeout = 1000 * 60 * time.Second
	VideoTypeRandom  = "random"
	VideoTypeManual  = "manual"
)

type Config struct {
	TokenFile        string
	SecretsFile      string
	OutputFile       string
	PrivacyStatus    string
	VideoInfoType    string
	VideoTitle       string
	VideoDescription string
	DatabaseURI      string
}

func (c *Config) GetTitle() (string, error) {
	switch c.VideoInfoType {
	case VideoTypeManual:
		return c.VideoTitle, nil
	case VideoTypeRandom:
		db, cancel, err := getDbClient(c)
		if err != nil {
			return "", err
		}
		defer cancel()

		var posts []types.Post
		err = db.FindRandomOfToday(1, &posts)
		if err != nil {
			return "", err
		}
		p := posts[0]
		return fmt.Sprintf("Top 10 best/cringiest Tiktoks today: %s", p.Title), nil
	default:
		return "", fmt.Errorf("unknown video type %v", c.VideoInfoType)
	}
}

func (c *Config) GetDescription() (string, error) {
	switch c.VideoInfoType {
	case VideoTypeManual:
		return c.VideoDescription, nil
	case VideoTypeRandom:
		return c.VideoDescription, nil
	default:
		return "", fmt.Errorf("unknown video type %v", c.VideoInfoType)
	}
}

func getDbClient(c *Config) (*database.Database, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), database_timeout)
	db, err := database.Connect(c.DatabaseURI, ctx)
	return db, cancel, err
}

func LoadConfig() (*Config, error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		return nil, err
	}

	return &Config{
		TokenFile:        myEnv["TOKEN_FILE"],
		SecretsFile:      myEnv["SECRETS_FILE"],
		OutputFile:       myEnv["OUTPUT_FILE"],
		PrivacyStatus:    myEnv["PRIVACY_STATUS"],
		VideoInfoType:    myEnv["VIDEO_INFO_TYPE"],
		VideoTitle:       myEnv["VIDEO_TITLE"],
		VideoDescription: myEnv["VIDEO_DESCRIPTION"],
		DatabaseURI:      myEnv["DATABASE_URI"],
	}, nil
}

func ToInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		panic(err)
	}

	return int(i)
}
