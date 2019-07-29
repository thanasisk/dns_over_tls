# dns_over_tls
a DNS over TLS proxy written in Golang
# build locally
## Prerequisites
```go get github.com/jessevdk/go-flags```
## Building
```go build```
## cleaning
```go clean```
# Docker
## build using Docker
```docker build . -t n26_test```
if this is for production usage, tag accordingly and DO NOT use *latest* (EVER!)

As you can see, this is a 2 stage build, in order to keep the running image small
# security considerations
## SUID and Linux capabilities
Normally, to bind to ports below 1024, you need root rights. In the event of a 
compromise, this would give the attacker full root access on the machine, something that 
is clearly not acceptable from a security perspective.However, recent
Docker versions support Linux capabilities ```man 7 capabilities```.
Setting the right capability to *BOTH* the executable and the process, do the following

- ```setcap 'cap_net_bind_service=+ep'```
- ```docker run -it --cap-add NET_BIND_SERVICE --network host n26_test /go/proxy -a=192.168.178.118```
For convenience I used host networking:
By default the proxy will attempt to bind 0.0.0.0:53 - on Fedora this is taken
by *dnsmasq* so I specified eth0's address explicitly
## Rate limiting
Rate limiting (i.e. using a leaky bucket algorithm) is currently not implemented, which
creates an assymetric attack surface - an attacker can send cheap DNS queries over TCP,
each of which will trigger a TLS handshake, which can be considered expensive. Therefore,
a nice reliability and security pattern is rate limiting and using a queue. Understandably, given
the short timeframe from this exercise, this is NOT implemented.
# Future improvements
## TLS connections
as of now, there is a 1:1 mapping of incoming to proxied connections.
Given that TLS handshake can be considered expensive, this is clearly *inefficient*.
Therefore, an optimization that could be made is to use TLS sessions to reuse
TLS connections as much as possible, reducing the load and round-about time per 
incoming connection. Since, it was not clear to which extend I could use external
components, some leftovers can be found under ```vendor```
## Multithreading
Furthermore, the server implementation is naive. Given that this is an I/O bounded problem,
standard Golang CSP measures such as channels, can be used for a more elegant solution. Dropping down
to low level synchronization level might gain a performance boost, but again, given that this is an
I/O bounded problem, this optimization might not be essentials.


# Known bugs
- verbosity is ignored as of now
- ```bufio.read``` error message is kind of spammy
# Debugging
In case you want to debug, send a SIGUSR1 to the running process. It will drop
to stdout certain runtime/debug statistics
# Further questions
- Imagine this proxy being deployed in an infrastructure. What would be the security concerns you would raise?
Please see section security concerns above.
- How would you integrate that solution in a distributed, microservices-oriented and containerized architecture?
TBC
-  What other improvements do you think would be interesting to add to the project?
Please see section Future Improvements above.
