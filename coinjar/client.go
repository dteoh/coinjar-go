package coinjar

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

type User struct {
	UUID               string
	Email              string
	FullName           string `json:"full_name"`
	AvailableBalance   string `json:"available_balance"`
	UnconfirmedBalance string `json:"unconfirmed_balance"`
}

func (c *client) Account() (obj *User, err error) {
	body, err := c.read("account.json")
	if err != nil {
		return
	}

	var wrapper struct{ User *User }
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return
	}
	return wrapper.User, nil
}

type BitcoinAddress struct {
	Label          string
	TotalConfirmed string `json:"total_confirmed"`
	TotalReceived  string `json:"total_received"`
	Address        string
}

func (c *client) BitcoinAddresses() ([]BitcoinAddress, error) {
	return c.ListBitcoinAddresses(100, 0)
}

func (c *client) ListBitcoinAddresses(limit, offset int) (obj []BitcoinAddress, err error) {
	body, err := c.read("bitcoin_addresses.json",
		"limit", strconv.Itoa(limit),
		"offset", strconv.Itoa(offset))
	if err != nil {
		return
	}

	var wrapper struct {
		Addresses []BitcoinAddress `json:"bitcoin_addresses"`
	}
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return
	}
	return wrapper.Addresses, nil
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

type Transaction struct {
	Confirmations      string
	Status             string
	Amount             string
	Reference          string
	CounterpartyType   string `json:"counterparty_type"`
	UpdatedAt          string `json:"updated_at"`
	UUID               string
	BitcoinTxid        string `json:"bitcoin_txid"`
	RelatedPaymentUUID string `json:"related_payment_uuid"`
	CounterpartyName   string `json:"counterparty_name"`
	CreatedAt          string `json:"created_at"`
}

func (c *client) Transactions() ([]Transaction, error) {
	return c.ListTransactions(100, 0)
}

func (c *client) ListTransactions(limit, offset int) (obj []Transaction, err error) {
	body, err := c.read("transactions.json",
		"limit", strconv.Itoa(limit),
		"offset", strconv.Itoa(offset))
	if err != nil {
		return
	}

	var wrapper struct{ Transactions []Transaction }
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return
	}
	return wrapper.Transactions, nil
}

func (c *client) Transaction(uuid string) (obj *Transaction, err error) {
	body, err := c.read("transactions/" + uuid + ".json")
	if err != nil {
		return
	}
	if strings.Contains(string(body), "\"status\":\"404\"") {
		return nil, errors.New(fmt.Sprintf("Transaction not found, response body: %q", body))
	}

	var wrapper struct{ Transaction *Transaction }
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return
	}
	return wrapper.Transaction, nil
}

type FairRate struct {
	Bid  string
	Ask  string
	Spot string
}

func (c *client) FairRate(currency string) (obj *FairRate, err error) {
	body, err := c.read("fair_rate/" + currency + ".json")
	if err != nil {
		return
	}
	obj = new(FairRate)
	err = json.Unmarshal(body, obj)
	if err != nil {
		return
	}
	return
}

func (c *client) read(api string, params ...string) (body []byte, err error) {
	request, _ := http.NewRequest("GET", c.endpoint+"/"+api, nil)
	request.SetBasicAuth(c.apiKey, "")
	request.URL.RawQuery = createQuery(params)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func createQuery(params []string) string {
	plen := len(params)
	if plen%2 == 1 {
		plen = plen - 1
	}
	values := url.Values{}
	for i := 0; i < plen; i += 2 {
		values.Set(params[i], params[i+1])
	}
	return values.Encode()
}
