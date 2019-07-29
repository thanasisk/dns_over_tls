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
		toCF, err := env.pool.Get()
		log.Println("toCF acquired")
		if err != nil {
			log.Println("Unable to get connection from pool: " + err.Error())
			break
		}
		toCF.Write(netData)
		log.Println("toCF wrote")
		foo := make([]byte, 512)
		bytesRead, err := toCF.Read(foo)
		if err != nil {
			log.Println("Error Reading from TLS endpoint: " + err.Error())
			break
		}
		log.Println("toCF Read")
		toCF.Close()
		log.Println(string(bytesRead))
		log.Println(string(foo))
		c.Write(foo)
	}
	c.Close()
}
