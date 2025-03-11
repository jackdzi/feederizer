package home

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jackdzi/feederizer/ui/internal/api"
	"github.com/jackdzi/feederizer/ui/internal/config"
	"github.com/jackdzi/feederizer/ui/internal/pages/page"
	"github.com/jackdzi/feederizer/ui/internal/theme"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	styles theme.Styles
	help   help.Model
	keys   keyMap
	err    error
	table  table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func New(styles theme.Styles) page.Model {
	columns := []table.Column{
		{Title: "", Width: 30},
	}

	rows := []table.Row{
		{"Login "},
		{"Create New User "},
		{"Edit Config "},
	}
  tableStyles := table.DefaultStyles()
  tableStyles.Cell.Foreground(lipgloss.Color(styles.TextColor))
  tableStyles.Header.Foreground(lipgloss.Color(styles.TextColor))
  tableStyles.Selected = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(styles.TextHighlight)).Foreground(lipgloss.Color(styles.TextColor))

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(6),
    table.WithStyles(tableStyles),
	)


	return &model{
		help:   help.New(),
		keys:   keys_config,
		styles: styles,
		table:  t,
	}
}

func (m model) UpdateSize() page.Model {
	m.styles.ApplySizes()
	return m
}

func (m model) Update(msg tea.Msg) (page.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
    case key.Matches(msg, m.keys.ClearDatabase):
      return m, page.ReturnClearDatabase
    case key.Matches(msg, m.keys.InitializeFeed):
      api.InitDb()
    case key.Matches(msg, m.keys.Quit):
      return m, page.ReturnQuit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Enter):
			switch selected := m.table.SelectedRow()[0]; selected {
			case "Login ":
				return m, page.ReturnLogin
			case "Edit Config ":
		    fmt.Print("\033[H\033[2J")
				openFileWithDefaultEditor(config.GetFilePath())
			case "Create New User ":
				return m, page.ReturnUser
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	helpview := m.help.View(m.keys)
  var spaces = ""
  if !m.help.ShowAll {
    spaces = "\n\n\n"
  }
	return m.styles.Centering.Render("\n\n"+m.styles.Homepage.Box.Render(m.table.View())) + m.styles.Homepage.BottomPadding.Render() + spaces + helpview
}

func (m model) GetPageTitle() string {
	return "Homepage"
}

func openFileWithDefaultEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
