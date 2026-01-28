package app

import (
	"context"
	"time"

	"github.com/bwalheim1205/chatty/internal/llm"
	"github.com/bwalheim1205/chatty/internal/llm/ollama"
)

type Mode int

const (
	ModeRead Mode = iota
	ModeCommand
	ModeChat
)

type State struct {
	Mode          Mode
	Messages      []llm.Message
	CursorYOffset int
	CursorXOffset int
	Lines         []string
	LastKey       string
	LastKeyTime   time.Time
	Command       string
	CommandType   string
	Model         string
	LLM           llm.Client
	cancel        context.CancelFunc
	Stream        <-chan llm.StreamChunk
}

func NewState() *State {
	defaultClient := ollama.New("", "")
	return &State{
		Mode:  ModeRead,
		LLM:   defaultClient,
		Model: string(defaultClient.DefaultModel()),
	}
}
