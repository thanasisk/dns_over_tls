package main

import "log"
import "crypto/tls"
import "bufio"
import "net"

type outgoingConnection struct {
	endpoint   string
	connection *tls.Conn
}

func (conn *outgoingConnection) SetEndpoint(endp string) {
	conn.endpoint = endp
}

func (conn *outgoingConnection) New(error) {

}

func (conn *outgoingConnection) Close() {
	conn.Close()
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
		log.Println(string(netData))
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
		log.Println(string(bytesRead))
		log.Println(string(foo))
		conn.Close()
		c.Write(foo)
	}
	c.Close()
}
