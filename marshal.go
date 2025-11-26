package hiccup

import "fmt"

/*
Marshaler function to encode a http response body.
Functions intended for marshaling response body content
must implement this interface.
*/
type Marshaler func(v any) ([]byte, error)

/*
Helper struct to quickly define a ResponseEncoder.

See the [ResponseMarshaler] function.
*/
type responseMarshaler struct {
	contentType string
	marshaler   Marshaler
}

/*
ResponseMarshaler is a helper function to return an object that conforms
to the [ResponseEncoder] interface, and can be passed into the [hiccup.HTTPHandler]
function to support multiple content encodings based on the "Accept" header
sent in a http request.
*/
func ResponseMarshaler(contentType string, m Marshaler) *responseMarshaler {
	return &responseMarshaler{
		contentType: contentType,
		marshaler:   m,
	}
}

/*
MarshalText is a utility which implements a [Marshaler] interface function
to support plain text responses.
*/
func MarshalText(v any) ([]byte, error) {
	s := fmt.Sprint(v)
	return []byte(s), nil
}

func (r *responseMarshaler) ContentType() string {
	return r.contentType
}

func (r *responseMarshaler) Marshal(v any) ([]byte, error) {
	return r.marshaler(v)
}
