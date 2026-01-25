// internal/ui/model.go
package ui

import (
	"github.com/bwalheim1205/chatty/internal/app"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	State    *app.State
	Viewport viewport.Model
}

func InitialModel() Model {
	vp := viewport.New(0, 0)
	vp.YPosition = 0

	return Model{
		State:    app.NewState(),
		Viewport: vp,
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
		cmd := m.State.HandleKey(msg)
		return m, cmd

	}

	return m, nil
}
