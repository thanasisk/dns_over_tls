package main

import "os"
import "log"
import "net"
import "syscall"
import "runtime"
import "runtime/debug"
import "os/signal"
import "encoding/json"
import "github.com/jessevdk/go-flags"

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
	endpoint string
	verbose  bool
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go dumpInfo(c)
	dnsEndpoint := opts.Dns + ":" + opts.Sport
	env := &Env{endpoint: dnsEndpoint, verbose: opts.Verbose}
	log.Println("Firing up TCP4 listener")
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

// the following function catches SIGUSR1 and dumps runtime statistics
// there is no performance penalty for collecting these stats
func dumpInfo(c chan os.Signal) {
	for {
		<-c
		log.Println("Signal caught - dumping runtime stats")
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		s, _ := json.Marshal(m)
		log.Println("MemStats JSON follows")
		log.Printf("%s\n", s)
		var garC debug.GCStats
		debug.ReadGCStats(&garC)
		log.Printf("\nLastGC:\t%s", garC.LastGC)         // time of last collection
		log.Printf("\nNumGC:\t%d", garC.NumGC)           // number of garbage collections
		log.Printf("\nPauseTotal:\t%s", garC.PauseTotal) // total pause for all collections
		log.Printf("\nPause:\t%s", garC.Pause)           // pause history, most recent first
		log.Println("debug.Stack: " + string(debug.Stack()))
		log.Println("runtime.NumGoroutine: " + string(runtime.NumGoroutine()))
	}
}
