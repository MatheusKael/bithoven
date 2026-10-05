package chat

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	messages    []string
	textarea    textarea.Model
	viewport    viewport.Model
	senderStyle lipgloss.Style
	llmStyle    lipgloss.Style
	codeStyle   lipgloss.Style
	history     *History
	provider    Provider
	err         error
}

func InitialModel(p Provider, history *History) model {

	ta := textarea.New()

	ta.Placeholder = "Send a message..."

	ta.SetVirtualCursor(false)
	ta.Focus()

	ta.Prompt = "┃ "

	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()

	ta.SetStyles(s)

	ta.ShowLineNumbers = false

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent("Chat with LLM")

	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	m := model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		llmStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		codeStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("14")).PaddingLeft(2),
		history:     history,
		provider:    p,
		err:         nil,
	}

	for _, msg := range history.Messages {
		switch msg.Role {
		case RoleUser:
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+msg.Text())
		case RoleAssistant:
			m.messages = append(m.messages, renderAssistance(msg, p.Name(), m.llmStyle, m.codeStyle))
		}
	}

	if len(m.messages) > 0 {
		m.viewport.SetContent(
			lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")),
		)
		m.viewport.GotoBottom()
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

type ResponseMsg struct {
	response Response
	err      error
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ResponseMsg:
		if msg.err != nil {
			// tratar erro
			m.err = msg.err
			return m, nil
		}

		rendered := renderAssistance(msg.response.Message, m.provider.Name(), m.llmStyle, m.codeStyle)

		m.messages = append(
			m.messages,
			rendered,
		)

		m.viewport.SetContent(
			lipgloss.NewStyle().
				Width(m.viewport.Width()).
				Render(strings.Join(m.messages, "\n")),
		)

		m.viewport.GotoBottom()

		return m, nil
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.textarea.Height())

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case "enter", "ctrl+j":
			text := m.textarea.Value()

			if strings.TrimSpace(text) == "" {
				return m, nil
			}

			userMsg := Message{
				Role:  RoleUser,
				Parts: []Part{{Kind: PartText, Text: text}},
			}

			if err := m.history.addMsg(userMsg); err != nil {
				m.err = err
				return m, nil
			}

			m.messages = append(m.messages, m.senderStyle.Render("You: ")+text)
			m.textarea.Reset()
			m.viewport.SetContent(
				lipgloss.NewStyle().Width(m.viewport.Width()).
					Render(strings.Join(m.messages, "\n")),
			)
			m.viewport.GotoBottom()

			// snapshot! evita corrida entre goroutine da UI e a do Cmd
			snapshot := m.history.Snapshot()

			return m, sendRequest(m.provider, snapshot)
		default:
			// Send all other keypresses to the textarea.
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}

	case cursor.BlinkMsg:
		// Textarea should also process cursor blinks.
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() tea.View {
	if m.err != nil {
		log.Print(m.err)
		v := tea.NewView("Error: " + m.err.Error() + "\n")
		v.AltScreen = true
		return v
	}

	viewportView := m.viewport.View()
	v := tea.NewView(viewportView + "\n" + m.textarea.View())
	c := m.textarea.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(viewportView)
	}
	v.Cursor = c
	v.AltScreen = true

	return v
}

func sendRequest(provider Provider, history []Message) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)

		defer cancel()

		resp, err := provider.Generate(ctx, history)

		return ResponseMsg{
			response: resp,
			err:      err,
		}
	}

}

func renderAssistance(m Message, name string, nameStyle, codeStyle lipgloss.Style) string {

	var sb strings.Builder

	sb.WriteString(nameStyle.Render(name + ": "))

	for i, p := range m.Parts {

		if i > 0 {
			sb.WriteString("\n")
		}
		switch p.Kind {
		case PartCode:

			sb.WriteString(codeStyle.Render(p.Text))
		case PartText:
			sb.WriteString(p.Text)
		}

	}

	return sb.String()
}
