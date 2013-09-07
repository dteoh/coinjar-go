package coinjar

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type client struct {
	apiKey     string
	endpoint   string
	httpClient *http.Client
}

func NewClient(apiKey string) (c *client) {
	c = new(client)
	c.apiKey = apiKey
	c.endpoint = "https://api.coinjar.io/v1"
	c.httpClient = new(http.Client)
	return
}

type Account struct {
	User struct {
		Uuid               string
		Email              string
		FullName           string `json:"full_name"`
		AvailableBalance   string `json:"available_balance"`
		UnconfirmedBalance string `json:"unconfirmed_balance"`
	}
}

func (c *client) Account() (*Account, error) {
	body, err := c.read("account.json")
	if err != nil {
		return nil, err
	}
	account := new(Account)
	err = json.Unmarshal(body, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (c *client) read(api string) ([]byte, error) {
	request, _ := http.NewRequest("GET", c.endpoint+"/"+api, nil)
	request.SetBasicAuth(c.apiKey, "")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
