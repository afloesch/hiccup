package hiccup_test

import (
	"net/http"
	"testing"

	"github.com/WP-beta/be-core/src/api/hiccup"
)

func ExampleRespond() {
	hiccup.Handler(func(r *http.Request) *hiccup.Response {
		// return a Response object and hiccup will
		// handle writing the response.
		res := hiccup.Respond(200).
			SetBody("Great Success!").
			SetHeader("Cookie", "key=value;")
		return res
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
