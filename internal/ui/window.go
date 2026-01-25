package ui

import (
	"strings"
)

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(ChatView(m))
	b.WriteString("\n")
	b.WriteString(CommandBar(m))

	return b.String()
}
