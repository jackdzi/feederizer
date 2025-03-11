package config

import (
	"github.com/pelletier/go-toml"
  "log"
  "os"
)


func ReturnConfig() *toml.Tree {
	config, err := os.ReadFile("../config.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	tree, err := toml.Load(string(config))
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
  return tree
}

func GetFilePath() string {
	return "../config.toml"
}

func LoadFile() (*toml.Tree, error){
  return toml.LoadFile(GetFilePath())
}
