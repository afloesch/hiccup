package hiccup_test

import (
	"encoding/json"
	"fmt"

	"github.com/afloesch/hiccup"
)

func ExampleWithEncoder() {
	m := hiccup.WithEncoder("application/json", json.Marshal)

	b, _ := m.Marshal(map[string]string{
		"Message": "Hello World!",
	})

	fmt.Println(string(b))
	// Output: {"Message":"Hello World!"}
}
