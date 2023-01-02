package cmd

import (
	"fmt"
	"os"
)

func getSecret() (string, error) {
	env, exists := os.LookupEnv("NOTION_API_KEY")
	if !exists {
		return "", fmt.Errorf("required env variables NOTION_API_KEY not set")
	}
	return env, nil
}

func getDBID() (string, error) {
	env, exists := os.LookupEnv("NOTION_DATABASE")
	if !exists {
		return "", fmt.Errorf("required env variables NOTION_DATABASE not set")
	}
	return env, nil
}
