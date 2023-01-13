package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const terminalWidth = 90

type (
	model struct {
		viewport  viewport.Model
		textinput textinput.Model

		surface *surface
	}

	forceUpdateMsg struct{}
)

var _ tea.Model = new(model)

func initialModel(surface *surface) model {
	m := model{
		textinput: textinput.New(),
		surface:   surface,
	}

	m.textinput.Width = terminalWidth
	m.textinput.Placeholder = "Send a command..."
	m.textinput.Prompt = "| "
	m.textinput.Focus()

	m.viewport = viewport.New(terminalWidth, m.height(40))
	m.viewport.SetContent(m.surface.buf.String())
	m.viewport.GotoBottom()

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.textinput.Reset()
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.surface.glas.SendInput(m.textinput.Value())
			m.textinput.Reset()
		}

	case forceUpdateMsg:
		m.viewport.SetContent(m.surface.buf.String())
		m.viewport.GotoBottom()

	case tea.WindowSizeMsg:
		m.viewport.Height = m.height(msg.Height)

	case error:
		m.surface.errCh <- msg
		return m, nil
	}

	m.viewport, cmd = m.viewport.Update(message)
	cmds = append(cmds, cmd)

	m.textinput, cmd = m.textinput.Update(message)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewport.View(),
		m.textinput.View(),
	)
}

func (m model) height(h int) int {
	return h - lipgloss.Height(m.textinput.View())
}
