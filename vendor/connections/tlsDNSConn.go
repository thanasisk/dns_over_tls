package tlsDNSConn

import "crypto/tls"
import "log"
import "net"

type OutgoingConnection struct {
	Connection *tls.Conn
	Verbosity  bool
}

func NewConnection(endpoint string) (*OutgoingConnection, error) {
	// rewrite using tls.Client
	tcpAddr, err := net.ResolveTCPAddr("tcp", endpoint)
	if err != nil {
		log.Println("Unable to resolve endpoint: " + endpoint)
		return nil, err
	}

	tcpConn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		log.Println("Unable to dial TCP " + endpoint)
		return nil, err
	}
	//func (c *TCPConn) SetKeepAlive(keepalive bool) error
	err = tcpConn.SetKeepAlive(true)
	if err != nil {
		log.Println("Unable to set KeepAlive on TCP transport")
		return nil, err
	}
	c := tls.Client(tcpConn, &tls.Config{InsecureSkipVerify: true})
	p := &OutgoingConnection{Connection: c}
	return p, nil
}

func (conn *OutgoingConnection) Close() {
	conn.Connection.Close()
}

func (conn *OutgoingConnection) Read(buf []byte) (int, error) {
	bytesRead, err := conn.Connection.Read(buf)
	return bytesRead, err
}

func (conn *OutgoingConnection) Write(payload []byte) {
	conn.Connection.Write(payload)
}
