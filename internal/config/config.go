package config

import (
	"fmt"
	"os"
)

func Load() string {
	key := os.Getenv("GEMINI_KEY")
	url := os.Getenv("GEMINI_URL")

	// todo: check for empty envs

	return fmt.Sprintf("%s%s", url, key)
}
