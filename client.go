package engineio

import (
	"crypto/tls"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/mushroomsir/engine.io/transports"
)

// Options ...
type Options struct {
	// TLSClientConfig specifies the TLS configuration to use with tls.Client.
	// If nil, the default configuration is used.
	TLSClientConfig *tls.Config
}

// NewClient ...
func NewClient(urlStr string, opts ...*Options) (client *transports.Client, err error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("transport", "websocket")
	query.Set("EIO", "3")
	u.RawQuery = query.Encode()

	dialer := &websocket.Dialer{}
	if len(opts) > 0 {
		dialer.TLSClientConfig = opts[0].TLSClientConfig
	}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	client = transports.NewWebSocket(conn)
	return
}
