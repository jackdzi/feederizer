package driver

import (
	"github.com/jackdzi/feederizer/ui/internal/pages/home"
	"github.com/jackdzi/feederizer/ui/internal/pages/login"
	"github.com/jackdzi/feederizer/ui/internal/pages/newUser"
	"github.com/jackdzi/feederizer/ui/internal/pages/page"
	"github.com/jackdzi/feederizer/ui/internal/theme"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	mode        int8
	currentpage int8
)

const (
	Home currentpage = iota
	Login
	Loading
	Feeds
	addUser
)

type model struct {
	page  currentpage
	pages map[currentpage]page.Model

	styles theme.Styles
}

func (m *model) Init() tea.Cmd {
	tea.SetWindowTitle("Feederizer")
	style := theme.NewStyles()
	style.ApplySizes()
	m.pages[addUser] = newUser.New(style)
	m.pages[Home] = home.New(style)
	m.pages[Login] = login.New(style)
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg.(type) {
	case page.Quit:
		fmt.Print("\033[H\033[2J")
		cmds = append(cmds, tea.Quit)
	case page.Login:
		m.page = Login
		m.pages[m.page] = login.New(m.styles)
	case page.User:
		m.page = addUser
		m.pages[m.page] = newUser.New(m.styles)
	case page.BackMsg:
		m.page = Home
	case tea.WindowSizeMsg:
		for index, page := range m.pages {
			newModel := page.UpdateSize()
			m.updatePageModel(newModel, index)
		}
		m.styles.ApplySizes()
	}

	currentPageModel := m.getCurrentPageModel()
	if currentPageModel != nil {
		newPageModel, cmd := currentPageModel.Update(msg)
		m.updatePageModel(newPageModel, m.page)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *model) updatePageModel(newModel page.Model, index currentpage) {
	m.pages[index] = newModel
}

func (m *model) getCurrentPageModel() page.Model {
	return m.pages[m.page]
}

func (m *model) View() string {
	currentPageModel := m.getCurrentPageModel()
	return " " + m.styles.Title.Render(" Feederizer") + currentPageModel.View()
}

func New(styles theme.Styles) *tea.Program {
	return tea.NewProgram(
		&model{
			page:   Home,
			pages:  make(map[currentpage]page.Model),
			styles: styles,
		})
}
