package gemini

import (
	"bytes"
	"encoding/json"
	"log"
)

type Response struct {
	Candidates []Candidate `json:"candidates,omitempty"`
}

type Candidate struct {
	Content      Content `json:"content"`
	FinishReason string  `json:"finish_reason,omitempty"`
	Index        int64   `json:"index,omitempty"`
}

type Payload struct {
	Contents []Content `json:"contents,omitempty"`
}

type Content struct {
	Parts []Part `json:"parts,omitempty"`
	Role  string `json:"role,omitempty"`
}

type Part struct {
	Text string `json:"text,omitempty"`
}

// Encode payload to fetch the gemini's api endpoint
func EncodePayload(text string) (*bytes.Buffer, error) {
	payload := new(bytes.Buffer)

	p := &Payload{
		Contents: []Content{
			{
				Parts: []Part{
					{
						Text: text,
					},
				},
			},
		},
	}

	err := json.NewEncoder(payload).Encode(p)

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
