package config

import (
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TokenFile     string
	SecretsFile   string
	OutputFile    string
	PrivacyStatus string
}

func LoadConfig() (*Config, error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		return nil, err
	}

	return &Config{
		TokenFile:     myEnv["TOKEN_FILE"],
		SecretsFile:   myEnv["SECRETS_FILE"],
		OutputFile:    myEnv["OUTPUT_FILE"],
		PrivacyStatus: myEnv["PRIVACY_STATUS"],
	}, nil
}

func ToInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		panic(err)
	}

	return int(i)
}
