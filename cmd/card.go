package cmd

import (
	"fmt"
	"strings"
	"trello-cli/trello"

	"github.com/spf13/cobra"
)

var cardCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "card <card-id-or-url>",
	Short: "View a card with a given ID or url",
	Long: `View a card with a given ID or url. For example:
$ trello-cli card 123-card-id
ID:    123-card-id
URL:   https://trello.com/c/123-card-id/card-name
Name:  Card Name
Desc:  This is the description of the card.
Checklist: Checklist Name
  - [ ] Checklist Item 1
  - [x] Checklist Item 2
`,
	Run: func(cmd *cobra.Command, args []string) {
		cardId := args[0]
		if cardIdIsUrl(cardId) {
			var err error
			cardId, err = parseCardIdFromUrl(cardId)
			if err != nil {
				fmt.Println("Error parsing card ID from URL:", err)
				return
			}
		}

		card, err := trello.NewTrelloClient(cmd).GetCard(cardId)
		if err != nil {
			fmt.Println("Error fetching card:", err)
			return
		}

		fmt.Println("ID:   ", card.Id)
		fmt.Println("URL:  ", card.Url)
		fmt.Println("Name: ", card.Name)
		fmt.Println("Desc: ", card.Description)
		for _, checklist := range card.Checklists {
			fmt.Printf("Checklist: %s\n", checklist.Name)
			for _, item := range checklist.Items {
				completeChar := ' '
				if item.State == "complete" {
					completeChar = 'x'
				}
				fmt.Printf("  - [%c] %s\n", completeChar, item.Name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cardCmd)
}

func cardIdIsUrl(cardId string) bool {
	return strings.HasPrefix(cardId, "http://") || strings.HasPrefix(cardId, "https://")
}

func parseCardIdFromUrl(url string) (string, error) {
	// We just assume the card url is in the typical format: https://trello.com/c/{cardId}/{cardName}
	// TODO: Create a utility where we can specify url formats and extract variables from them, so we can be more flexible with the url formats we accept
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return "", fmt.Errorf("invalid Trello card URL")
	}
	return parts[4], nil
}
