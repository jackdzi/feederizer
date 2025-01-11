package login

import (
	"bytes"
	"encoding/json"
	"github.com/jackdzi/feederizer/ui/internal/pages/page"
	"github.com/jackdzi/feederizer/ui/internal/theme"
	"fmt"
	"net/http"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	inputs    []textinput.Model
	focused   int
	formatter textinput.Model

	serverError     bool
	credentialError bool
	styles          theme.Styles
	err             error
}

const (
	username = iota
	password
)

func New(styles theme.Styles) page.Model {
	inputs := make([]textinput.Model, 2)
	inputs[username] = textinput.New()
	inputs[username].Placeholder = "Username Here"
	inputs[username].Focus()
	inputs[username].CharLimit = 20
	inputs[username].Width = 13
	inputs[username].Prompt = ""

	inputs[password] = textinput.New()
	inputs[password].Placeholder = "Password Here"
	inputs[password].Blur()
	inputs[password].CharLimit = 35
	inputs[password].Width = 13
	inputs[password].EchoMode = textinput.EchoPassword
	inputs[password].Prompt = ""

	formatter := textinput.New()
	formatter.Blur()
	formatter.Prompt = ""

	return &model{
		inputs:    inputs,
		styles:    styles,
		formatter: formatter,
	}
}

func (m model) Update(msg tea.Msg) (page.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "esc":
			m.reset()
			return m, page.ReturnBackMsg
		case "enter":
			if m.focused == password {
				if err := m.submitUserData(); err != nil {
					fmt.Println("Error: ", err)
					return m, tea.Quit
				}
				m.reset()
				if m.serverError {
					m.inputs[username].Prompt = "\033[31mUser \033[31mnot \033[31mFound\033[0m\n"
				} else if m.credentialError {
					m.inputs[username].Prompt = "\033[31mInvalid \033[31mCredentials\033[0m\n"
				} else {
					return m, page.ReturnBackMsg
				}
			} else {
				m.focused = password
			}
		case tea.KeyTab.String():
			m.nextInput()
		case tea.KeyShiftTab.String():
			m.prevInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	case error:
		m.err = msg
		return m, nil
	}
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m model) UpdateSize() page.Model {
	m.styles.ApplySizes()
	return m
}

func (m *model) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m *model) reset() {
	m.inputs[username].Reset()
	m.inputs[username].Prompt = ""
	m.inputs[password].Reset()
	m.focused = username
}

func (m *model) submitUserData() error {
	data := map[string]string{"name": m.inputs[username].Value(), "password": m.inputs[password].Value()}
	m.reset()

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error parsing string to JSON")
		return err
	}

	resp, err := http.Post("http://localhost:8080/credentials", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	m.serverError = (resp.StatusCode == http.StatusInternalServerError)
	m.credentialError = (resp.StatusCode == http.StatusUnauthorized)
	return nil
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) GetPageTitle() string {
	return "Input"
}

func (m model) View() string {
	if m.err != nil {
		return "help"
	}

	return m.styles.Centering.Render(
		m.styles.Input.Box.Render(m.inputs[username].View()+
			"\n"+"\n"+
			m.inputs[password].View()+"\n"+m.formatter.View())) +
		m.styles.RenderInstructions() + m.styles.Input.BottomPadding.Render()
}
