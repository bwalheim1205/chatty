// internal/ui/model.go
package ui

import (
	"context"

	"github.com/bwalheim1205/chatty/internal/app"
	"github.com/bwalheim1205/chatty/internal/llm"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	State    *app.State
	Viewport viewport.Model
	lines    []string
}

func InitialModel() Model {
	vp := viewport.New(0, 0)
	vp.YPosition = 0

	return Model{
		State:    app.NewState(),
		Viewport: vp,
		lines:    []string{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Viewport.Width = msg.Width
		m.Viewport.Height = msg.Height - 2
		return m, nil

	case tea.KeyMsg:
		cmd := m.State.HandleKey(msg, &m.Viewport)
		m.updateViewportContent()
		return m, cmd

	case app.LLMCompleteMsg:
		m.State.Messages = append(m.State.Messages, llm.Message{
			Role:    "assistant",
			Content: msg.Text,
		})

		m.updateViewportContent()
		m.Viewport.GotoBottom()
		m.State.CursorYOffset = max(0, len(m.State.Lines)-1)
		m.State.CursorXOffset = max(0, len(m.State.Lines[len(m.State.Lines)-1])-1)
		m.updateViewportContent()

		return m, nil

	case app.LLMStreamChunk:
		if msg.Err != nil {
			m.State.Messages = append(m.State.Messages, llm.Message{
				Role:  "assistant",
				Error: msg.Err.Error(),
			})
			return m, nil
		}

		// ensure assistant message exists
		if len(m.State.Messages) == 0 ||
			m.State.Messages[len(m.State.Messages)-1].Role != "assistant" {
			m.State.Messages = append(m.State.Messages, llm.Message{
				Role: "assistant",
			})
		}

		if msg.Text != "" {
			m.State.Messages[len(m.State.Messages)-1].Content += msg.Text
		}

		if msg.Done {
			return m, nil
		}

		m.updateViewportContent()
		m.Viewport.GotoBottom()
		m.State.CursorYOffset = max(0, len(m.State.Lines)-1)
		m.State.CursorXOffset = max(0, len(m.State.Lines[len(m.State.Lines)-1])-1)
		m.updateViewportContent()

		// schedule next chunk
		return m, app.ReadNextChunk(context.Background(), m.State.Stream)
	}

	return m, nil
}
