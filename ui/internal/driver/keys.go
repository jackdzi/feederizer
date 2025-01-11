package driver

import "github.com/charmbracelet/bubbles/key"

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
	Enter          key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.AddUser, k.InitializeFeed, k.DeleteFeed},
		{k.Help, k.Quit},
	}
}

var keys_config = keyMap{
  Enter: key.NewBinding(
    key.WithKeys("enter"),
    key.WithHelp("↵", "choose entry"),
    ),
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
		key.WithHelp("I", "initialize database"),
	),
	DeleteFeed: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete users"),
	),
	EnterInputMode: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "enter input mode"),
	),
}
