package main

import (
	"fmt"
	"os"
  "time"

	"github.com/jackdzi/feederizer/ui/internal/driver"
	"github.com/jackdzi/feederizer/ui/internal/theme"
  //"github.com/jackdzi/feederizer/server/shared"
	"github.com/spf13/cobra"
)

var (
	version = "beta"

	rootCmd = &cobra.Command{
		Use:   "feederizer",
		Short: "Feederizer is a command line UI tool for aggregating and viewing RSS feeds.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				os.Exit(1)
			}
      // TODO: Check if user is logged in already and then skip login page to go right to the user's feed. Will have to implement currentPage in driver.New(). Same with concurrent server running if docker == false
      //
      //go shared.RunServer()
      time.Sleep(5 * time.Millisecond)

      program := driver.New(theme.NewStyles())
      program.Run()
      program.ReleaseTerminal()
		},
	}
)

func Init() {
	rootCmd.Version = version

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:     "help",
		Aliases: []string{"h", "--h"},
		Short:   "Show help for feederizer",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Root().Help()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "about",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("\nFeederizer was created to hopefully serve\nas a tool that makes viewing news that you\nwant to see easy, fast, and customizable.")
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "example",
		Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("Example")
		},
	})
}

func main() {
	Init()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
