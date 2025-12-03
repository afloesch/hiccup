package hiccup_test

import (
	"net/http"
	"testing"

	"github.com/afloesch/hiccup"
)

func ExampleRespond() {
	hiccup.Handler(func(r *http.Request) *hiccup.Response {
		// return a Response object and hiccup will
		// handle writing the response.
		return hiccup.Respond(200).
			SetBody("Great Success!").
			SetHeader("Cookie", "key=value;")
	})
}

func ExampleRespond_setRedirectURI() {
	hiccup.Handler(func(r *http.Request) *hiccup.Response {
		// return a Response object with a valid 3XX status code
		// and set the redirect URI.
		return hiccup.Respond(http.StatusTemporaryRedirect).
			SetRedirectURI("http://www.acme.com")
	})
}

func TestResponse(t *testing.T) {
	r := &hiccup.Response{}
	r.SetHeader("key", "value")
	if r.Headers["key"] != "value" {
		t.Error("header value not set")
		t.FailNow()
	}

	r = &hiccup.Response{}
	r.SetHeaders(map[string]string{
		"key": "value",
	})
	if r.Headers["key"] != "value" {
		t.Error("header value not set")
		t.FailNow()
	}
}
