package theme

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var border_color = lipgloss.Color(border_color_string)
var text_color = lipgloss.Color(text_color_string)

type Styles struct {
	Input     InputStyles
	Footer    lipgloss.Style
	Centering lipgloss.Style
	Homepage  Homepage
  Title     lipgloss.Style
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
	s.Input.Box = s.Input.Box.Width(w / 3).MarginTop(h/3 + 1)
	s.Homepage.Box = s.Input.Box.Width(w / 3).MarginTop(h/3 + 1)
	s.Centering = s.Centering.Width(w)
	s.Input.Instructions = s.Input.Instructions.Width(w)
	s.Input.BottomPadding = s.Input.BottomPadding.Margin(0, 0, h/3+2, 0)
	s.Homepage.BottomPadding = s.Input.BottomPadding.Margin(0, 0, h/3-1, 0)
}

func (s *Styles) RenderInstructions() string {
	return s.Input.Instructions.Render(s.Input.InstructionText)
}

func NewStyles() Styles {
	s := Styles{}
  s.Title = lipgloss.NewStyle().
    Background(lipgloss.Color(border_color_string)).Width(10).
    Foreground(lipgloss.Color(text_color_string))
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
