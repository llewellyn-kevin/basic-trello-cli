package trello

import (
	"fmt"
)

type Card struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"desc"`
	Checklists  []Checklist `json:"checklists"`
	Url         string      `json:"url"`
}

type Checklist struct {
	Id    string          `json:"id"`
	Name  string          `json:"name"`
	Items []ChecklistItem `json:"checkItems"`
}
type ChecklistItem struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

func (c TrelloClient) GetCard(cardID string) (*Card, error) {
	var card Card
	request := fmt.Sprintf("/1/cards/%s", cardID)
	err := c.get(&card, request, map[string]string{
		"fields":     "name,desc,url",
		"checklists": "all",
	})
	if err != nil {
		return nil, err
	}
	return &card, nil
}
