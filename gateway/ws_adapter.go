package gateway

import (
	"golang.org/x/net/websocket"
	"log"
)

type wsAdapter struct {
	conn *websocket.Conn
}

func (p *wsAdapter) Read() ([]byte, error) {
	var b []byte
	err := websocket.Message.Receive(p.conn, &b)
	return b, err
}

func (p *wsAdapter) Write(b []byte) error {
	return websocket.Message.Send(p.conn, b)
}

func (p *wsAdapter) Close() {
	p.conn.Close()
}
func (p *wsAdapter) RemoteIp() string {
	return ""
}
