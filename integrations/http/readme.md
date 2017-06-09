# HTTP
Various ways to serve http requests in go.

## simple.go
Let's write a simple server and client using what is provided to us in [`net/http`](https://golang.org/pkg/net/http/).

Point of interest: it is recommended to [not use http.Client in production](https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779).

## sling.go
[Sling](https://github.com/dghubble/sling) is a Go HTTP client library for creating and sending API requests.


For this example, we'll use sling to make some requests to httpbin.


