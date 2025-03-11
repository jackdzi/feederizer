package feed

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackdzi/feederizer/ui/internal/api"
	"github.com/jackdzi/feederizer/ui/internal/pages/page"
	"github.com/jackdzi/feederizer/ui/internal/theme"
)

type model struct {
	tableStyles table.Styles
	table       table.Model
	ids         []int

	styles theme.Styles
	help   help.Model
	keys   keyMap
	err    error
}

func (m model) Init() tea.Cmd {
	return nil
}

func New(styles theme.Styles) page.Model {
	columns := []table.Column{
		{Title: "Article", Width: 20},
		{Title: "Source", Width: 20},
		{Title: "Date", Width: 20},
		{Title: "Ids", Width: 20},
	}

	rows := []table.Row{
		{"Feed Item 1", "WSJ", "1/12", "1"},
		{"Feed Item 2", "AllPoets", "1/12", "2"},
		{"Feed Item 3", "NYT", "1/12", "3"},
	}
	tableStyles := table.DefaultStyles()
	tableStyles.Cell.Foreground(lipgloss.Color(styles.TextColor))
	tableStyles.Header.Foreground(lipgloss.Color(styles.TextColor))
	tableStyles.Selected = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(styles.TextHighlight)).Foreground(lipgloss.Color(styles.TextColor)).Width(styles.Feed.TableWidth)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(styles.Feed.TableHeight),
		table.WithStyles(tableStyles),
	)

	return &model{
		help:        help.New(),
		keys:        keys_config,
		styles:      styles,
		table:       t,
		tableStyles: tableStyles,
	}
}

func (m model) UpdateSize() page.Model {
	m.styles.ApplySizes()
	m.tableStyles.Selected = m.tableStyles.Selected.Width(m.styles.Feed.TableWidth)
	m.table.SetStyles(m.tableStyles)
	m.table.SetWidth(m.styles.Feed.TableWidth)
	m.table.SetHeight(m.styles.Feed.TableHeight)
	return m
}

func (m model) Update(msg tea.Msg) (page.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.ClearDatabase):
			return m, page.ReturnClearDatabase
		case key.Matches(msg, m.keys.Quit):
			return m, page.ReturnQuit
		case key.Matches(msg, m.keys.LogOut):
			return m, page.ReturnBackMsg
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Enter):
      api.GetSubscribedFeedItems()
      return m, tea.Quit
      return m, page.ReturnViewer
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	helpview := m.help.View(m.keys)
	spaces := ""
	if !m.help.ShowAll {
		spaces = "\n\n\n"
	}
	newTable := m.returnRemoveColumn("Ids")
	return "\n\n\n" + lipgloss.NewStyle().MarginLeft(4).Render(newTable.View()) + "\n" + spaces + helpview
}

func (m *model) returnRemoveColumn(title string) table.Model {
	index := -1
	for i, col := range m.table.Columns() {
		if strings.EqualFold(col.Title, title) {
			index = i
			break
		}
	}

	if index == -1 {
		return m.table
	}

	// fmt.Println(m.table.SelectedRow()[1])

	newRows := make([]table.Row, len(m.table.Rows()))
	for i, row := range m.table.Rows() {
		newRows[i] = table.Row{row[0], row[1], row[2]}
	}

	newCols := make([]table.Column, len(m.table.Columns()))
	for i, col := range m.table.Columns() {

		if i == index {
			continue
		}
		newCols[i] = table.Column{Title: col.Title, Width: col.Width}
	}

	t := table.New(
		table.WithColumns(newCols),
		table.WithRows(newRows),
		table.WithFocused(true),
		table.WithHeight(m.table.Height()),
		table.WithStyles(m.tableStyles),
	)

  selectedRow := m.table.Cursor()
	if selectedRow >= len(newRows) {
		selectedRow = len(newRows) - 1
	}
	if selectedRow < 0 {
		selectedRow = 0
	}
	t.SetCursor(selectedRow)
	return t
}

func (m model) GetPageTitle() string {
	return "Feed"
}
