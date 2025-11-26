package hiccup

import (
	"fmt"
	"mime"
	"net/http"
)

/*
Handler function for http requests. Functions intended to be
used in a [hiccup.Handler] must implement this interface.

See also [Response] and [Respond] for returning a response from
a handler.
*/
type HandlerFunc func(r *http.Request) *Response

/*
ResponseHandler unifies http response body encoding for any
[HandlerFunc].

The [Handler] function returns a ResponseHander.
*/
type ResponseHandler struct {
	handler        HandlerFunc
	encoder        map[string]ResponseEncoder
	defaultEncoder ResponseEncoder
}

var contentTypeText = mime.TypeByExtension(".txt")

/*
Handler returns a [http.Handler] for the passed [HandlerFunc],
and any [ResponseEncoder]. The first ResponseEncoder passed
will be used as the default encoder if no match can be made
with the "Accept" header value sent in a [http.Request], or if
no "Accept" header value is sent.

If no ResponseEncoders are passed the body will be sent as
plain text.

See [ResponseHandler.ServeHTTP] for more info.
*/
func Handler(h HandlerFunc, w ...ResponseEncoder) http.Handler {
	rh := new(ResponseHandler)
	rh.handler = h

	if len(w) > 0 {
		rh.defaultEncoder = w[0]
	}

	rh.encoder = make(map[string]ResponseEncoder)
	for _, e := range w {
		rh.encoder[e.ContentType()] = e
	}
	return rh
}

/*
ServeHTTP implements the [http.Handler] interface for handling http
requests.

It encodes the [Handler] returned response content based on the configured
Encoders, and the "Accept" header value sent by a client in a http
request. If no matching encoder for the requested content type is
found then plain text is sent.

If a configured encoder in the [ResponseHandler] cannot successfully
marshal response body content the error encountered will be sent
as plain text with a 500 status code.
*/
func (h *ResponseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := h.handler(r)
	accept, _, _ := mime.ParseMediaType(r.Header.Get("Accept"))

	if encFunc := h.encoder[accept]; encFunc != nil {
		writeEncodedBody(w, res, encFunc)
	} else if encFunc := h.defaultEncoder; encFunc != nil {
		writeEncodedBody(w, res, encFunc)
	} else {
		writeTextBody(w, res)
	}
}

func writeEncodedBody(w http.ResponseWriter, r *Response, enc ResponseEncoder) {
	b, err := enc.Marshal(r.Body)
	if err != nil {
		writeTextBody(w, &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		})
		return
	}

	for k, v := range r.Headers {
		w.Header().Set(k, v)
	}
	w.Header().Set("Content-Type", enc.ContentType())
	w.WriteHeader(r.StatusCode)

	if len(b) > 0 {
		w.Write(b)
	}
}

func writeTextBody(w http.ResponseWriter, r *Response) {
	for k, v := range r.Headers {
		w.Header().Set(k, v)
	}
	w.Header().Set("Content-Type", contentTypeText)
	w.WriteHeader(r.StatusCode)

	if r.Body != nil {
		w.Write([]byte(fmt.Sprint(r.Body)))
	}
}
