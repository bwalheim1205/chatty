package app

import (
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
	Mode        Mode
	Messages    []llm.Message
	Command     string
	CommandType string
	Model       string
	LLM         llm.Client
}

func NewState() *State {
	defaultClient := ollama.New("", "")
	return &State{
		Mode:  ModeRead,
		LLM:   defaultClient,
		Model: string(defaultClient.DefaultModel()),
	}
}
