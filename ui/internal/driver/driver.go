package driver

import (
	"feederizer/ui/internal/pages/home"
	"feederizer/ui/internal/pages/page"
	"feederizer/ui/internal/pages/newUser"
	"feederizer/ui/internal/pages/login"
	"feederizer/ui/internal/theme"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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

	help   help.Model
	keys   keyMap
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
	updatedMsg := msg
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			if m.page == Home {
				m.help.ShowAll = !m.help.ShowAll
			}
		case key.Matches(msg, m.keys.Quit):
			if m.page == Home {
				fmt.Print("\033[H\033[2J")
				cmds = append(cmds, tea.Quit)
			}
		case key.Matches(msg, m.keys.AddUser):
			if m.page == Home {
				m.page = addUser
				updatedMsg = tea.KeyMsg{Type: tea.KeyBackspace}
			}
		case key.Matches(msg, m.keys.DeleteFeed):
			if m.page == Home {
				if err := deleteUsers(); err != nil {
					print("Error: ", err)
					cmds = append(cmds, tea.Quit)
				}
			}
		}
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
		newPageModel, cmd := currentPageModel.Update(updatedMsg)
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
	helpView := m.help.View(m.keys)
	currentPageModel := m.getCurrentPageModel()
	return currentPageModel.View() + helpView
}

func New(styles theme.Styles) *tea.Program {
	return tea.NewProgram(
		&model{
			page:   Home,
			pages:  make(map[currentpage]page.Model),
			help:   help.New(),
			keys:   keys_config,
			styles: styles,
		})
}
