package cmd

import (
	"fmt"
	"trello-cli/trello"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get the cards in a list with a given ID",
	Long: `Get the cards in a list with a given ID. For example:

$ trello-cli list 123-list-id

- Card 0: Card Name
  ID: 123-abc
  URL: https://trello.com/c/123-abc/card-name

- Card 1: Second Card
  ID: 456-xyz
  URL: https://trello.com/c/456-xyz/second-card
`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := trello.NewTrelloClient(cmd).GetList(args[0])
		if err != nil {
			fmt.Println("Error fetching list:", err)
			return
		}
		for i, card := range *list {
			fmt.Println(fmt.Sprintf("- Card %d: %s\n  ID: %s\n  URL: %s\n", i+1, card.Name, card.Id, card.Url))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
