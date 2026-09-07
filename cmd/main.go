package main

import (
	"beethoven/internal/config"
	"beethoven/internal/gemini"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")

	payload, err := gemini.EncodePayload("ola")

	if err != nil {
		log.Fatal(err)
	}

	b, err := Fetch(payload)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)
}

func Fetch(payload *bytes.Buffer) ([]byte, error) {

	url := config.Load()

	res, err := http.Post(url, "application/json", payload)

	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return b, nil
}
