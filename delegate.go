package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var prevItemIndex *int

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return m.NewStatusMessage(statusMessageStyle("You chose " + title))

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(statusMessageStyle("Deleted " + title))

			case key.Matches(msg, keys.showPass):
				index := m.Index()
				prevItemIndex = &index
				// Create a modified copy of the selected item
				i := m.Items()[index].(item)
				i.showPass = !i.showPass
				// Set the modified item back to the slice
				m.Items()[index] = i
				return m.NewStatusMessage(statusMessageStyle("shows " + i.pass))

			case key.Matches(msg, keys.showPass2):

				b := m.Items()[m.Index()].(item)

				// i := m.SelectedItem().(item)
				return m.NewStatusMessage(statusMessageStyle("shows " + b.pass))

			case key.Matches(msg, keys.movement):
				if prevItemIndex != nil {
					// Create a modified copy of the selected item
					i := m.Items()[*prevItemIndex].(item)
					if i.showPass {
						i.showPass = false
						// Set the modified item back to the slice
						m.Items()[*prevItemIndex] = i
						return nil
					}
					return nil

				}
				return nil

			}

		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove, keys.showPass, keys.showPass2}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	choose    key.Binding
	remove    key.Binding
	showPass  key.Binding
	showPass2 key.Binding
	movement  key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
		d.showPass,
		d.showPass2,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
			d.showPass,
			d.showPass2,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
		showPass: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "show pass"),
		),
		showPass2: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("h", "show pass"),
		),
		movement: key.NewBinding(
			key.WithKeys("up", "down", "k", "j", "g", "ctrl+g", "home", "end"),
			key.WithHelp("h", "show pass"),
		),
	}
}
