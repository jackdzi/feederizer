package page

import tea "github.com/charmbracelet/bubbletea"

type BackMsg struct{}

type Reset struct{}

type User struct{}

type Login struct{}

type Quit struct{}

func ReturnQuit() tea.Msg {
  return Quit{}
}

func ReturnLogin() tea.Msg {
  return Login{}
}

func ReturnUser() tea.Msg {
  return User{}
}

func ReturnBackMsg() tea.Msg {
	return BackMsg{}
}

func ReturnReset() tea.Msg {
	return Reset{}
}

type Model interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Model, tea.Cmd)
	View() string
	GetPageTitle() string
  UpdateSize() Model
}
