package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) updateViewportContent() {
	var b strings.Builder

	text_style := lipgloss.NewStyle().Width(m.Viewport.Width)
	user_text_style := lipgloss.NewStyle().Bold(true).Underline(true).Foreground(lipgloss.Color("#c543c9"))
	assistant_text_style := lipgloss.NewStyle().Bold(true).Underline(true).Foreground(lipgloss.Color("#43a1c9"))
	error_style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f14d4d"))

	for _, m := range m.State.Messages {
		switch m.Role {
		case "user":
			b.WriteString(user_text_style.Render("You:") + "\n")
		case "assistant":
			b.WriteString(assistant_text_style.Render("Assistant:") + "\n")
		}

		if m.Error != "" {
			b.WriteString(error_style.Render("Error occured calling an llm please try and resend message:") + "\n")
			b.WriteString(error_style.Render(m.Error))
		} else {
			b.WriteString(text_style.Render(m.Content))
		}

		b.WriteString("\n\n")
	}

	m.Viewport.SetContent(b.String())
	m.Viewport.GotoBottom()
}

func ChatView(m Model) string {
	m.updateViewportContent()
	return m.Viewport.View()
}
