package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func CommandBar(m Model) string {

	modeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Align(lipgloss.Center).
		Width(10)

	var b strings.Builder

	if m.Mode == ModeCommand {
		// Add Title
		modeStyle = modeStyle.Background(lipgloss.Color("#f4d756ff"))
		b.WriteString(modeStyle.Render(" COMMAND "))
		b.WriteString("\n")

		// Add Content
		b.WriteString(m.CommandType + m.Command)
	} else if m.Mode == ModeChat {
		// Add Title
		modeStyle = modeStyle.Background(lipgloss.Color("#ad67dfff"))
		b.WriteString(modeStyle.Render(" CHAT "))
		b.WriteString("\n")

		// Add Content
		if len(m.Command)+2 > m.Viewport.Width {
			b.WriteString(".." + m.Command[len(m.Command)+2-m.Viewport.Width:])
		} else {
			b.WriteString("> " + m.Command)
		}
	} else {
		// Add Title
		modeStyle = modeStyle.Background(lipgloss.Color("#3cb171ff")).Width(10)
		b.WriteString(modeStyle.Render(" NORMAL "))
		b.WriteString("\n press ':' for commands")
	}

	return b.String()
}
