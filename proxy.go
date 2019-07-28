package main

import "os"
import "log"
import "net"
import "github.com/jessevdk/go-flags"
import "connections"
import "crypto/tls"

func init() {
	// required for TLS 1.3 support
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")
}

var opts struct {
	Verbose bool   `short:"v" long:"verbose" description:"verbose mode"`
	Port    string `short:"p" long:"port" description:"the port to bind to" default:"53"`
	Address string `short:"a" long:"address" description:"the address to listen to" default:"0.0.0.0"`
	Dns     string `short:"d" long:"dns" description:"the DNS server to connect to" default:"1.1.1.1"`
	Sport   string `short:"s" long:"sport" description:"the remote DNS port to connect to" default:"853"`
}

type Env struct {
	oConn *tlsDNSConn.OutgoingConnection
}

func (e *Env) SetConnection(c *tls.Conn) {
	e.oConn.Connection = c
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}
	// this can be optimized by using strings.Builder
	listenerEndpoint := opts.Address + ":" + opts.Port
	// fireup our TCP Listener
	l, err := net.Listen("tcp4", listenerEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	// time to setup our TCP TLS connection
	dnsEndpoint := opts.Dns + ":" + opts.Sport
	oconn, err := tlsDNSConn.NewConnection(dnsEndpoint)
	if err != nil {
		log.Fatal("Error establishing connection to DNS endpoint: " + err.Error())
	}
	env := &Env{oConn: oconn}
	for {
		c, err := l.Accept()
		if err != nil {
			// TODO: revisit this!
			log.Println("Accept()")
			log.Fatal(err)
		}
		go handleConnection(c, *env)
	}
}
