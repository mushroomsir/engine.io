package engineio

import (
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/mushroomsir/engine.io/transports"
)

// Options ...
type Options struct {
}

// NewClient ...
func NewClient(urlStr string, opts ...Options) (client *transports.Client, err error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("transport", "websocket")
	query.Set("EIO", "3")
	u.RawQuery = query.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		return
	}
	client = transports.NewWebSocket(conn)
	return
}
