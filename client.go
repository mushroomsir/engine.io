package engineio

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/mushroomsir/engine.io/transports"
)

// Options ...
type Options struct {
	// TLSClientConfig specifies the TLS configuration to use with tls.Client.
	// If nil, the default configuration is used.
	TLSClientConfig *tls.Config
	LocalAddr       string
	RequestHeader   http.Header
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

	var header http.Header
	if len(opts) > 0 {
		op := opts[0]
		if op.TLSClientConfig != nil {
			dialer.TLSClientConfig = op.TLSClientConfig
		}
		if op.LocalAddr != "" {
			dialer.NetDial = func(network, addr string) (net.Conn, error) {
				localAddr, err := net.ResolveIPAddr("ip", op.LocalAddr)
				if err != nil {
					panic(err)
				}
				localTCPAddr := net.TCPAddr{
					IP: localAddr.IP,
				}
				netDialer := &net.Dialer{LocalAddr: &localTCPAddr}
				return netDialer.Dial(network, addr)
			}
		}
		if len(op.RequestHeader) > 0 {
			header = op.RequestHeader
		}
	}
	conn, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		return
	}
	client = transports.NewWebSocket(conn)
	return
}
