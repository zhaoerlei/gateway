package gateway

import (
	"fmt"
	"gate/wrap"
	"log"
	"time"
)

type conState struct {
	socket *Socket
	online bool
}

type ConAdapter interface {
	Read() ([]byte, error)
	Write([]byte) error
	Close()
	RemoteIp() string
}

type readData struct {
	dataArr []byte
	err     error
}

type Socket struct {
	SocketId    string          
	ticker      *time.Ticker    
	adapter     ConAdapter      
	ch          chan []byte     
	ClientRooms map[string]bool 
	router      *Router         
	SocketState bool
	RemoteIp    string
}

func (socket *Socket) write(chState chan<- *conState) {
	for {
		select {
		case b := <-socket.ch:
			if b == nil { 
				socket.ticker.Stop()
				socket.adapter.Close() 
				return
			}
			socket.innerWrite(chState, b)
		case <-socket.ticker.C:
			socket.innerWrite(chState, socket.ping())
		}
	}
}

func (socket *Socket) innerWrite(chState chan<- *conState, b []byte) {
	if err := socket.adapter.Write(b); err != nil {
		chState <- &conState{socket: socket, online: false} 
	}
}

func (socket *Socket) read(chState chan<- *conState, chUp chan<- *Packet) {
	noReadTicker := time.NewTicker(time.Duration(3) * time.Second)
	chReadData := make(chan *readData)
	go readFromClient(socket, chReadData)
	for {
		if !socket.SocketState {
			log.Println("read socketState:", socket.SocketState)
		}
		select {
		case rData := <-chReadData:
			if rData.err != nil {
				log.Println("read err", rData.err)
				chState <- &conState{socket: socket, online: false}
				noReadTicker.Stop()
				return
			}
			chUp <- &Packet{SocketId: socket.SocketId, Sock: socket, DataArr: rData.dataArr}
		case <-noReadTicker.C:
		}
	}
}

func readFromClient(socket *Socket, chReadData chan<- *readData) {
	for {
		b, err := socket.adapter.Read()
		chReadData <- &readData{dataArr: b, err: err}
		if err != nil {
			break
		}
	}
}

func (socket *Socket) close() {
	close(socket.ch)
}

func (socket *Socket) ping() []byte {
	cmd := wrap.CreateCmd(wrap.Cmd_PING, "", Tick())
	b, err := wrap.Boxing(cmd)
	if err != nil {
	}
	return b
}

func (socket *Socket) sendHandshake(chSocket chan *Socket) {
	if !socket.SocketState {
		return
	}
	cmd := wrap.CreateCmd(wrap.Cmd_HANDSHAKE, fmt.Sprintf(socket.SocketId), Tick())
	b, err := wrap.Boxing(cmd)
	socket.ch <- b
	chSocket <- socket
}

func Tick() int64 {
	return time.Now().UnixNano() / 1e6
}
