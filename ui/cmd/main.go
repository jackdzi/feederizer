package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// You generally won't need this unless you're processing stuff with
// complicated ANSI escape sequences. Turn it on if you notice flickering.
//
// Also keep in mind that high performance rendering only works for programs
// that use the full size of the terminal. We're enabling that below with
// tea.EnterAltScreen().
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
		switch m.mode {
		case normalMode:
			k := msg.String()
			switch k {
			case "ctrl+c", "q", "esc":
				return m, tea.Quit
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
		case inputMode:
			switch msg.String() {
			case "backspace":
				if len(m.input) > 0 && m.cursorPos > 0 {
					m.input = m.input[:m.cursorPos-1] + m.input[m.cursorPos:]
					m.cursorPos--
				}
			case "left":
				if m.cursorPos > 0 {
					m.cursorPos--
				}
			case "right":
				if m.cursorPos < len(m.input) {
					m.cursorPos++
				}
			case "enter":
				// Submit the input and return to normal mode
				data := map[string]string{
					"key1": m.input,
					"key2": m.input,
				}
				jsonData, err := json.Marshal(data)
				if err != nil {
					fmt.Println("Error parsing string to JSON")
				}
				if _, err := http.Post("http://localhost:8080/insert", "application/json", bytes.NewBuffer(jsonData)); err != nil {
					fmt.Println("Error:", err)
				}
				m.mode = normalMode
				m.input = ""
				m.cursorPos = 0
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
			case "esc":
				// Exit input mode without submitting
				m.mode = normalMode
				m.input = ""
				m.cursorPos = 0
			default:
				// Add typed character
				m.input = m.input[:m.cursorPos] + msg.String() + m.input[m.cursorPos:]
				m.cursorPos++
			}
		}

	// Handle window size messages
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
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := titleStyle.Render("Feederizer")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
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

	p := tea.NewProgram(
		model{content: string(prettyJSON.String())},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
