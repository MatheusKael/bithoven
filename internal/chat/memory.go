package chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type History struct {
	Messages []Message
	File     *os.File
}

func (h *History) addMsg(msg Message) error {
	b, err := json.Marshal(msg)

	if err != nil {
		log.Printf("failed to marshal: %v", err)
		return err
	}

	b = append(b, '\n')

	if _, err := h.File.Write(b); err != nil {
		log.Printf("failed to write to file: %v", err)
		return err
	}

	h.Messages = append(h.Messages, msg)
	return nil
}

func (h *History) Snapshot() []Message {

	return append([]Message(nil), h.Messages...)
}

func (h *History) Close() error {

	if h.File == nil {
		return nil
	}

	return h.File.Close()
}

func (h *History) load() error {
	if _, err := h.File.Seek(0, 0); err != nil {
		return err
	}
	sc := bufio.NewScanner(h.File)

	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {

		line := sc.Bytes()

		if len(line) == 0 {
			continue
		}
		var m Message

		if err := json.Unmarshal(line, &m); err != nil {
			return err
		}

		h.Messages = append(h.Messages, m)
	}

	if err := sc.Err(); err != nil {
		return err
	}

	_, err := h.File.Seek(0, 2)

	return err
}

func OpenHistory(dir, filename string) (*History, error) {

	err := os.MkdirAll(dir, 0o750)

	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, filename)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)

	if err != nil {

		if info, statErr := os.Stat(path); statErr == nil {
			return nil, fmt.Errorf("open %s: %w  (%v)", path, err, info.Mode())
		}

		if info, statErr := os.Stat(dir); statErr == nil {
			return nil, fmt.Errorf("open %s: %w  (%v)", path, err, info.Mode())
		}

		return nil, err
	}

	h := &History{File: f, Messages: []Message{}}

	if err := h.load(); err != nil {
		f.Close()
		return nil, err
	}

	return h, nil
}
