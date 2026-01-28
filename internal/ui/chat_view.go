package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type visialLine struct {
	content string
	style   lipgloss.Style
}

func (m *Model) updateViewportContent() {

	text_style := lipgloss.NewStyle()
	user_text_style := text_style.Bold(true).Underline(true).Foreground(lipgloss.Color("#c543c9"))
	assistant_text_style := text_style.Bold(true).Underline(true).Foreground(lipgloss.Color("#43a1c9"))
	error_style := text_style.Bold(true).Foreground(lipgloss.Color("#f14d4d"))
	cursor_style := text_style.Reverse(true)

	visual_lines := make([]visialLine, 0)
	lines := make([]string, 0)

	for _, msg := range m.State.Messages {
		switch msg.Role {
		case "user":
			content := wrapLines("You:", m.Viewport.Width)
			lines = append(lines, content...)
			for _, line := range content {
				visual_lines = append(visual_lines, visialLine{content: line, style: user_text_style})
			}
		case "assistant":
			content := wrapLines("Assistant:", m.Viewport.Width)
			lines = append(lines, content...)
			for _, line := range content {
				visual_lines = append(visual_lines, visialLine{content: line, style: assistant_text_style})
			}
		}
		if msg.Error != "" {
			content := wrapLines("Error occured calling an llm please try and resend message:\n"+msg.Error+"\n", m.Viewport.Width)
			lines = append(lines, content...)
			for _, line := range content {
				visual_lines = append(visual_lines, visialLine{content: line, style: error_style})
			}
		} else {
			content := wrapLines(msg.Content+"\n", m.Viewport.Width)
			lines = append(lines, content...)
			for _, line := range content {
				visual_lines = append(visual_lines, visialLine{content: line, style: text_style})
			}
		}
	}

	// Save state of lines
	m.State.Lines = lines

	// Add cursor and render content
	var b strings.Builder
	for i, line := range visual_lines {
		if i == m.State.CursorYOffset {
			if len(line.content) == 0 {
				b.WriteString(cursor_style.Render(" "))
			} else {
				if m.State.CursorXOffset > 0 {
					b.WriteString(line.style.Render(line.content[:m.State.CursorXOffset]))
				}
				if m.State.CursorXOffset < len(line.content) {
					b.WriteString(cursor_style.Render(string(line.content[m.State.CursorXOffset])))
				}
				if m.State.CursorXOffset != len(line.content)-1 {
					b.WriteString(line.style.Render(line.content[m.State.CursorXOffset+1:]))
				}
			}

		} else {
			b.WriteString(line.style.Width(m.Viewport.Width).Render(line.content))
		}

		if i < len(visual_lines)-1 {
			b.WriteString("\n")
		}
	}

	m.Viewport.SetContent(b.String())

}

func wrapLines(s string, width int) []string {
	var lines []string

	for _, line := range strings.Split(s, "\n") {
		words := strings.Fields(line) // split by whitespace
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		current := words[0]

		for _, word := range words[1:] {
			if len(current)+1+len(word) > width {
				lines = append(lines, current)
				current = word
			} else {
				current += " " + word
			}
		}

		lines = append(lines, current)
	}

	return lines
}

func ChatView(m Model) string {
	return m.Viewport.View()
}
