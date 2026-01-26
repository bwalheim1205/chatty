package app

import (
	"context"
	"strings"

	"github.com/bwalheim1205/chatty/internal/llm"
	tea "github.com/charmbracelet/bubbletea"
)

func (s *State) executeCommand() tea.Cmd {

	commandParts := strings.Split(s.Command, " ")

	switch commandParts[0] {
	case "q", "q!":
		return tea.Quit
	case "m", "model":
		if len(commandParts) > 1 {
			s.Model = commandParts[1]
		}
	case "c", "c!", "clear":
		s.Messages = []llm.Message{}
	}

	s.Command = ""
	s.Mode = ModeRead

	return nil
}

func (s *State) chat() tea.Cmd {
	// cancel any in-flight request
	if s.cancel != nil {
		s.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	// append user message
	s.Messages = append(s.Messages, llm.Message{
		Role:    "user",
		Content: s.Command,
	})

	req := llm.Request{
		Model:    llm.ModelID(s.Model),
		Messages: s.Messages,
	}

	s.Command = ""
	s.Mode = ModeRead

	stream, err := s.LLM.Stream(ctx, req)
	if err != nil {
		return func() tea.Msg {
			return LLMStreamChunk{Err: err}
		}
	}

	s.Stream = stream

	return ReadNextChunk(ctx, stream)
}

func ReadNextChunk(
	ctx context.Context,
	stream <-chan llm.StreamChunk,
) tea.Cmd {
	return func() tea.Msg {
		select {
		case <-ctx.Done():
			// context was cancelled or timed out
			return LLMStreamChunk{Err: ctx.Err()}

		case chunk, ok := <-stream:
			if !ok {
				// stream closed naturally
				return LLMStreamChunk{Done: true}
			}

			// forward any chunk error
			if chunk.Err != nil {
				return LLMStreamChunk{Err: chunk.Err}
			}

			// normal chunk
			return LLMStreamChunk{
				Text: chunk.Text,
				Done: chunk.Done,
			}
		}
	}
}
