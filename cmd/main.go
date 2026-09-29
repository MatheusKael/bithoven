package main

import (
	"beethoven/internal/chat"
	"fmt"
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

	model := chat.InitialModel()
	p := tea.NewProgram(model)

	p.Run()

}
