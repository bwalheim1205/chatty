package app

import (
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
	s.Messages = append(s.Messages, llm.Message{
		Role:    "user",
		Content: s.Command,
	})
	s.Command = ""
	s.Mode = ModeRead
	return nil
}
