package gemini

import (
	"beethoven/internal/chat"
	"context"
)

type Client struct {
	Model string
}

func (c *Client) Name() string {
	return "Gemini"
}

func (c *Client) Generate(ctx context.Context, history []chat.Message) (chat.Response, error) {
	contents := make([]*Message, 0, len(history))

	for _, m := range history {

		role := "user"
		if m.Role == chat.RoleAssistant {
			role = "model"
		}

		var parts []Part

		for _, p := range m.Parts {

			part := Part{
				Text: p.Text,
			}

			parts = append(parts, part)
		}

		contents = append(contents, &Message{Role: role, Parts: parts})

	}

	resp, err := GenerateContent(contents)

	if err != nil {
		return chat.Response{}, err
	}

	var out []chat.Part

	for _, cand := range resp.Candidates {
		if cand.Message == nil {
			continue
		}
		for _, p := range cand.Message.Parts {

			switch {
			case p.Text != "":
				out = append(out, chat.Part{Kind: chat.PartText, Text: p.Text})
			}
		}
	}

	return chat.Response{
		Message: chat.Message{Role: chat.RoleAssistant, Parts: out},
		Raw:     resp,
	}, nil

}
