package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle   = focusedStyle
	noStyle       = lipgloss.NewStyle()
	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type inputModel struct {
	inputs    []textinput.Model
	focusIdx  int
	Submitted bool
}

func newInputModel() inputModel {
	m := inputModel{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Title"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "desc"

		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}
func (m inputModel) Update(msg tea.Msg) (inputModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIdx == len(m.inputs) {
				m.Submitted = true
				return m, nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIdx--
			} else {
				m.focusIdx++
			}

			if m.focusIdx > len(m.inputs) {
				m.focusIdx = 0
			} else if m.focusIdx < 0 {
				m.focusIdx = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIdx {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Update text inputs
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m inputModel) View() string {
	var b strings.Builder

	for i, input := range m.inputs {
		b.WriteString(input.View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := blurredButton
	if m.focusIdx == len(m.inputs) {
		button = focusedButton
	}
	b.WriteString("\n\n" + button + "\n")

	return b.String()
}
