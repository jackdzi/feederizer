package driver

import (
	"fmt"

	"github.com/jackdzi/feederizer/ui/internal/api"
	"github.com/jackdzi/feederizer/ui/internal/config"
	"github.com/jackdzi/feederizer/ui/internal/pages/confirmation"
	"github.com/jackdzi/feederizer/ui/internal/pages/feed"
	"github.com/jackdzi/feederizer/ui/internal/pages/reader"
	"github.com/jackdzi/feederizer/ui/internal/pages/home"
	"github.com/jackdzi/feederizer/ui/internal/pages/login"
	"github.com/jackdzi/feederizer/ui/internal/pages/newUser"
	"github.com/jackdzi/feederizer/ui/internal/pages/page"
	"github.com/jackdzi/feederizer/ui/internal/theme"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	Home page.Currentpage = iota
	Login
	addUser
	Confirmation
	Feed
  Reader
)

type model struct {
	page  page.Currentpage
	pages map[page.Currentpage]page.Model
	user  string

	styles theme.Styles
}

func (m *model) Init() tea.Cmd {
	tea.SetWindowTitle("Feederizer")
	style := theme.NewStyles()
	style.ApplySizes()
	m.pages[addUser] = newUser.New(style)
	m.pages[Home] = home.New(style)
	m.pages[Login] = login.New(style)
	m.pages[Confirmation] = confirmation.New(style)
	m.pages[Feed] = feed.New(style)
  m.pages[Reader] = reader.New(style)
	autoLog := m.handleLogin("", "", true)
	if autoLog {
		if config.ReturnConfig().Get("authentication.autoLogin").(bool) {
			m.page = Feed
		}
	}
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case page.Authentication:
		user, pass := msg.AuthData()
		m.handleLogin(user, pass, false)
		m.page = Feed
	case page.Yes:
		api.ClearDatabase()
		m.page = msg.YesData()
	case page.No:
		m.page = msg.NoData()
	case page.ClearDatabase:
		m.page = Confirmation
	case page.Quit:
		fmt.Print("\033[H\033[2J")
		cmds = append(cmds, tea.Quit)
  case page.Feed:
    m.page = Feed
  case page.Viewer:
    m.page = Reader
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

func (m *model) updatePageModel(newModel page.Model, index page.Currentpage) {
	m.pages[index] = newModel
}

func (m *model) getCurrentPageModel() page.Model {
	return m.pages[m.page]
}

func (m *model) View() string {
	currentPageModel := m.getCurrentPageModel()
	return "\n   " + m.styles.Title.Render(" Feederizer") + currentPageModel.View()
}

func New(styles theme.Styles) *tea.Program {
	fmt.Print("\033[H\033[2J")
	return tea.NewProgram(
		&model{
			page:   Home,
			pages:  make(map[page.Currentpage]page.Model),
			styles: styles,
		})
}
