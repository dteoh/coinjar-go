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

type Client struct {
	apiKey     string
	endpoint   string
	httpClient *http.Client
}

func NewClient(apiKey string) (c *Client) {
	c = new(Client)
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

func (c *Client) Account() (obj *User, err error) {
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

func (c *Client) BitcoinAddresses() ([]BitcoinAddress, error) {
	return c.ListBitcoinAddresses(100, 0)
}

func (c *Client) ListBitcoinAddresses(limit, offset int) (obj []BitcoinAddress, err error) {
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

func (c *Client) BitcoinAddress(address string) (obj *BitcoinAddress, err error) {
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

type Contact struct {
	UpdatedAt string `json:"updated_at"`
	UUID      string
	Name      string
	PayeeName string
	PayeeType string
	CreatedAt string
}

func (c *Client) Contacts() ([]Contact, error) {
	return c.ListContacts(100, 0)
}

func (c *Client) ListContacts(limit, offset int) (obj []Contact, err error) {
	body, err := c.read("contacts.json",
		"limit", strconv.Itoa(limit),
		"offset", strconv.Itoa(offset))
	if err != nil {
		return
	}

	var wrapper struct{ Contacts []Contact }
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return
	}
	return wrapper.Contacts, nil
}

func (c *Client) Contact(uuid string) (obj *Contact, err error) {
	body, err := c.read("contacts/" + uuid + ".json")
	if err != nil {
		return
	}

	if string(body) == "null" {
		return nil, errors.New("Contact not found")
	}

	var wrapper struct{ Contact *Contact }
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return
	}
	return wrapper.Contact, nil
}

type Payment struct {
	Status             string
	Amount             string
	CreatedAt          string `json:"created_at"`
	PayeeName          string `json:"payee_name"`
	PayeeType          string `json"payee_type"`
	Reference          string
	RelatedTransaction *Transaction `json:"related_transaction"`
	UUID               string
	UpdatedAt          string `json:"updated_at"`
}

func (c *Client) Payments() (obj []Payment, err error) {
	return c.ListPayments(100, 0)
}

func (c *Client) ListPayments(limit, offset int) (obj []Payment, err error) {
	body, err := c.read("payments.json",
		"limit", strconv.Itoa(limit),
		"offset", strconv.Itoa(offset))
	if err != nil {
		return
	}

	var wrapper struct{ Payments []Payment }
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return
	}
	return wrapper.Payments, nil
}

func (c *Client) Payment(uuid string) (obj *Payment, err error) {
	body, err := c.read("payments/" + uuid + ".json")
	if err != nil {
		return
	}
	if string(body) == "null" {
		return nil, errors.New("Payment not found")
	}

	var wrapper struct{ Payment *Payment }
	err = json.Unmarshal(body, &wrapper)
	if err != nil {
		return
	}
	return wrapper.Payment, nil
}

type Transaction struct {
	Amount              string
	BitcoinTxid         string `json:"bitcoin_txid"`
	Confirmations       string
	CounterpartyAddress string `json:"counterparty_address"`
	CounterpartyName    string `json:"counterparty_name"`
	CounterpartyType    string `json:"counterparty_type"`
	CounterpartyUserID  int    `json:"counterparty_user_id"`
	CreatedAt           string `json:"created_at"`
	ID                  int
	PaymentID           int `json:"payment_id"`
	Reference           string
	RelatedPaymentUUID  string `json:"related_payment_uuid"`
	Status              string
	UUID                string
	UpdatedAt           string `json:"updated_at"`
	UserID              int    `json:"user_id"`
}

func (c *Client) Transactions() ([]Transaction, error) {
	return c.ListTransactions(100, 0)
}

func (c *Client) ListTransactions(limit, offset int) (obj []Transaction, err error) {
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

func (c *Client) Transaction(uuid string) (obj *Transaction, err error) {
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

func (c *Client) FairRate(currency string) (obj *FairRate, err error) {
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

func (c *Client) read(api string, params ...string) (body []byte, err error) {
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
