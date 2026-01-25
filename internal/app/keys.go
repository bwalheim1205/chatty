// internal/app/keys.go
package app

import tea "github.com/charmbracelet/bubbletea"

func (s *State) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch s.Mode {

	case ModeCommand:
		switch msg.String() {
		case "enter":
			return s.executeCommand()
		case "ctrl+c", "esc":
			s.Command = ""
			s.Mode = ModeRead
		case "backspace":
			if len(s.Command) == 0 {
				s.Mode = ModeRead
			} else {
				s.Command = s.Command[:len(s.Command)-1]
			}
		default:
			if len(msg.String()) == 1 {
				s.Command += msg.String()
			}
		}

	case ModeChat:
		switch msg.String() {
		case "enter":
			return s.chat()
		case "ctrl+c":
			s.Command = ""
			s.Mode = ModeRead
		case "backspace":
			if len(s.Command) == 0 {
				s.Mode = ModeRead
			} else {
				s.Command = s.Command[:len(s.Command)-1]
			}
		default:
			if len(msg.String()) == 1 {
				s.Command += msg.String()
			}
		}

	case ModeRead:
		switch msg.String() {
		case ":":
			s.Mode = ModeCommand
			s.CommandType = ":"
		case "/":
			s.Mode = ModeCommand
			s.CommandType = "/"
		case "i", "c":
			s.Mode = ModeChat
		}
	}

	return nil
}
