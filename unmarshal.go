package hiccup

import (
	"io"
	"mime"
	"net/http"
)

/*
BodyDecoder defines an interface to describe different [Unmarshaler]
functions for different content types.
*/
type BodyDecoder interface {
	ContentType() string
	Unmarshal(data []byte, v any) error
}

/*
Unmarshaler function to decode a [http.Request] body.
Functions intended to unmarshal request content must implement this
interface to use it in a [BodyDecoder] or [RequestUnmarshaler].
*/
type Unmarshaler func(data []byte, v any) error

/*
RequestDecoder provides methods to decode request body content of different
content types.

See the [Decoder] function for more info.
*/
type RequestDecoder struct {
	decoder        map[string]BodyDecoder
	defaultDecoder BodyDecoder
}

/*
requestUnmarshaler is a helper struct to easily define an object which
conforms to the [BodyDecoder] interface.
*/
type requestUnmarshaler struct {
	contentType string
	unmarshaler Unmarshaler
}

func (r *requestUnmarshaler) ContentType() string {
	return r.contentType
}

func (r *requestUnmarshaler) Unmarshal(data []byte, v any) error {
	return r.unmarshaler(data, v)
}

/*
WithDecoder is a helper function that returns an object which conforms
to the [BodyDecoder] interface.
*/
func WithDecoder(contentType string, u Unmarshaler) *requestUnmarshaler {
	return &requestUnmarshaler{
		contentType: contentType,
		unmarshaler: u,
	}
}

/*
Decoder returns a [RequestDecoder] configured with the passed BodyDecoders. If no
BodyDecoders are configured then calls to DecodeBody will return the body content
as a string. The first [BodyDecoder] passed will be used as the default decoder if
no "Content-Type" header is sent, or if the value cannot be matched with a
[BodyDecoder].
*/
func Decoder(d ...BodyDecoder) *RequestDecoder {
	dec := new(RequestDecoder)
	if len(d) > 0 {
		dec.defaultDecoder = d[0]
	}

	dec.decoder = make(map[string]BodyDecoder)
	for _, v := range d {
		dec.decoder[v.ContentType()] = v
	}
	return dec
}

/*
DecodeBody with the matched BodyDecoder for the specified "Content-Type" header value
sent in the [http.Request]. If no match is found the default decoder will be used.

It requires the request object be passed, as well as a pointer to the object the
request body will be unmarshaled to.

It returns the raw bytes of the request body, as well as any error if one was
encountered during unmarshaling.
If no decoders are configured the passed value will not be modified, and only
the raw bytes of the request body will be returned.
If the request is nil, or if the body is empty, it returns a nil byte array and
a nil error.
*/
func (r *RequestDecoder) DecodeBody(req *http.Request, v any) ([]byte, error) {
	if req == nil {
		return nil, nil
	}

	b, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	if len(b) == 0 {
		return nil, nil
	}

	contype, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if decFunc := r.decoder[contype]; decFunc != nil {
		return b, decFunc.Unmarshal(b, v)
	} else if decFunc := r.defaultDecoder; decFunc != nil {
		return b, decFunc.Unmarshal(b, v)
	} else {
		return b, nil
	}
}
