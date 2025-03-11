package reader

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackdzi/feederizer/ui/internal/pages/page"
	"github.com/jackdzi/feederizer/ui/internal/theme"
)

type model struct {
	styles theme.Styles
  viewer viewport.Model
	help   help.Model
	keys   keyMap
	err    error
}

func (m model) Init() tea.Cmd {
	return nil
}

func New(styles theme.Styles) page.Model {
  v := viewport.New(styles.Viewer.Width, styles.Viewer.Height)
  v.SetContent(content)
	return &model{
		help:   help.New(),
    viewer: v,
		keys:   keys_config,
		styles: styles,
	}
}

func (m model) UpdateSize() page.Model {
	m.styles.ApplySizes()
  m.viewer.Height = m.styles.Viewer.Height
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
    case key.Matches(msg, m.keys.Back):
      return m, page.ReturnFeed
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
    default:
      var cmd tea.Cmd
      m.viewer, cmd = m.viewer.Update(msg)
      return m, cmd
			}
		}
	return m, cmd
}

func (m model) View() string {
	helpview := m.help.View(m.keys)
  if !m.help.ShowAll {
    m.viewer.Height = m.styles.Viewer.Height + 3
  }
	return "\n\n" + m.viewer.View() + "\n\n" +  helpview
}

func (m model) GetPageTitle() string {
	return "Viewer"
}
