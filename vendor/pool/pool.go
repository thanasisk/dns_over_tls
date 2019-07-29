package tlsPool

// heavily influenced from https://github.com/eternnoir/gncp/blob/master/pool.go
// influenced because original code deals with plain TCP connections ...
import "errors"
import "crypto/tls"
import "sync"
import "context"
import "time"
import "fmt"

type ConnPool interface {
	Get() (tls.Conn, error)
	Close() error
	Remove(conn tls.Conn) error
	GetWithTimeoute(timeoute time.Duration) (tls.Conn, error)
}

type GTLSPool struct {
	lock         sync.Mutex
	conns        chan tls.Conn
	minConnNum   int
	maxConnNum   int
	totalConnNum int
	closed       bool
	connCreator  func() (tls.Conn, error)
}

var (
	errPoolIsClose = errors.New("Connection pool has been closed")
	// Error for get connection time out.
	errTimeOut      = errors.New("Get Connection timeout")
	errContextClose = errors.New("Get Connection close by context")
)

func NewPool(minConn int, maxConn int, connCreator func() (tls.Conn, error)) (*GTLSPool, error) {
	if minConn > maxConn || minConn < 0 || maxConn <= 0 {
		return nil, errors.New("Conn # bounding error")
	}
	pool := &GTLSPool{}
	pool.minConnNum = minConn
	pool.maxConnNum = maxConn
	pool.connCreator = connCreator
	pool.conns = make(chan tls.Conn, maxConn)
	pool.closed = false
	pool.totalConnNum = 0
	err := pool.init()
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (p *GTLSPool) init() error {
	for i := 0; i < p.minConnNum; i++ {
		conn, err := p.createConn()
		if err != nil {
			return err
		}
		p.conns <- conn
	}
	return nil
}

func (p *GTLSPool) Get() (tls.Conn, error) {
	if p.isClosed() == true {
		return tls.Conn{}, errPoolIsClose
	}
	go func() {
		conn, err := p.createConn()
		if err != nil {
			return
		}
		p.conns <- conn
	}()
	select {
	case conn := <-p.conns:
		return p.packConn(conn), nil
	}
}

func (p *GTLSPool) GetWithTimeOut(timeout time.Duration) (tls.Conn, error) {
	if p.isClosed() == true {
		return tls.Conn{}, errPoolIsClose
	}
	go func() {
		conn, err := p.createConn()
		if err != nil {
			return
		}
		p.conns <- conn
	}()
	select {
	case conn := <-p.conns:
		return p.packConn(conn), nil
	case <-time.After(timeout):
		return tls.Conn{}, errTimeOut
	}
}

func (p *GTLSPool) GetWithContext(ctx context.Context) (tls.Conn, error) {
	if p.isClosed() == true {
		return tls.Conn{}, errPoolIsClose
	}
	go func() {
		conn, err := p.createConn()
		if err != nil {
			return
		}
		p.conns <- conn
	}()
	select {
	case conn := <-p.conns:
		return p.packConn(conn), nil
	case <-ctx.Done():
		return tls.Conn{}, errContextClose
	}
}
func (p *GTLSPool) Close() error {
	if p.isClosed() == true {
		return errPoolIsClose
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.closed = true
	close(p.conns)
	for conn := range p.conns {
		conn.Close()
	}
	return nil
}

func (p *GTLSPool) Put(conn tls.Conn) error {
	if p.isClosed() == true {
		return errPoolIsClose
	}
	if conn == nil {
		p.lock.Lock()
		p.totalConnNum = p.totalConnNum - 1
		p.lock.Unlock()
		return errors.New("Cannot put nil to connection pool")
	}
	select {
	case p.conns <- conn:
		return nil
	default:
		return conn.Close()
	}
}

func (p *GTLSPool) isClosed() bool {
	p.lock.Lock()
	ret := p.closed
	p.lock.Unlock()
	return ret
}

func (p *GTLSPool) Remove(conn tls.Conn) error {
	if p.isClosed() == true {
		return errPoolIsClose
	}
	p.lock.Lock()
	p.totalConnNum = p.totalConnNum - 1
	p.lock.Unlock()
	conn.Close()
	return nil
}

func (p *GTLSPool) createConn() (tls.Conn, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.totalConnNum >= p.maxConnNum {
		return tls.Conn{}, fmt.Errorf("Cannot create new connection: Current: %d Max: %d", p.totalConnNum, p.maxConnNum)
	}
	conn, err := p.connCreator()
	if err != nil {
		return tls.Conn{}, fmt.Errorf("Cannot create new connection: %s", err)
	}
	p.totalConnNum = p.totalConnNum + 1
	return conn, nil
}

func (p *GTLSPool) packConn(conn tls.Conn) tls.Conn {
	ret := &CpConn{pool: p}
	ret.Conn = conn
	return ret
}

//
type CpConn struct {
	Conn tls.Conn
	pool *GTLSPool
}

// Destroy will close connection and release connection from connection pool.
func (conn *CpConn) Destroy() error {
	if conn.pool == nil {
		return errors.New("Connection not belong any connection pool.")
	}
	err := conn.pool.Remove(conn.Conn)
	if err != nil {
		return err
	}
	conn.pool = nil
	return nil
}

// Close will push connection back to connection pool. It will not close the real connection.
func (conn *CpConn) Close() error {
	if conn.pool == nil {
		return errors.New("Connection not belong any connection pool.")
	}
	return conn.pool.Put(conn.Conn)
}
