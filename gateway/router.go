package gateway

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

const (
	ROOT_ROOM = "root"
)

type Router struct {
	sync.Mutex

	chConState chan *conState 
	chUp       chan *Packet   
	chDown     chan *Packet   
	rooms map[string]map[string]*Socket 

	handler          Handler
	RoomUserCountMap map[string]int 
}

type Packet struct {
	SocketId string
	DataArr  []byte
	Sock     *Socket
}

func (r *Router) init() {
	rootRoom := r.rooms[ROOT_ROOM]
	if rootRoom == nil {
		rootRoom = make(map[string]*Socket)
		r.rooms[ROOT_ROOM] = rootRoom
	}
}

func (r *Router) runUp() {
	for {
		select {
		case state := <-r.chConState:
			go r.handler.OnChange(state.socket, state.online)

		case pkt := <-r.chUp:
			go r.handler.OnReceive(pkt)
		}
	}
}

func (r *Router) runDown() {
	for {
		select {
		case pkt := <-r.chDown:
			func() {
				if pkt.Sock.SocketState {
					chanLen := len(pkt.Sock.ch)
					pkt.Sock.ch <- pkt.DataArr
				}
			}()

		}
	}
}

func (r *Router) online(socket *Socket) {
	rootRoom := r.rooms[ROOT_ROOM]
	if rootRoom == nil {
		rootRoom = make(map[string]*Socket)
		rootRoom[socket.SocketId] = socket
	} else {
		exists := false
		if c := rootRoom[socket.SocketId]; c != nil && c == socket {
			exists = true
		}
		if !exists {
			rootRoom[socket.SocketId] = socket
		}
	}
	r.rooms[ROOT_ROOM] = rootRoom        
	socket.ClientRooms[ROOT_ROOM] = true 
}

func (r *Router) offline(socket *Socket) {
	for roomName, _ := range socket.ClientRooms {
		room := r.rooms[roomName]
		if room == nil {
			continue
		}
		delete(room, socket.SocketId)

		if len(room) == 0 && roomName != ROOT_ROOM {
			delete(r.rooms, roomName)
		}
	}
	if socket.SocketState {
		socket.close()
		socket.SocketState = false
	}
}

func (r *Router) DisconnectSocket(socketId string) error {
	rootRoom := r.rooms[ROOT_ROOM]
	socket := rootRoom[socketId]
	r.chConState <- &conState{socket: socket, online: false}
	return nil
}

func (r *Router) JoinRoom(socketId string, roomId string) error {
	r.Lock()

	defer r.Unlock()

	rootRoom := r.rooms[ROOT_ROOM]
	socket := rootRoom[socketId]
	room := r.rooms[roomId]
	if room == nil {
		room = make(map[string]*Socket)
		r.rooms[roomId] = room
	}
	room[socketId] = socket
	socket.ClientRooms[roomId] = true
	return nil
}

func (r *Router) SendToRoom(roomId string, exceptId string, dataArr []byte) {
	room := r.rooms[roomId]
	if room == nil {
		return
	}
	for keySocketId, valSocket := range room {
		if keySocketId == exceptId {
			continue
		}
		r.chDown <- &Packet{SocketId: keySocketId, Sock: valSocket, DataArr: dataArr}
	}

}

func (r *Router) SendToClient(socketId string, dataArr []byte) {
	rootRoom := r.rooms[ROOT_ROOM]
	r.chDown <- &Packet{SocketId: socketId, Sock: socket, DataArr: dataArr}
}

func (r *Router) SendToGroup(roomId string, socketIdGroup string, dataArr []byte) {
	room := r.rooms[roomId]
	socketIdArr := strings.Split(socketIdGroup, ",")
	if room == nil || socketIdGroup == "" {
		return
	}
	for i := 0; i < len(socketIdArr); i++ {
		socketId := socketIdArr[i]
		socket := room[socketId]
		if socket == nil {
			return
		}
		r.chDown <- &Packet{SocketId: socketId, Sock: socket, DataArr: dataArr}
	}

}

func (r *Router) GetRooms() map[string]map[string]*Socket {
	return r.rooms
}
