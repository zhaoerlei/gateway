package main

import (
	"bytes"
	"flag"
	"fmt"
	"gate/gateway"
	"gate/web"
	"gate/wrap"
	"log"
	"net/http"
	_ "net/http/pprof"
)

const (
	SocketPort = ":8114"
	HttpPort   = ":8115"
)

var (
	router        *gateway.Router
	NodeServerUrl = ""
	GoServerID    = ""
)

func main() {

	nodeUrl := flag.String("nodeUrl", "", "eg http://221.122.71.51:8103")
	goId := flag.String("goId", "", "this is Go server id 61.182.132.11 is 11")

	flag.Parse()
	NodeServerUrl = *nodeUrl
	GoServerID = *goId

	if NodeServerUrl == "" || GoServerID == "" {
		log.Println("parameter is invalid")
		return
	}
	gateway.NodeServerUrl = NodeServerUrl
	gateway.GoServerID = GoServerID

	testHandler := &TestHandler{}
	router = gateway.CreateRouter(testHandler)
	go gateway.ReportUserCount(router)

	log.Println("welcome dabaojian", NodeServerUrl, GoServerID)
	gateway.CreateTcpListener(SocketPort, router)
	log.Println("tcp server ok", SocketPort)
	webServer := &web.WebServer{
		Addr:   HttpPort,
		Router: router,
	}
	log.Println("http server ok ", HttpPort)
	webServer.StartWebServer()
}

type TestHandler struct {
}

func (handler *TestHandler) OnReceive(pkt *gateway.Packet) {
	m, err := wrap.Unboxing(pkt.DataArr)
	if err != nil {
		log.Println("unbox error", err)
		return
	}
	switch m.(type) {
	case *wrap.Msg:
		msg := m.(*wrap.Msg)
		clientMsgHandler(msg, pkt)
	case *wrap.Cmd:
		cmd := m.(*wrap.Cmd)
		if cmd.GetTp() == wrap.Cmd_PING {
			log.Println("receive ping:", pkt.SocketId)
			pong := wrap.CreateCmd(wrap.Cmd_PONG, fmt.Sprintf(pkt.SocketId), gateway.Tick())
			b, err := wrap.Boxing(pong)
			if err != nil {
				log.Println("pong unbox:err", err)
				return
			}
			router.SendToClient(pkt.SocketId, b)

		}
	default:
		log.Println("unknown msg type")
	}
}
func (handler *TestHandler) OnChange(socket *gateway.Socket, online bool) {
	for roomName, inRoom := range socket.ClientRooms {
		if roomName != gateway.ROOT_ROOM && inRoom {
			resp, err := http.Get(NodeServerUrl + "/connChg?socketId=" + socket.SocketId + "&roomId=" + roomName + "&goId=" + GoServerID + "&rmIp=" + socket.RemoteIp)
			if err != nil {
				log.Println(err)
				continue
			}
			defer resp.Body.Close()
		}
	}

}

func clientMsgHandler(msg *wrap.Msg, pkt *gateway.Packet) {
	if msg.GetTp() == wrap.Msg_CLIENT_EVENT {
		go forwardNodejs(msg, pkt)
	} else if msg.GetTp() == wrap.Msg_USER_COUNT {
		go sendUserCount(msg, pkt)
	} else {
		log.Println("unknow msg type ")
	}
}

func forwardNodejs(msg *wrap.Msg, pkt *gateway.Packet) {

	msgBuf := bytes.NewBufferString(msg.GetEvent())
	resp, err := http.Post(NodeServerUrl+"/gateway?socketId="+pkt.SocketId+"&goId="+GoServerID+"&rmIp="+pkt.Sock.RemoteIp, "application/json", msgBuf)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
}

func sendUserCount(msg *wrap.Msg, pkt *gateway.Packet) {
	roomName := msg.GetEvent()
	userCount := router.RoomUserCountMap[roomName]
	repMsg := wrap.CreateUserCountMsg(roomName, userCount)
	b, err := wrap.Boxing(repMsg)
	if err != nil {
		log.Println("sendUserCount:err", err)
		return
	}

	router.SendToClient(pkt.SocketId, b)
}
