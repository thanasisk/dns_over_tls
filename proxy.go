package main

import "os"
import "log"
import "net"
import "github.com/jessevdk/go-flags"
import "pool"
import "crypto/tls"

///func init() {
// required for TLS 1.3 support
//	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")
//}

var opts struct {
	Verbose bool   `short:"v" long:"verbose" description:"verbose mode"`
	Port    string `short:"p" long:"port" description:"the port to bind to" default:"53"`
	Address string `short:"a" long:"address" description:"the address to listen to" default:"0.0.0.0"`
	Dns     string `short:"d" long:"dns" description:"the DNS server to connect to" default:"1.1.1.1"`
	Sport   string `short:"s" long:"sport" description:"the remote DNS port to connect to" default:"853"`
}

type Env struct {
	//oConn *tlsDNSConn.OutgoingConnection
	oConn tls.Conn
	pool  *tlsPool.GTLSPool
}

func (e *Env) SetConnection(c *tls.Conn) {
	e.oConn = *c
}

// connCreator let connection know how to create new connection.
func connCreator() (tls.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "1.1.1.1:853")
	if err != nil {
		log.Println("Unable to resolve endpoint:")
		return tls.Conn{}, err
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println("Unable to dial TCP")
		return tls.Conn{}, err
	}
	// upgrade standard TCP connection to insecure TLS one
	err = tcpConn.SetKeepAlive(true)
	if err != nil {
		log.Println("Unable to set KeepAlive on TCP transport")
		return tls.Conn{}, err
	}
	c := tls.Client(tcpConn, &tls.Config{InsecureSkipVerify: true})
	return *c, nil
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
	log.Println(dnsEndpoint)
	// Create new connection pool. It will initialize 3 connection in pool when pool created.
	// If connection not enough in pool, pool will call creator to create new connection.
	// But when total connection number pool created reach 10 connection, pool will not creat
	// any new connection until someone call Remove().
	pool, err := tlsPool.NewPool(3, 10, connCreator)

	// Get connection from pool. If pool has no connection and total connection reach max number
	// of connections, this method will block until someone put back connection to pool.
	oconn, err := pool.Get()
	if err != nil {
		log.Fatal("Error establishing connection to DNS endpoint: " + err.Error())
	}
	env := &Env{oConn: oconn, pool: pool}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal("Accept: " + err.Error())
		}
		go handleConnection(c, *env)
	}
}
