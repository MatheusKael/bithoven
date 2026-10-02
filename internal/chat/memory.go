package chat

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type History struct {
	Messages []string
	LastMsg  string
	File     *os.File
}

func (h *History) addMsg(msg string) {

	if h.Messages == nil {
		log.Fatal("forgot to initialize chat history")
	}

	err := h.write(msg)

	if err != nil {
		log.Fatalf("failed to write history file: %v", err)
	}
	prevMsg := h.LastMsg
	h.LastMsg = msg

	h.Messages = append(h.Messages, prevMsg)
}

func (h *History) write(msg string) error {

	if h.File == nil {
		f, err := historyFile()

		if err != nil {
			return err
		}
		h.File = f
	}
	_, err := h.File.Write([]byte(msg + "\n"))

	if err != nil {
		return err
	}

	return nil
}

func historyFile() (*os.File, error) {

	dir := "chathistory"

	err := os.MkdirAll(dir, 0o750)

	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, "history.txt")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)

	if err != nil {

		if info, statErr := os.Stat(path); statErr == nil {
			return nil, fmt.Errorf("open %s: %w  (%v)", path, err, info.Mode())
		}

		if info, statErr := os.Stat(dir); statErr == nil {
			return nil, fmt.Errorf("open %s: %w  (%v)", path, err, info.Mode())
		}

		return nil, err
	}

	return f, nil
}
