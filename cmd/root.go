package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "trello-cli",
	Short: "A light cli wrapper for the Trello API",
	Long: `This tool enables AI to easily interact with the Trello API,
allowing for seamless integration and automation of Trello tasks. It provides
a simple interface for managing boards, lists, cards, and more, making it
easier to organize and track your projects.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("profile", "p", "default", "The Trello account to use")
}
