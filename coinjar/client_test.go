package coinjar

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertRequestUsesApiKey(t, r, "someapikey")
		if url := r.URL.Path; url == "/account.json" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprint(w, `
				{
					"user": {
						"available_balance": "1.0",
						"email": "test@example.com",
						"full_name": "John Doe",
						"unconfirmed_balance": "0.3",
						"uuid": "29d7f276-ba50-11e3-b016-7eddf9792095"
					}
				}
			`)
		} else {
			t.Errorf("Requested unexpected endpoint: %v", url)
		}
	}))
	defer ts.Close()

	client := NewCustomClient("someapikey", ts.URL)
	user, err := client.Account()
	assertNotNil(t, err)
	assertEqual(t, user.UUID, "29d7f276-ba50-11e3-b016-7eddf9792095")
	assertEqual(t, user.Email, "test@example.com")
	assertEqual(t, user.FullName, "John Doe")
	assertEqual(t, user.AvailableBalance, "1.0")
	assertEqual(t, user.UnconfirmedBalance, "0.3")
}

func assertRequestUsesApiKey(t *testing.T, r *http.Request, key string) {
	if !strings.HasPrefix(r.Header.Get("Authorization"), "Basic") {
		t.Error("Not using Basic Authentication")
		return
	}
	assertEqual(t, parseApiKey(r), key)
}

func parseApiKey(r *http.Request) string {
	val := r.Header.Get("Authorization")
	parts := strings.SplitN(val, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return ""
	}
	return strings.SplitN(string(data), ":", 2)[0]
}

func assertNotNil(t *testing.T, actual interface{}) {
	if actual != nil {
		t.Errorf("Assertion 'notNil' failed\n\tActual: %v", actual)
	}
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("Assertion 'equal' failed\n\tActual: %v\n\tExpected: %v", actual, expected)
	}
}
