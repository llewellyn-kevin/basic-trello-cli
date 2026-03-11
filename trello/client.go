package trello

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const baseURL = "https://api.trello.com"

type TrelloClient struct {
	client *http.Client
	cmdCtx *cobra.Command
}

type TrelloAccountConfig struct {
	Name  string `yaml:"name"`
	Key   string `yaml:"key"`
	Token string `yaml:"token"`
}

type TrelloAccountConfigurationList struct {
	Accounts []TrelloAccountConfig `yaml:"configurations"`
}

func NewTrelloClient(cmdCtx *cobra.Command) *TrelloClient {
	return &TrelloClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		cmdCtx: cmdCtx,
	}
}

func (c TrelloClient) makeRequest(marshalTarget any, method string, route string, query map[string]string) error {
	req, _ := http.NewRequest(method, baseURL+route, nil)
	req.Header.Set("Accept", "application/json")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	// TODO: Check a list of valid config paths (e.g. XDG_CONFIG_HOME, etc.) instead of hardcoding this path
	configPath := filepath.Join(homeDir, ".config", "trello-cli", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var configList TrelloAccountConfigurationList
	if err = yaml.Unmarshal(data, &configList); err != nil {
		return err
	}

	profile, err := c.cmdCtx.Flags().GetString("profile")
	if err != nil {
		return err
	}
	account, err := configList.getAccountByName(profile)
	if err != nil {
		return err
	}

	queryValues := req.URL.Query()
	queryValues.Set("key", account.Key)
	queryValues.Set("token", account.Token)
	for key, value := range query {
		queryValues.Set(key, value)
	}
	req.URL.RawQuery = queryValues.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		err = fmt.Errorf(`Trello API request failed
Made request to: %s
Response status: %s
Response body: %s`, req.URL.String(), resp.Status, string(body))
		return err
	}
	if err != nil {
		return err
	}

	json.Unmarshal(body, marshalTarget)

	return nil
}

func (c TrelloClient) get(marshalTarget any, route string, query map[string]string) error {
	return c.makeRequest(marshalTarget, http.MethodGet, route, query)
}

func (c TrelloAccountConfigurationList) getAccountByName(name string) (*TrelloAccountConfig, error) {
	for _, account := range c.Accounts {
		if account.Name == name {
			return &account, nil
		}
	}
	return nil, fmt.Errorf("account with name '%s' not found", name)
}
