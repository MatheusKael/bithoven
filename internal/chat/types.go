package chat

import (
	"context"
	"strings"
)

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type PartKind string

const (
	PartText PartKind = "text"
	PartCode PartKind = "code"
)

type Part struct {
	Kind     PartKind `json:"kind"`
	Text     string   `json:"text"`
	Language string   `json:"language,omitempty"`
}

type Message struct {
	Role  Role   `json:"role"`
	Parts []Part `json:"parts"`
}

func (m Message) Text() string {

	var sb strings.Builder

	for _, p := range m.Parts {

		sb.WriteString(p.Text)
	}

	return sb.String()
}

type Response struct {
	Message Message
	Raw     any
}

type Provider interface {
	Name() string
	Generate(ctx context.Context, history []Message) (Response, error)
}
