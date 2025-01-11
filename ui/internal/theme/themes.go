package theme

import (
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

var (
	text_color_string       string
	border_color_string     string
	title_color_string      string
	title_text_color_string string
	text_highlight_string   string
)

func init() {
	config, err := os.ReadFile("../config.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	tree, err := toml.Load(string(config))
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	border_color_string = tree.Get("theme.borderColor").(string)
	text_color_string = tree.Get("theme.textColor").(string)
  title_text_color_string = tree.Get("theme.titleTextColor").(string)
  title_color_string = tree.Get("theme.titleColor").(string)
  text_color_string = tree.Get("theme.textColor").(string)
  text_highlight_string = tree.Get("theme.textHighlight").(string)
}
