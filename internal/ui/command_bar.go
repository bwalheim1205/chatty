package ui

import (
	"strconv"
	"strings"

	"github.com/bwalheim1205/chatty/internal/app"
	"github.com/charmbracelet/lipgloss"
)

func AddModeTitle(m Model, b *strings.Builder) {
	modeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Align(lipgloss.Center).
		Width(10)

	providerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA"))

	modelStyle := lipgloss.NewStyle().
		Bold(false).
		Foreground(lipgloss.Color("#FAFAFA"))

	switch m.State.Mode {
	case app.ModeCommand:
		modeStyle = modeStyle.Background(lipgloss.Color("#f4d756ff"))
		b.WriteString(modeStyle.Render(" COMMAND "))
	case app.ModeChat:
		modeStyle = modeStyle.Background(lipgloss.Color("#ad67dfff"))
		b.WriteString(modeStyle.Render(" CHAT "))
	case app.ModeRead:
		modeStyle = modeStyle.Background(lipgloss.Color("#3cb171ff"))
		b.WriteString(modeStyle.Render(" READ "))
	}

	b.WriteString(" ")
	b.WriteString(providerStyle.Render(strings.ToUpper(m.State.LLM.Name())))
	b.WriteString(" - ")
	b.WriteString(modelStyle.Render(m.State.Model))
	b.WriteString("\n")
}

func CommandBar(m Model) string {

	var b strings.Builder

	if m.State.Mode == app.ModeCommand {
		AddModeTitle(m, &b)
		b.WriteString(m.State.CommandType + m.State.Command)
	} else if m.State.Mode == app.ModeChat {
		AddModeTitle(m, &b)
		if len(m.State.Command)+2 > m.Viewport.Width {
			b.WriteString(".." + m.State.Command[len(m.State.Command)+2-m.Viewport.Width:])
		} else {
			b.WriteString("> " + m.State.Command)
		}
	} else {
		AddModeTitle(m, &b)

		// Add line numbers
		rightText := "Ln " + strconv.Itoa(m.State.CursorYOffset+1) + ", Col " + strconv.Itoa(m.State.CursorXOffset)
		statusLine := lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(m.Viewport.Width/2).Align(lipgloss.Left).Render("press ':' for commands"),
			lipgloss.NewStyle().Width(m.Viewport.Width/2).Align(lipgloss.Right).Render(rightText),
		)
		b.WriteString(statusLine)

	}

	return b.String()
}
