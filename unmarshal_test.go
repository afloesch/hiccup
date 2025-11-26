package hiccup_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/WP-beta/be-core/src/api/hiccup"
	"go.yaml.in/yaml/v3"
)

func ExampleRequestUnmarshaler() {
	// example unmarshaler
	ex := hiccup.RequestUnmarshaler("application/json", json.Unmarshal)

	// test json object
	msg := []byte(`{
		"Message": "Hello World!"
	}`)

	// unmarshal json to a map
	d := make(map[string]string)
	ex.Unmarshal(msg, &d)
	fmt.Println(d["Message"])
	// Output: Hello World!
}

func ExampleRequestDecoder_DecodeBody() {
	// create request body decoders
	dec := hiccup.Decoder(
		hiccup.RequestUnmarshaler("application/json", json.Unmarshal), // the first entry is the default
		hiccup.RequestUnmarshaler("application/yaml", yaml.Unmarshal),
	)

	// create test request of yaml data
	req := httptest.NewRequest("GET", "/", bytes.NewBufferString(`Message: Hello World!`))
	req.Header.Set("Content-Type", "application/yaml")

	// decode the request body to a map
	var data = make(map[string]string)
	dec.DecodeBody(req, &data)

	fmt.Println(data["Message"])
	// Output: Hello World!
}

type testReader struct{}

func (t *testReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test reader error")
}

func TestRequestDecoder(t *testing.T) {
	dec := hiccup.Decoder(
		hiccup.RequestUnmarshaler("application/json", json.Unmarshal),
		hiccup.RequestUnmarshaler("application/yaml", yaml.Unmarshal),
	)

	_, req := testRequest("GET", "/", nil)
	var b []byte
	_, err := dec.DecodeBody(req, &b)
	if err != nil || b != nil {
		t.Error(err, b)
		t.FailNow()
	}

	_, req = testRequest("GET", "/", bytes.NewBufferString(`{"Message": "Hello World!"}`))

	data := make(map[string]string)
	_, err = dec.DecodeBody(req, &data)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if data["Message"] != "Hello World!" {
		t.Error("unexpected decode value")
		t.FailNow()
	}

	_, req = testRequest("GET", "/", bytes.NewBufferString(`Hello World!`))
	req.Header.Set("Content-Type", "application/yaml")

	data = make(map[string]string)
	_, err = dec.DecodeBody(req, &data)
	if err == nil {
		t.Error("expected a decoding error", data)
		t.FailNow()
	}

	dec = hiccup.Decoder()
	_, req = testRequest("GET", "/", bytes.NewBufferString("Hello World!"))
	data = make(map[string]string)
	b, err = dec.DecodeBody(req, &data)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if string(b) != "Hello World!" {
		t.Error("incorrect body", string(b))
		t.FailNow()
	}

	_, req = testRequest("GET", "/", new(testReader))
	req.Header.Set("Content-Type", "application/yaml")

	data = make(map[string]string)
	_, err = dec.DecodeBody(req, &data)
	if err == nil {
		t.Error("expected a reader error", data)
		t.FailNow()
	}

	data = make(map[string]string)
	b, err = dec.DecodeBody(nil, &data)
	if err != nil || len(b) != 0 {
		t.Error("expected empty params for empty request")
		t.FailNow()
	}
}
