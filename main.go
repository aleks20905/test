package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	platform, username, password string
	showPassword                 bool
}

func (i item) Title() string { return i.platform }
func (i item) Description() string {
	if i.showPassword {
		return fmt.Sprintf("%s | %s", i.username, i.password)
	}
	return fmt.Sprintf("%s | %s", i.username, maskPassword(i.password))
}
func (i item) FilterValue() string { return i.platform }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			index := m.list.Index()
			items := m.list.Items()
			if index >= 0 && index < len(items) {
				currentItem := items[index].(item)
				currentItem.showPassword = !currentItem.showPassword
				items[index] = currentItem
				m.list.SetItems(items)
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func maskPassword(password string) string {
	return strings.Repeat("*", len(password))
}

func generateCustomItems() []list.Item {
	return []list.Item{
		item{platform: "Git", username: "dev@example.com", password: "gitpass123", showPassword: false},
		item{platform: "Gmail", username: "user@gmail.com", password: "gmailpass456", showPassword: false},
		item{platform: "GitHub", username: "githubuser", password: "githubpass789", showPassword: false},
		item{platform: "Gitlab", username: "gitlabuser", password: "gitlabpass101", showPassword: false},
		item{platform: "Bitbucket", username: "bitbucketuser", password: "bitbucketpass202", showPassword: false},
		item{platform: "Jira", username: "jirauser@company.com", password: "jirapass303", showPassword: false},
		item{platform: "Slack", username: "slackuser", password: "slackpass404", showPassword: false},
		item{platform: "AWS", username: "awsuser", password: "awspass505", showPassword: false},
		item{platform: "DigitalOcean", username: "douser", password: "dopass606", showPassword: false},
		item{platform: "Heroku", username: "herokuuser", password: "herokupass707", showPassword: false},
	}
}

func main() {
	items := generateCustomItems()

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Platform Credentials (Press Enter to show/hide password)"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
