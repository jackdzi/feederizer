package theme

import (
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

var (
	text_color_string   string
	border_color_string string
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
}
