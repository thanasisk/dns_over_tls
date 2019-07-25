package main

import "os"
import "log"
import "net"
import "strings"
import "bufio"
import "github.com/jessevdk/go-flags"

func init() {
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")
}

var opts struct {
	Verbose bool   `short:"v" long:"verbose" description:"verbose mode"`
	Port    string `short:"p" long:"port" description:"the port to bind to" default:"0.0.0.0:53"`
}

func handleConnection(c net.Conn) {
	log.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
	}
	c.Close()
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}
	port := opts.Port
	// fireup our TCP Listener
	l, err := net.Listen("tcp4", port)
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
