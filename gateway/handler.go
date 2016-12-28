package gateway

type MsgHandler struct {
}

type Handler interface {
	OnReceive(pkt *Packet)
	OnChange(socket *Socket, online bool)
}

func (handler *MsgHandler) OnReceive(pkt *Packet) {
}

func (handler *MsgHandler) OnChange(socket *Socket, online bool) {

}

func (handler *MsgHandler) SendToRoom(roomName string, blackId string) {
}



