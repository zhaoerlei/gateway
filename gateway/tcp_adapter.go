package gateway

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const (
	TCP_HEADER_LEN   = 4
	TCP_MAX_PACK_LEN = 4 * 1024 * 1024
)

type tcpAdapter struct {
	conn net.Conn
}

func (p *tcpAdapter) Read() ([]byte, error) {
	//read header
	b := make([]byte, TCP_HEADER_LEN)
	if _, err := io.ReadFull(p.conn, b); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(b)
	var n uint32
	if err := binary.Read(buf, binary.BigEndian, &n); err != nil {
		return nil, err
	}

	if n > TCP_MAX_PACK_LEN {
		return nil, fmt.Errorf("may be error pack length: %dk", n/1024/1024)
	}

	b = make([]byte, n)
	_, err := io.ReadFull(p.conn, b)
	return b, err

}

func (p *tcpAdapter) Write(b []byte) error {
	n := uint32(len(b))

	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, &n); err != nil {
		return err
	}

	if _, err := p.conn.Write(buf.Bytes()); err != nil {
		return err
	}

	_, err := p.conn.Write(b)
	return err
}

func (p *tcpAdapter) Close() {
	p.conn.Close()
}

func (p *tcpAdapter) RemoteIp() string {
	return ""
}
