package transports

import (
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
)

// Event ...
type Event struct {
	Type string
	Data []byte
}

// Client ...
type Client struct {
	Conn     *websocket.Conn
	sender   chan *Packet
	receiver chan *Packet
	Event    chan *Event
	opened   bool
	done     chan bool
	info     *connectionInfo
}
type connectionInfo struct {
	Sid          string        `json:"sid"`
	Upgrades     []string      `json:"upgrades"`
	PingInterval time.Duration `json:"pingInterval"`
	PingTimeout  time.Duration `json:"pingTimeout"`
}

// NewWebSocket ...
func NewWebSocket(conn *websocket.Conn) *Client {
	client := &Client{
		Conn:     conn,
		opened:   false,
		sender:   make(chan *Packet, 100),
		receiver: make(chan *Packet, 100),
		Event:    make(chan *Event, 100),
	}
	client.done = make(chan bool)
	go client.msgReceiveLoop()
	go client.msgReadLoop()
	return client
}
func (c *Client) msgReceiveLoop() {
	// listens for open event, and then make it open
	for {
		_, res, err := c.Conn.ReadMessage()
		if nil != err {
			c.receiver <- NewErrorPacket(err)
			c.receiver <- NewClosePacket()
			break
		} else {
			packet := BytesToPacket(res)
			c.receiver <- packet
		}
	}
}
func (c *Client) msgReadLoop() {
	for {
		select {
		case packet := <-c.receiver:
			event := &Event{Data: packet.Data}
			switch packet.Type {
			case Open:
				c.opened = true
				event.Type = "open"
				c.info = &connectionInfo{}
				err := json.Unmarshal(event.Data, c.info)
				if err != nil {
					event.Type = "error"
				}
			case Close:
				c.opened = false
				event.Type = "close"
				c.Event <- event
				return
			case Upgrade:
				event.Type = "upgrade"
			case Message:
				event.Type = "message"
			case Ping:
				event.Type = "ping"
				c.sender <- &Packet{Type: Pong, Data: packet.Data}
			case Pong:
				event.Type = "pong"
			case Error:
				event.Type = "error"
			default:
				event = nil
			}
			if nil != event {
				c.Event <- event
			}
		case p := <-c.sender:
			if err := c.Conn.WriteMessage(websocket.TextMessage, PacketToBytes(p)); nil != err {
				c.receiver <- NewErrorPacket(err)
			}
		case <-c.done:
			return
		}
	}
}

// SendMessage ...
func (c *Client) SendMessage(data []byte) {
	c.sender <- NewPacket(Message, data)
}

// SendPacket ...
func (c *Client) SendPacket(packet *Packet) {
	c.sender <- packet
}

// GetSID ...
func (c *Client) GetSID() string {
	return c.info.Sid
}

// Close ...
func (c *Client) Close() error {
	close(c.done)
	return c.Conn.Close()
}
