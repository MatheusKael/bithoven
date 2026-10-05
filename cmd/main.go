package main

import (
	"beethoven/internal/chat"
	"beethoven/internal/gemini"
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	history, err := chat.OpenHistory("chathistory", "history.json1")

	if err != nil {
		log.Fatalf("failed to open history: %v", err)
	}

	c := &gemini.Client{
		Model: "gemini",
	}

	model := chat.InitialModel(c, history)
	p := tea.NewProgram(model)

	p.Run()

}
