package ui

import (
	"github.com/bwalheim1205/chatty/internal/messages"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	ModeRead Mode = iota
	ModeCommand
	ModeChat
)

type Model struct {
	Mode        Mode
	Messages    []messages.ChatMessage
	Command     string
	CommandType string
	Viewport    viewport.Model
}

func InitialModel() Model {
	vp := viewport.New(0, 0)
	vp.YPosition = 0 // top aligned

	return Model{
		Mode:     ModeRead,
		Viewport: vp,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.Mode {
	case ModeCommand:
		switch msg.String() {
		case "enter":
			return m.executeCommand()
		case "ctrl+c":
			m.Command = ""
			m.Mode = ModeRead
		case "esc":
			m.Mode = ModeRead
		case "backspace":
			if len(m.Command) == 0 {
				m.Mode = ModeRead
			} else {
				m.Command = m.Command[:len(m.Command)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.Command += msg.String()
			}
		}
	case ModeChat:
		switch msg.String() {
		case "enter":
			return m.chat()
		case "ctrl+c":
			m.Command = ""
			m.Mode = ModeRead
		case "backspace":
			if len(m.Command) == 0 {
				m.Mode = ModeRead
			} else {
				m.Command = m.Command[:len(m.Command)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.Command += msg.String()
			}
		}
	case ModeRead:
		switch msg.String() {
		case ":":
			m.Mode = ModeCommand
			m.CommandType = ":"
		case "/":
			m.Mode = ModeCommand
			m.CommandType = "/"
		case "i", "c":
			m.Mode = ModeChat
		}
	}
	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Viewport.Width = msg.Width
		m.Viewport.Height = msg.Height - 2
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case messages.NewLLMResponse:
		m.Messages = append(m.Messages, messages.ChatMessage{
			Role:    "assistant",
			Content: msg.Text,
		})
		return m, nil
	}

	return m, nil
}
