package home

import (
	"feederizer/ui/internal/pages/page"
	"feederizer/ui/internal/theme"
	"os/exec"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	styles theme.Styles
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
		{"Login"},
		{"Create New User"},
		{"Edit Config"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(6),
	)

	return &model{styles: styles, table: t}
}

func (m model) UpdateSize() page.Model {
	m.styles.ApplySizes()
	return m
}

func (m model) Update(msg tea.Msg) (page.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "enter":
			switch selected := m.table.SelectedRow()[0]; selected {
			case "Login":
        return m, page.ReturnLogin
			case "Edit Config":
        openFileWithDefaultEditor("../config.toml") // TODO: Change for docker
			case "Create New User":
				return m, page.ReturnUser
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.styles.Title.Render("Feederizer") + m.styles.Centering.Render("\n\n"+m.styles.Homepage.Box.Render(m.table.View())) + m.styles.Homepage.BottomPadding.Render()
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
