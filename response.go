package hiccup

/*
ResponseEncoder defines an interface to describe different marshalers
for different content types.

See also the [ResponseMarshaler] function.
*/
type ResponseEncoder interface {
	// The content type to respond with.
	ContentType() string
	// Marshal function to encode response body content.
	Marshal(v any) ([]byte, error)
}

/*
Response object returned by a [Handler] function.
*/
type Response struct {
	// Response body content.
	Body any
	// Response headers to set.
	Headers map[string]string
	// HTTP status code to send.
	StatusCode int
}

/*
Respond returns a [Response] object for a [Handler] function return value.
*/
func Respond(statusCode int) *Response {
	return &Response{
		StatusCode: statusCode,
	}
}

/*
Set the response body to the passed value.
*/
func (r *Response) SetBody(value any) *Response {
	r.Body = value
	return r
}

/*
Set a header value. Any existing value will be overwritten.
*/
func (r *Response) SetHeader(key string, value string) *Response {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[key] = value
	return r
}

/*
Set all header values at once. Overwrites all existing header values.
*/
func (r *Response) SetHeaders(headers map[string]string) *Response {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers = headers
	return r
}
