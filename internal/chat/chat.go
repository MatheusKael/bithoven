package chat

import (
	"beethoven/internal/gemini"
	"fmt"
	"log"
	"strings"

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
	history     *History
	err         error
}

func InitialModel(history *History) model {

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

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		llmStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		history:     history,
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ResponseMsg:
		if msg.err != nil {
			// tratar erro
			return m, nil
		}
		render := m.llmStyle.Render("Gemini: ") + msg.text

		m.history.addMsg(render)

		m.messages = append(
			m.messages,
			render,
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

			render := m.senderStyle.Render("You: ") + text

			m.history.addMsg(render)
			m.messages = append(m.messages, render)

			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			return m, sendRequest(m, text)
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

type ResponseMsg struct {
	text string
	err  error
}

func sendRequest(m model, text string) tea.Cmd {
	return func() tea.Msg {
		if len(text) == 0 {
			log.Print("empty message")
			return nil
		}

		response, err := gemini.Gemini(
			fmt.Sprintf("history: %s\nLast message: %s", m.history.Messages, m.history.LastMsg),
		)

		if err != nil {
			log.Fatal(err)
		}

		return ResponseMsg{
			text: response.Candidates[0].Content.Parts[0].Text,
		}

	}

}
