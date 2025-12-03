package hiccup_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afloesch/hiccup"
	"go.yaml.in/yaml/v3"
)

func Example_handler() {
	// create request body decoders
	decoders := hiccup.Decoder(
		hiccup.WithDecoder("application/json", json.Unmarshal), // the first entry is the default
		hiccup.WithDecoder("application/yaml", yaml.Unmarshal),
	)

	// define the supported response body formats.
	encoders := hiccup.Encoder(
		hiccup.WithEncoder("application/json", json.Marshal), // the first entry is the default
		hiccup.WithEncoder("application/yaml", yaml.Marshal),
	)

	// create a handler that conforms to the hiccup.HandlerFunc interface.
	myHandler := func(r *http.Request) *hiccup.Response {
		// decode the request body
		var data = make(map[string]string)
		decoders.DecodeBody(r, &data)

		// do handler stuff...
		fmt.Println(data)

		// send response
		return hiccup.Respond(http.StatusOK).SetBody(map[string]string{
			"Message": "Hello World!",
		})
	}

	// create the http.Handler that can be served or integrated with a router.
	hiccup.Handler(myHandler, encoders...)
}

func ExampleHandler_basic() {
	// create a handler that conforms to the hiccup.HandlerFunc interface.
	myHandler := func(r *http.Request) *hiccup.Response {
		return hiccup.Respond(http.StatusOK).SetBody("Hello World!")
	}

	// make a simple test request.
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// create the http.Handler and make the test request.
	hiccup.Handler(myHandler).ServeHTTP(w, req)

	// print the response body
	body, _ := io.ReadAll(w.Result().Body)
	fmt.Println(string(body))
	// Output: Hello World!
}

func ExampleHandler_differentContentTypes() {
	// create a handler that conforms to the hiccup.HandlerFunc interface.
	myHandler := func(r *http.Request) *hiccup.Response {
		return hiccup.Respond(http.StatusOK).SetBody(map[string]string{
			"Message": "Hello World!",
		})
	}

	// define the supported response formats.
	// this example demonstrates json, yaml, and plain text.
	en := hiccup.Encoder(
		hiccup.WithEncoder("application/json", json.Marshal), // the first entry is the default
		hiccup.WithEncoder("application/yaml", yaml.Marshal),
		hiccup.WithEncoder("text/plain", hiccup.MarshalText),
	)

	// make a simple test request asking for a yaml response.
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/yaml")
	w := httptest.NewRecorder()

	// create the http.Handler and make the test request.
	hiccup.Handler(myHandler, en...).ServeHTTP(w, req)

	// print the response content type header and body
	body, _ := io.ReadAll(w.Result().Body)
	fmt.Println(w.Result().Header.Get("Content-Type"), string(body))
	// Output: application/yaml Message: Hello World!
}

func testRequest(method string, url string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(method, url, body)
	w := httptest.NewRecorder()
	return w, req
}

func testFailMarshal(v any) ([]byte, error) {
	return nil, errors.New("marshal failed")
}

func TestHandler_Redirect(t *testing.T) {
	myHandler := func(r *http.Request) *hiccup.Response {
		return hiccup.
			Respond(http.StatusMovedPermanently).
			SetRedirectURI("http://www.acme.com/home")
	}

	handler := hiccup.Handler(myHandler, hiccup.WithEncoder("text/plain", hiccup.MarshalText))
	w, req := testRequest("GET", "/", nil)
	handler.ServeHTTP(w, req)

	loc := w.Result().Header.Get("Location")
	if loc != "http://www.acme.com/home" {
		t.Error("invalid redirect location", loc)
	}
}

func TestHandler(t *testing.T) {
	myHandler := func(r *http.Request) *hiccup.Response {
		return hiccup.Respond(http.StatusOK).SetBody(map[string]string{
			"Message": "Hello World!",
		}).SetHeader("Cookie", "key=value;").SetHeader("x-test", "value")
	}

	// define the supported response formats.
	// this example demonstrates json, yaml, and plain text.
	en := []hiccup.ResponseEncoder{
		hiccup.WithEncoder("application/json", json.Marshal),
		hiccup.WithEncoder("application/yaml", yaml.Marshal),
		hiccup.WithEncoder("text/plain", hiccup.MarshalText),
		hiccup.WithEncoder("test/failed", testFailMarshal),
	}

	handler := hiccup.Handler(myHandler, en...)

	w, req := testRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/yaml")
	handler.ServeHTTP(w, req)
	body, _ := io.ReadAll(w.Result().Body)
	if string(body) != "Message: Hello World!\n" {
		t.Error(string(body))
		t.FailNow()
	}

	w, req = testRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	handler.ServeHTTP(w, req)
	body, _ = io.ReadAll(w.Result().Body)
	if string(body) != `{"Message":"Hello World!"}` {
		t.Error("incorrect body encoding", body)
		t.FailNow()
	}

	w, req = testRequest("GET", "/", nil)
	handler.ServeHTTP(w, req)
	body, _ = io.ReadAll(w.Result().Body)
	if string(body) != `{"Message":"Hello World!"}` {
		t.Error("incorrect body encoding", body)
		t.FailNow()
	}

	w, req = testRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/plain")
	handler.ServeHTTP(w, req)
	body, _ = io.ReadAll(w.Result().Body)
	if string(body) != `map[Message:Hello World!]` {
		t.Error("incorrect body encoding", string(body))
		t.FailNow()
	}

	w, req = testRequest("GET", "/", nil)
	handler.ServeHTTP(w, req)
	body, _ = io.ReadAll(w.Result().Body)
	if string(body) != `{"Message":"Hello World!"}` {
		t.Error("incorrect body encoding", body)
		t.FailNow()
	}

	w, req = testRequest("GET", "/", nil)
	req.Header.Set("Accept", "test/failed")
	handler.ServeHTTP(w, req)
	body, _ = io.ReadAll(w.Result().Body)
	if string(body) != `marshal failed` {
		t.Error("expected an error from marshaler")
		t.FailNow()
	}

	handler = hiccup.Handler(myHandler)
	w, req = testRequest("GET", "/", nil)
	handler.ServeHTTP(w, req)
	body, _ = io.ReadAll(w.Result().Body)
	if w.Header().Get("Cookie") != "key=value;" {
		t.Error("header value not set")
		t.FailNow()
	}
	if string(body) != `map[Message:Hello World!]` {
		t.Error("incorrect body value")
		t.FailNow()
	}
}
