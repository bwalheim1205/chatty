package ui

import (
	"strings"
)

func (m *Model) updateViewportContent() {
	var b strings.Builder

	if len(m.Messages) == 0 {
		// Fill with empty lines to occupy viewport
		for i := 0; i < m.Viewport.Height; i++ {
			b.WriteString("\n")
		}
	} else {
		start := m.Viewport.YOffset
		end := start + m.Viewport.Height
		if end > len(m.Messages) {
			end = len(m.Messages)
		}
		for i := start; i < end; i++ {
			line := m.Messages[i].Role + ": " + m.Messages[i].Content
			b.WriteString(line + "\n")
		}
	}

	m.Viewport.SetContent(b.String())
}

func ChatView(m Model) string {
	m.updateViewportContent()
	return m.Viewport.View()
}
