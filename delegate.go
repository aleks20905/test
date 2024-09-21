package main

import (
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

var prevItemIndex *int
var (
	mu       sync.Mutex
	lastCopy time.Time
)

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

			// show the password for the specific item
			case key.Matches(msg, keys.showPass):
				index := m.Index()
				prevItemIndex = &index
				// Create a modified copy of the selected item
				i := m.Items()[index].(item)
				i.showPass = !i.showPass
				// Set the modified item back to the slice
				m.Items()[index] = i
				return m.NewStatusMessage(statusMessageStyle("shows " + i.pass))

			// copyes the passowrd to the clipboard and afther spec_time clears the clipboard
			case key.Matches(msg, keys.showPass2):
				index := m.Index()
				i := m.Items()[index].(item)
				err := copyToClipboardWithTimeout(i.pass, 10*time.Second)
				if err != nil {
					return m.NewStatusMessage(statusMessageStyle("Problem: " + err.Error()))
				}
				return m.NewStatusMessage(statusMessageStyle("Copied to clipboard (will clear in 10 seconds)"))

			// after showind the password this hides it after u move
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

func copyToClipboardWithTimeout(text string, duration time.Duration) error {
	mu.Lock()
	defer mu.Unlock()

	if err := clipboard.Init(); err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(text))
	lastCopy = time.Now()

	go func() {
		time.Sleep(duration)
		mu.Lock()
		defer mu.Unlock()
		if time.Since(lastCopy) >= duration {
			clipboard.Write(clipboard.FmtText, []byte(""))
		}
	}()

	return nil
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
			key.WithKeys("y"),
			key.WithHelp("y", "copy pass"),
		),
		movement: key.NewBinding(
			key.WithKeys("up", "down", "k", "j", "g", "ctrl+g", "home", "end"),
			// key.WithHelp("h", "show pass"),
		),
	}
}
