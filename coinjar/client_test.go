package coinjar

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO check auth header
		if url := r.URL.Path; url == "/account.json" {
			// TODO read this from a file
			fmt.Fprint(w, `{"user":{"uuid":"29d7f276-ba50-11e3-b016-7eddf9792095","email":"test@example.com","full_name":"John Doe","available_balance":"1.0","unconfirmed_balance":"0.3"}}`)
		} else {
			t.Errorf("Requested unexpected endpoint: %v", url)
		}
	}))
	defer ts.Close()

	client := NewCustomClient("someapikey", ts.URL)
	user, err := client.Account()
	notNil(t, err)
	equal(t, user.UUID, "29d7f276-ba50-11e3-b016-7eddf9792095")
	equal(t, user.Email, "test@example.com")
	equal(t, user.FullName, "John Doe")
	equal(t, user.AvailableBalance, "1.0")
	equal(t, user.UnconfirmedBalance, "0.3")
}

func notNil(t *testing.T, actual interface{}) {
	if actual != nil {
		t.Errorf("Assertion 'notNil' failed\n\tActual: %v", actual)
	}
}

func equal(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("Assertion 'equal' failed\n\tActual: %v\n\tExpected: %v", actual, expected)
	}
}