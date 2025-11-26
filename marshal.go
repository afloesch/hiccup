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

See the [WithEncoder] function.
*/
type responseMarshaler struct {
	contentType string
	marshaler   Marshaler
}

/*
Encoder is a helper that returns a slice of ResponseEncoders for
inclusion in a [hiccup.Handler].
*/
func Encoder(e ...ResponseEncoder) []ResponseEncoder {
	enc := make([]ResponseEncoder, 0)
	enc = append(enc, e...)
	return enc
}

/*
WithEncoder is a helper function to return an object that conforms
to the [ResponseEncoder] interface, and can be passed into the [hiccup.HTTPHandler]
function to support multiple content encodings based on the "Accept" header
sent in a http request.
*/
func WithEncoder(contentType string, m Marshaler) *responseMarshaler {
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
