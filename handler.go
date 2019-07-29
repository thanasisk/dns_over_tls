package main

import "log"
import "bufio"
import "net"
import "crypto/tls"

func handleConnection(c net.Conn, env Env) {
	log.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		lenDNSPacket := 1024
		netData := make([]byte, lenDNSPacket)
		cfData := make([]byte, lenDNSPacket)
		_, err := bufio.NewReader(c).Read(netData)
		if err != nil {
			log.Println("bufio.Read(): " + err.Error())
			break
		}
		log.Println(string(netData))
		// now that we have netData, let's send them to endpoint
		log.Println("Connecting to: %s", env.endpoint)
		oc, err := tls.Dial("tcp4", env.endpoint, nil)
		if err != nil {
			log.Println("Cannot Dial %s: %s", env.endpoint, err.Error())
			break
		}
		oc.Write(netData)
		bytesRead, err := oc.Read(cfData)
		if err != nil {
			log.Println("Error Reading from TLS endpoint: " + err.Error())
			break
		}
		log.Println("Bytes read: %d", bytesRead)
		oc.Close()
		log.Println(string(bytesRead))
		log.Println(string(cfData))
		c.Write(cfData)
	}
	c.Close()
}
