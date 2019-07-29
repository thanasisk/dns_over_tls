# dns_over_tls
a DNS over TLS proxy written in Golang
# build locally
## Prerequisites
```go get github.com/jessevdk/go-flags```
## Building
```go build```
## cleaning
```go clean```
# build using Docker
TBC
# run using Docker
TBC
# security considerations
TBC
# Future improvements
as of now, there is a 1:1 mapping of incoming to proxied connections.
Given that TLS handshake can be considered expensive, this is clearly *inefficient*.
Therefore, an optimization that could be made is to use TLS sessions to reuse
TLS connections as much as possible, reducing the load and round-about time per 
incoming connection.
