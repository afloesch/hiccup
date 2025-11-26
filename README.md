# Hiccup
![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)

Hiccup is a go package to simplify handling http requests and responses with multiple content types in a RESTful API, and eliminate some of the boilerplate code in writing responses and decoding request body content. The main handler method in hiccup returns a standard go library [http.Handler](https://pkg.go.dev/net/http#Handler), so it can easily integrate with any go router or framework of your choosing.

It supports any data format that can conform to the go standard library [json.Marshal]https://pkg.go.dev/encoding/json#Marshal) and [json.Unmarshal](https://pkg.go.dev/encoding/json#Unmarshal) function parameters for different request or response body content types, so virtually any content type can be supported with little code or impact to individual handlers. Request body content is decoded based on the standard http "Content-Type" header value sent in the request, and response body content is encoded based on the standard http "Accept" header value sent in the request. If either of these headers is not sent in a request, or if the request headers cannot be matched with what is configured, hiccup will fallback to the configured defaults.

```sh
go get -u github.com/afloesch/hiccup
```