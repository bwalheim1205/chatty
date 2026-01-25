package ui

import (
	"github.com/bwalheim1205/chatty/internal/messages"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) executeCommand() (tea.Model, tea.Cmd) {

	switch m.Command {
	case "q", "q!":
		return m, tea.Quit
	case "c", "c!", "clear":
		m.Messages = []messages.ChatMessage{}
	}

	m.Command = ""
	m.Mode = ModeRead

	return m, nil
}

func (m *Model) chat() (tea.Model, tea.Cmd) {
	m.Messages = append(m.Messages, messages.ChatMessage{
		Role:    "user",
		Content: m.Command,
	})
	m.Command = ""
	m.Mode = ModeRead
	return m, nil
}
