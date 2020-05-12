package socketBox

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)



type Connection struct {
	conn      *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan []byte

	key string

	closeLocker *sync.Mutex
	isClose     bool
}

func NewConnection(wsConn *websocket.Conn, key string) (connection *Connection, err error) {
	connection = &Connection{
		conn:        wsConn,
		inChan:      make(chan []byte, 1000),
		outChan:     make(chan []byte, 1000),
		closeChan:   make(chan []byte, 1),
		key:         key,
		closeLocker: &sync.Mutex{},
		isClose:     false,
	}
	onlineConnections[key] = connection

	go connection.readLoop()
	go connection.writeLoop()
	return
}

//listenning
func (c *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = c.conn.ReadMessage(); err != nil {
			goto ERR
		}

		select {
		case c.inChan <- data:
		case <-c.closeChan:
			goto ERR
		}
	}

ERR:
	c.Close()
}

func (c *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-c.outChan:
		case <-c.closeChan:
			goto ERR
		}
		//sendTo

		if err = c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	c.Close()
}

func (c *Connection) WriteMsg(data []byte) (err error) {
	select {
	case  c.outChan<-data:
	case <-c.closeChan:
		err = errors.New("Connection is closed!")
	}
	return
}

func (c *Connection) ReadMsg() (data []byte, err error) {
	select {
	case data = <-c.inChan:
	case <-c.closeChan:
		err = errors.New("Connection is closed!")
	}
	return
}

func (c *Connection) Close() {
	defer c.closeLocker.Unlock()
	c.closeLocker.Lock()
	if !c.isClose {
		c.conn.Close()
		delete(onlineConnections, c.key)
		close(c.closeChan)
		c.isClose = true
	}
}
