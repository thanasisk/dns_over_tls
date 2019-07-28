package main

import "log"
import "bufio"
import "net"

//import "io"

//import "syscall"

func handleConnection(c net.Conn, env Env) {
	log.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		// TODO of course FIX this
		netData := make([]byte, 512)
		_, err := bufio.NewReader(c).Read(netData)
		if err != nil {
			log.Println("bufio.Read(): " + err.Error())
			break
		}
		log.Println(string(netData))
		// now that we have netData, let's send them to cloudflare
		env.oConn.Write(netData)
		foo := make([]byte, 512)
		bytesRead, err := env.oConn.Read(foo)
		if err != nil {
			log.Println("Error Reading from TLS endpoint: " + err.Error())
			break
		}
		log.Println(string(bytesRead))
		log.Println(string(foo))
		c.Write(foo)
	}
	c.Close()
}
