package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	title       string
	description string
	pass        string
	showPass    bool
}

func (i item) Title() string { return i.title }
func (i item) Description() string { // todo change to fit mail
	if i.showPass {
		return fmt.Sprintf("%-15s | %-10s", i.description, i.pass)
	}
	return fmt.Sprintf("%-15s | %-10s", i.description, maskPassword(i.pass))
}

// func (i item) Pass() string { return i.pass }

func (i item) FilterValue() string { return i.title }

func maskPassword(password string) string {
	return strings.Repeat("*", len(password))
}

type listKeyMap struct {
	toggleTitleBar  key.Binding
	toggleStatusBar key.Binding
	toggleHelpMenu  key.Binding
	insertItem      key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle help"),
		),
	}
}

type inputState int

const (
	normalState inputState = iota
	addingItemState
)

type model struct {
	list          list.Model
	itemGenerator *randomItemGenerator
	keys          *listKeyMap
	delegateKeys  *delegateKeyMap
	inputModel    inputModel
	state         inputState
}

func newModel() model {
	var (
		itemGenerator randomItemGenerator
		delegateKeys  = newDelegateKeyMap()
		listKeys      = newListKeyMap()
	)

	// Make initial list of items
	const numItems = 24
	items := make([]list.Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = itemGenerator.next()
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	groceryList := list.New(items, delegate, 200, 200)
	groceryList.Title = "Accaunts"
	groceryList.Styles.Title = titleStyle
	groceryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.toggleHelpMenu,
		}
	}

	// Setup text inputs
	inputModel := newInputModel()

	return model{
		list:          groceryList,
		keys:          listKeys,
		delegateKeys:  delegateKeys,
		itemGenerator: &itemGenerator,
		inputModel:    inputModel,
		state:         normalState,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Handle adding item mode
		if m.state == addingItemState {
			var cmd tea.Cmd
			m.inputModel, cmd = m.inputModel.Update(msg)
			if m.inputModel.Submitted {
				newItem := item{
					title:       m.inputModel.inputs[0].Value(),
					description: m.inputModel.inputs[1].Value(),
					pass:        m.inputModel.inputs[2].Value(),
					showPass:    false,
				}
				insCmd := m.list.InsertItem(0, newItem)
				statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
				m.state = normalState // Return to normal state
				return m, tea.Batch(insCmd, statusCmd)
			}
			return m, cmd
		}

		// Normal list interaction
		switch {
		case key.Matches(msg, m.keys.insertItem):
			m.state = addingItemState
			m.inputModel = newInputModel() // Reset input fields
			return m, nil

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}
	}

	// Normal list update
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel

	return m, cmd
}

func (m model) View() string {
	if m.state == addingItemState {
		return appStyle.Render(m.inputModel.View())
	}
	return appStyle.Render(m.list.View())
}

func main() {
	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
