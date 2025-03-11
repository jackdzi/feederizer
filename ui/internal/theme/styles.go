package theme

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	border_color = lipgloss.Color(border_color_string)
	text_color   = lipgloss.Color(text_color_string)
)

type Styles struct {
	Input         InputStyles
	Footer        lipgloss.Style
	Centering     lipgloss.Style
	Homepage      Homepage
	Feed          Feed
	Title         lipgloss.Style
	TextColor     string
	TextHighlight string
	Viewer        Viewport
}

type Feed struct {
	BottomPadding lipgloss.Style
	TableWidth    int
	TableHeight   int
}

type Viewport struct {
	Height int
	Width  int
}

type InputStyles struct {
	Box             lipgloss.Style
	Instructions    lipgloss.Style
	InstructionText string
	BottomPadding   lipgloss.Style
}

type Homepage struct {
	Box           lipgloss.Style
	BottomPadding lipgloss.Style
}

func (s *Styles) ApplySizes() {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error:", err)
	}
	var width int
	if w/3 < 38 {
		width = 38
	} else {
		width = w / 3
	}
	heightLeft := h - 2 - 20
	offset := 0
	if heightLeft%2 == 1 {
		offset = 1
	}
	s.Viewer.Height = h - 8
	s.Viewer.Width = w
	s.Input.Box = s.Input.Box.Width(width).MarginTop(heightLeft/2 + 3 + offset)
	s.Homepage.Box = s.Input.Box.Width(width).MarginTop(heightLeft/2 + 1 + offset)
	s.Centering = s.Centering.Width(w)
	s.Input.Instructions = s.Input.Instructions.Width(w)
	s.Input.BottomPadding = s.Input.BottomPadding.Margin(0, 0, heightLeft/2, 0)
	s.Homepage.BottomPadding = s.Input.BottomPadding.Margin(0, 0, heightLeft/2+2, 0)
	s.Feed.BottomPadding = s.Feed.BottomPadding.MarginBottom(heightLeft / 2)
	s.Feed.TableWidth = w - 8
	s.Feed.TableHeight = h - 7
}

func (s *Styles) RenderInstructions() string {
	return s.Input.Instructions.Render(s.Input.InstructionText)
}

func NewStyles() Styles {
	s := Styles{}
	s.Feed.TableHeight = 1
	s.Feed.TableWidth = 1
	s.Feed.BottomPadding = lipgloss.NewStyle().Margin(0, 0, 1, 0)
	s.TextColor = text_color_string
	s.TextHighlight = text_highlight_string
	s.Title = lipgloss.NewStyle().Margin(0).Width(12).
		Background(lipgloss.Color(title_color_string)).
		Foreground(lipgloss.Color(title_text_color_string))
	s.Footer = lipgloss.NewStyle().Margin(1, 0, 0, 1)
	s.Input.Box = lipgloss.NewStyle().
		Padding(3, 2).
		PaddingBottom(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(border_color_string)).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	s.Centering = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center)
	s.Input.Instructions = lipgloss.NewStyle().
		PaddingTop(2).
		Align(lipgloss.Center)
	s.Input.InstructionText = "\n(Press 'Enter' to submit, 'Backspace' to delete, 'Tab' to switch)"
	s.Input.BottomPadding = lipgloss.NewStyle().Margin(0, 0, 1, 0)
	s.Homepage.Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(border_color_string)).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	return s
}
