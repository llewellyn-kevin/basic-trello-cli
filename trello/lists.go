package trello

import (
	"fmt"
)

func (c TrelloClient) GetList(listID string) (*[]Card, error) {
	var list []Card
	request := fmt.Sprintf("/1/lists/%s/cards", listID)
	err := c.get(&list, request, map[string]string{})
	if err != nil {
		return nil, err
	}
	return &list, nil
}
