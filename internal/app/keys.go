// internal/app/keys.go
package app

import (
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (s *State) HandleKey(msg tea.KeyMsg, vp *viewport.Model) tea.Cmd {

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

		now := time.Now()

		switch msg.String() {
		case "j", "down":
			if s.CursorYOffset < len(s.Lines)-1 {
				s.CursorYOffset++

				// Reset to end of line if two long
				lineLength := max(0, len(s.Lines[s.CursorYOffset]))
				if s.CursorXOffset > lineLength {
					s.CursorXOffset = lineLength
				}

				// scroll down if needed
				if s.CursorYOffset > vp.YOffset+vp.VisibleLineCount()-1 {
					vp.ScrollDown(1)
				}
			}
		case "k", "up":
			if s.CursorYOffset > 0 {
				s.CursorYOffset--

				// Reset to end of line if two long
				lineLength := max(0, len(s.Lines[s.CursorYOffset])-1)
				if s.CursorXOffset > lineLength {
					s.CursorXOffset = lineLength
				}
				// scroll up if needed
				if s.CursorYOffset < vp.YOffset {
					vp.ScrollUp(1)
				}
			}
		case "h", "left":
			if s.CursorXOffset > 0 {
				s.CursorXOffset--
			}
		case "l", "right":
			if s.CursorXOffset < len(s.Lines[s.CursorYOffset])-1 {
				s.CursorXOffset++
			}
		case "g":
			if isDoubleTap(msg.String(), now, s) {
				vp.GotoTop()
				s.CursorYOffset = 0
				s.CursorXOffset = 0
			}
		case "G":
			vp.GotoBottom()
			s.CursorYOffset = vp.TotalLineCount() - 1
			s.CursorXOffset = 0
		case ":":
			s.Mode = ModeCommand
			s.CommandType = ":"
		case "/":
			s.Mode = ModeCommand
			s.CommandType = "/"
		case "i", "c":
			s.Mode = ModeChat
		}

		s.LastKey = msg.String()
		s.LastKeyTime = now
	}

	return nil
}

func isDoubleTap(key string, currTime time.Time, s *State) bool {

	if key == s.LastKey && currTime.Sub(s.LastKeyTime) < 250*time.Millisecond {
		return true
	}
	return false
}
