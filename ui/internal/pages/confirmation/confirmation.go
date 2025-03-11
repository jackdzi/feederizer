package confirmation

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackdzi/feederizer/ui/internal/pages/page"
	"github.com/jackdzi/feederizer/ui/internal/theme"
)

type model struct {
	styles theme.Styles
}

func (m model) Init() tea.Cmd {
	return nil
}

func New(styles theme.Styles) page.Model {
	return &model{
		styles: styles,
	}
}

func (m model) UpdateSize() page.Model {
	m.styles.ApplySizes()
	return m
}

func (m model) Update(msg tea.Msg) (page.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			return m, page.ReturnYes(0)
		case "n", "N":
			return m, page.ReturnNo(4)
		case "esc", "ctrl+c", "q":
			return m, page.ReturnBackMsg
		}
	}
	return m, nil
}

func (m model) View() string {
	return m.styles.Centering.Render("\n\n" + m.styles.Homepage.Box.Render("Are you sure you want to clear the entire database? Press 'y' to confirm or 'n' to cancel.\n"))
}

func (m model) GetPageTitle() string {
	return "Confirmation"
}
