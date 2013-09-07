package coinjar

import (
	"encoding/json"
	"errors"
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

type BitcoinAddress struct {
	Label          string
	TotalConfirmed string `json:"total_confirmed"`
	TotalReceived  string `json:"total_received"`
	Address        string
}

type bitcoinAddresses struct {
	Addresses []BitcoinAddress `json:"bitcoin_addresses"`
}

func (c *client) BitcoinAddresses() ([]BitcoinAddress, error) {
	body, err := c.read("bitcoin_addresses.json")
	if err != nil {
		return nil, err
	}
	bitcoinAddresses := new(bitcoinAddresses)
	err = json.Unmarshal(body, bitcoinAddresses)
	if err != nil {
		return nil, err
	}
	return bitcoinAddresses.Addresses, nil
}

func (c *client) BitcoinAddress(address string) (obj *BitcoinAddress, err error) {
	body, err := c.read("bitcoin_addresses/" + address + ".json")
	if err != nil {
		return
	}
	if string(body) == "null" {
		return nil, errors.New("Bitcoin address not found")
	}

	var wrapper struct {
		Address *BitcoinAddress `json:"bitcoin_address"`
	}
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return
	}
	return wrapper.Address, nil
}


type FairRate struct {
	Bid  string
	Ask  string
	Spot string
}

func (c *client) FairRate(currency string) (*FairRate, error) {
	body, err := c.read("fair_rate/" + currency + ".json")
	if err != nil {
		return nil, err
	}
	fairRate := new(FairRate)
	err = json.Unmarshal(body, fairRate)
	if err != nil {
		return nil, err
	}
	return fairRate, nil
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
