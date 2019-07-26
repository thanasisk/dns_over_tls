package main

import "os"
import "crypto/tls"
import "log"
import "net"
import "bufio"
import "github.com/jessevdk/go-flags"

func init() {
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")
}

var opts struct {
	Verbose bool   `short:"v" long:"verbose" description:"verbose mode"`
	Port    string `short:"p" long:"port" description:"the port to bind to" default:"53"`
	Address string `short:"a" long:"address" description:"the address to listen to" default:"0.0.0.0"`
	Dns     string `short:"d" long:"dns" description:"the DNS server to connect to" default:"1.1.1.1"`
	Sport   string `short:"s" long:"sport" description:"the remote DNS port to connect to" default:"853"`
}

func handleConnection(c net.Conn) {
	log.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		// TODO of course FIX this
		netData := make([]byte, 512)
		_, err := bufio.NewReader(c).Read(netData)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(netData)
		// now that we have netData, let's send them to cloudflare
		conn, err := tls.Dial("tcp", "1.1.1.1:853", nil)
		if err != nil {
			log.Fatal(err)
		}
		conn.Write(netData)
		foo := make([]byte, 512)
		bytesRead, err := conn.Read(foo)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(bytesRead)
		log.Println(foo)
		conn.Close()
		c.Write(foo)
	}
	c.Close()
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}
	port := opts.Port
	address := opts.Address
	// this can be optimized by using strings.Builder
	addrPort := address + ":" + port
	// fireup our TCP Listener
	l, err := net.Listen("tcp4", addrPort)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	// rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			// TODO: revisit this!
			log.Fatal(err)
		}
		go handleConnection(c)
	}
}
