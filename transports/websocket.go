package transports

import "github.com/gorilla/websocket"

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
	go client.msgReceiveLoop()
	go client.msgReadLoop()
	go client.msgWriteLoop()
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
		packet := <-c.receiver
		event := &Event{Data: packet.Data}
		switch packet.Type {
		case Open:
			c.opened = true
			event.Type = "open"
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
			event.Data = []byte("invalid message type")
		default:
			event = nil
		}
		if nil != event {
			c.Event <- event
		}
	}
}
func (c *Client) msgWriteLoop() {
	for {
		p := <-c.sender
		err := c.Conn.WriteMessage(websocket.TextMessage, PacketToBytes(p))
		if nil != err {
			break
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

// Close ...
func (c *Client) Close() error {
	return c.Conn.Close()
}
