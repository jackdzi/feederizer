package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type mode int

const (
	normalMode mode = iota
	inputMode
)

type model struct {
	content   string
	ready     bool
	viewport  viewport.Model
	mode      mode
	input     string
	cursorPos int
	name      string
	password  string
	state     int

	help help.Model
	keys keyMap
}

const (
	promptForName = iota
	promptForPassword
)

type keyMap struct {
	Up             key.Binding
	Down           key.Binding
	Left           key.Binding
	Right          key.Binding
	Help           key.Binding
	Quit           key.Binding
	AddUser        key.Binding
	InitializeFeed key.Binding
	DeleteFeed     key.Binding
	EnterInputMode key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.AddUser, k.InitializeFeed, k.DeleteFeed, k.EnterInputMode},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "move up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("↓/j", "move down"),
    ),
    Left: key.NewBinding(
        key.WithKeys("left", "h"),
        key.WithHelp("←/h", "move left"),
    ),
    Right: key.NewBinding(
        key.WithKeys("right", "l"),
        key.WithHelp("→/l", "move right"),
    ),
    Help: key.NewBinding(
        key.WithKeys("?"),
        key.WithHelp("?", "toggle help"),
    ),
    Quit: key.NewBinding(
        key.WithKeys("q", "esc", "ctrl+c"),
        key.WithHelp("q", "quit"),
    ),
    AddUser: key.NewBinding(
        key.WithKeys("u"),
        key.WithHelp("u", "add user"),
    ),
    InitializeFeed: key.NewBinding(
        key.WithKeys("I"),
        key.WithHelp("I", "initialize feed"),
    ),
    DeleteFeed: key.NewBinding(
        key.WithKeys("d"),
        key.WithHelp("d", "delete feed"),
    ),
    EnterInputMode: key.NewBinding(
        key.WithKeys("i"),
        key.WithHelp("i", "enter input mode"),
    ),
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.KeyMsg:

		if m.mode == normalMode {
			switch {
			case key.Matches(msg, m.keys.Help): // <-- NEW
				m.help.ShowAll = !m.help.ShowAll
				return m, nil

			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.keys.Up):
			case key.Matches(msg, m.keys.Down):
			case key.Matches(msg, m.keys.Left):
			case key.Matches(msg, m.keys.Right):
			}

			switch msg.String() {

			case "u":
				m.mode = inputMode
				m.state = promptForName
				m.input = ""
				m.cursorPos = 0
				m.content = ""

			case "I":
				if _, err := http.Post("http://localhost:8080/init", "application/json", nil); err != nil {
					fmt.Println("Error:", err)
				}

			case "d":
				if _, err := http.Post("http://localhost:8080/delete", "application/json", nil); err != nil {
					fmt.Println("Error:", err)
				}
				resp, err := http.Get("http://localhost:8080/feed")
				if err != nil {
					fmt.Println("Error:", err)
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Failed to read response body:", err)
				}

				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
					fmt.Printf("Failed to parse JSON: %v\n", err)
				}
				m.content = prettyJSON.String()
				m.viewport.SetContent(prettyJSON.String())

			case "i":
				m.mode = inputMode
				m.input = ""
				m.cursorPos = 0
			}

		} else if m.mode == inputMode {
			switch m.state {
			case promptForName:
				if msg.String() == "enter" {
					m.state = promptForPassword
					m.input = ""
					m.password = ""
					return m, nil
				}
				if msg.String() == "backspace" && len(m.name) > 0 {
					m.name = m.name[:len(m.name)-1]
				} else if msg.String() != "backspace" {
					m.name += msg.String()
				}

			case promptForPassword:
				if msg.String() == "enter" {
					data := map[string]string{
						"name":     m.name,
						"password": m.password,
					}
					jsonData, err := json.Marshal(data)
					if err != nil {
						fmt.Println("Error parsing string to JSON")
					}
					if _, err := http.Post("http://localhost:8080/adduser", "application/json", bytes.NewBuffer(jsonData)); err != nil {
						fmt.Println("Error:", err)
					}

					m.mode = normalMode
					m.input = ""
					m.cursorPos = 0
					m.name = ""
					m.password = ""
					m.state = promptForName
					// Fetch updated feed data
					resp, err := http.Get("http://localhost:8080/feed")
					if err != nil {
						fmt.Println("Error:", err)
					}
					defer resp.Body.Close()

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("Failed to read response body:", err)
					}

					var prettyJSON bytes.Buffer
					if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
						fmt.Printf("Failed to parse JSON: %v\n", err)
					}
					m.content = prettyJSON.String()
					m.viewport.SetContent(m.content)
					return m, nil
				}
				if msg.String() == "backspace" && len(m.password) > 0 {
					m.password = m.password[:len(m.password)-1]
				} else if msg.String() != "backspace" {
					m.password += msg.String()
				}
			}
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {

			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true
			m.viewport.YPosition = headerHeight + 1
			m.help.Width = msg.Width

		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.help.Width = msg.Width
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.mode == normalMode {
		helpView := m.help.View(m.keys)

		return fmt.Sprintf(
			"%s\n%s\n%s\n\n%s",
			m.headerView(),
			m.viewport.View(),
			m.footerView(),
			helpView,
		)
	}

	var prompt string
	switch m.state {
	case promptForName:
		prompt = fmt.Sprintf("Enter your name:\n %s", m.name)
	case promptForPassword:
		maskedPassword := strings.Repeat("*", len(m.password))
		prompt = fmt.Sprintf("Enter your password:\n %s", maskedPassword)
	}

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error:", err)
	}

	box := lipgloss.NewStyle().
		Padding(3, 2).
		MarginTop(height / 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		AlignHorizontal(lipgloss.Center).
		Width(width / 3).
		AlignVertical(lipgloss.Center).
		Render(prompt)

	centering := lipgloss.NewStyle().
		Width(width).
		AlignHorizontal(lipgloss.Center).
		Render(box)

	instructions := lipgloss.NewStyle().
		Width(width).
		PaddingTop(2).
		Align(lipgloss.Center).
		Render("\n(Press 'Enter' to submit, 'Backspace' to delete)")

	return centering + instructions + strings.Repeat("\n", 1)
}

func (m model) headerView() string {
	title := titleStyle.Render("Feederizer")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))

	helpView := m.help.View(m.keys)
	helpHeight := lipgloss.Height(helpView)
	padding := strings.Repeat("\n", helpHeight+1)
	return padding + lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	resp, err := http.Get("http://localhost:8080/feed")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return
	}

	initialModel := model{
		content: string(prettyJSON.String()),

		help: help.New(),
		keys: keys,
	}

	p := tea.NewProgram(
		initialModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
