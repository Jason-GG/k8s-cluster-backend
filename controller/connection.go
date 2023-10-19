package controller

import (
	"errors"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

func websockettest(c *websocket.Conn) {
	// c.Locals is added to the *websocket.Conn
	log.Println(c.Locals("allowed"))  // true
	log.Println(c.Params("id"))       // 123
	log.Println(c.Query("v"))         // 1.0
	log.Println(c.Cookies("session")) // ""

	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}

}

// WsMessage websocket消息
type WsMessage struct {
	MessageType int
	Data        []byte
}

// WsConnection
type WsConnection struct {
	wsSocket *websocket.Conn // basic websocket
	inChan   chan *WsMessage // in queue
	outChan  chan *WsMessage // out queue

	mutex     sync.Mutex // avoid repeatedly close channel
	isClosed  bool
	closeChan chan byte // close channel
}

func (wsConn *WsConnection) wsReadLoop() {
	var (
		msgType int
		data    []byte
		msg     *WsMessage
		err     error
	)
	for {
		// read a message
		if msgType, data, err = wsConn.wsSocket.ReadMessage(); err != nil {
			goto ERROR
		}
		msg = &WsMessage{
			msgType,
			data,
		}
		select {
		case wsConn.inChan <- msg:
		case <-wsConn.closeChan:
			goto CLOSED
		}
	}
ERROR:
	wsConn.WsClose()
CLOSED:
}

// 发送协程
func (wsConn *WsConnection) wsWriteLoop() {
	var (
		msg *WsMessage
		err error
	)
	for {
		select {
		// get a response
		case msg = <-wsConn.outChan:
			// pop it into channel
			if err = wsConn.wsSocket.WriteMessage(msg.MessageType, msg.Data); err != nil {
				goto ERROR
			}
		case <-wsConn.closeChan:
			goto CLOSED
		}
	}
ERROR:
	wsConn.WsClose()
CLOSED:
}

/************** parallel **************/
func InitWebsocket(c *websocket.Conn) (wsConn *WsConnection, err error) {
	wsSocket := c
	wsConn = &WsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *WsMessage, 100000),
		outChan:   make(chan *WsMessage, 100000),
		closeChan: make(chan byte),
		isClosed:  false,
	}

	go wsConn.wsReadLoop()
	go wsConn.wsWriteLoop()

	return
}

func (wsConn *WsConnection) WsWrite(messageType int, data []byte) (err error) {
	select {
	case wsConn.outChan <- &WsMessage{messageType, data}:
	case <-wsConn.closeChan:
		err = errors.New("websocket closed")
	}
	return
}

func (wsConn *WsConnection) WsRead() (msg *WsMessage, err error) {
	select {
	case msg = <-wsConn.inChan:
		return
	case <-wsConn.closeChan:
		err = errors.New("websocket closed")
	}
	return
}

func (wsConn *WsConnection) WsClose() {
	wsConn.wsSocket.Close()

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}
