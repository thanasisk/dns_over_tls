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
```docker build . -t n26_test```

as you can see, this is a 2 stage build, in order to keep the running image small
# run using Docker
For convenience I used host networking:
```docker run -it --network host n26_test /go/proxy -a=192.168.178.118```
By default the proxy will attempt to bind 0.0.0.0:53 - on Fedora this is taken
by *dnsmasq* so I specified eth0's address explicitly
# Debugging
In case you want to debug, send a SIGUSR1 to the running process. It will drop
to stdout certain runtime/debug statistics
# security considerations
TBC
# Future improvements
as of now, there is a 1:1 mapping of incoming to proxied connections.
Given that TLS handshake can be considered expensive, this is clearly *inefficient*.
Therefore, an optimization that could be made is to use TLS sessions to reuse
TLS connections as much as possible, reducing the load and round-about time per 
incoming connection. Since, it was not clear to which extend I could use external
components, some leftovers can be found under ```vendor```
# Known bugs
- verbosity is ignored as of now
- ```bufio.read``` error message is kind of spammy
