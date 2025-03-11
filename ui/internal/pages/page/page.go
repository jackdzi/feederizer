package page

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackdzi/feederizer/ui/internal/config"
)

type BackMsg struct{}

type Reset struct{}

type User struct{}

type Login struct{}

type Quit struct{}

type Feed struct{}

func ReturnFeed() tea.Msg {
  return Feed{}
}

type Viewer struct{}

func ReturnViewer() tea.Msg {
  return Viewer{}
}

type (
	Currentpage int8
)

type No struct{
  page Currentpage
}

type Yes struct{
  page Currentpage
}

type Authentication struct {
	user string
  pass string
}

func (no No) NoData() Currentpage {
  return no.page
}

func (yes Yes) YesData() Currentpage {
  return yes.page
}

func ReturnNo(page Currentpage) tea.Cmd {
  return func() tea.Msg {
    return No{page: page}
  }
}

func (auth Authentication) AuthData() (string, string) {
  return auth.user, auth.pass
}

func ReturnAuthentication(user string, pass string) tea.Cmd {
	return func() tea.Msg {
    passHidden, _ := config.Encrypt(pass)
    return Authentication{user: user, pass: passHidden}
	}
}

func ReturnYes(page Currentpage) tea.Cmd {
  return func() tea.Msg {
    return Yes{page: page}
  }
}

type ClearDatabase struct{}

func ReturnClearDatabase() tea.Msg {
	return ClearDatabase{}
}

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
