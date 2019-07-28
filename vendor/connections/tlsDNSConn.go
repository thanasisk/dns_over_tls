package tlsDNSConn

import "crypto/tls"

type OutgoingConnection struct {
	Connection *tls.Conn
	Verbosity  bool
}

func (conn *OutgoingConnection) SetConnection(Connection *tls.Conn) {
	conn.Connection = Connection
}

func NewConnection(endpoint string) (*OutgoingConnection, error) {
	c, err := tls.Dial("tcp4", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var p *OutgoingConnection
	p.SetConnection(c)
	return p, nil
}

func (conn *OutgoingConnection) New() *OutgoingConnection {
	return conn
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
