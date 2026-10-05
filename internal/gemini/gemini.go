package gemini

import (
	"beethoven/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Candidates []Candidate `json:"candidates,omitempty"`
}

type Candidate struct {
	Message      *Message `json:"content"`
	FinishReason string   `json:"finish_reason,omitempty"`
	Index        int64    `json:"index,omitempty"`
}

type Payload struct {
	Contents []Message `json:"contents,omitempty"`
}

type Message struct {
	Parts []Part `json:"parts,omitempty"`
	Role  string `json:"role,omitempty"`
}

type Req struct {
	Contents []*Message `json:"contents"`
}

type Part struct {
	Text string `json:"text,omitempty"`
}

// Encode payload to fetch the gemini's api endpoint
func EncodePayload(contents *Req) (*bytes.Buffer, error) {
	payload := new(bytes.Buffer)

	err := json.NewEncoder(payload).Encode(contents)

	if err != nil {
		log.Fatal(err)
	}

	return payload, nil
}

func DecodePayload(p []byte) (Response, error) {

	var response Response

	if err := json.Unmarshal(p, &response); err != nil {
		return Response{}, err
	}

	return response, nil
}

func GenerateContent(contents []*Message) (Response, error) {

	req := &Req{
		Contents: contents,
	}

	payload, err := EncodePayload(req)

	if err != nil {
		return Response{}, err
	}

	url := config.Load(".")

	b, err := Fetch(url, payload)

	if err != nil {
		return Response{}, err
	}

	return DecodePayload(b)
}

// not the place for this, but fuck it
func Fetch(url string, payload *bytes.Buffer) ([]byte, error) {
	p := make([]byte, payload.Len())

	payload.Read(p)

	bff := bytes.NewReader(p)

	res, err := http.Post(url, "application/json", bff)

	if res.StatusCode == http.StatusBadRequest {
		defer res.Body.Close()

		bBytes, err := io.ReadAll(res.Body)

		if err != nil {
			log.Printf("failed to read body of bad request: %v", err)
		}

		return nil, fmt.Errorf("bad request: %v\n json: %v", string(bBytes), string(p))
	}

	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return b, nil
}
