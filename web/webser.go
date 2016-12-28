package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gate/gateway"
	"gate/wrap"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type WebServer struct {
	Addr   string
	Router *gateway.Router
}

type SendMsgObj struct {
	SocketId string 
	RoomId   string
	JsonStr  string
	SendType int 
}

func (webSer *WebServer) StartWebServer() error {
	webSer.webRouter()
	if err := http.ListenAndServe(webSer.Addr, nil); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (webSer *WebServer) webRouter() {
	r := mux.NewRouter()

	xdf := r.PathPrefix("/xdf").Subrouter()
	xdf.HandleFunc("/sendmsg", webSer.sendPostHandler).Methods("POST")
	xdf.HandleFunc("/joinroom", webSer.joinHandler).Methods("GET")
	xdf.HandleFunc("/kickuser", webSer.kickHandler).Methods("GET")
	xdf.HandleFunc("/sendmsg2", webSer.sendGetHandler).Methods("GET")

	ws := r.PathPrefix("/ws").Subrouter()
	ws.Handle("/conn", gateway.CreateWsHandler(webSer.Router))

	util := r.PathPrefix("/util").Subrouter()
	util.HandleFunc("/rooms", webSer.roomsHandler).Methods("GET")
	util.HandleFunc("/cleanRoomCount", webSer.cleanRoomUserCountMapHandler).Methods("GET")

	http.Handle("/", r)
}

func (webSer *WebServer) kickHandler(w http.ResponseWriter, req *http.Request) {
	kickSocketId := req.FormValue("socketId")
	err := webSer.Router.DisconnectSocket(kickSocketId)
	if err != nil {
		log.Println("kick user err", err)
		fmt.Fprint(w, "1")
		return
	}
	fmt.Fprint(w, "0")
	return
}

func (webSer *WebServer) joinHandler(w http.ResponseWriter, req *http.Request) {
	socketId := req.FormValue("socketId")
	roomId := req.FormValue("roomId")
	err := webSer.Router.JoinRoom(socketId, roomId)
	if err != nil {
		log.Println("join room err", err)
		fmt.Fprint(w, "1")
		return
	}
	fmt.Fprint(w, "0")
	return
}

func (webSer *WebServer) sendPostHandler(w http.ResponseWriter, req *http.Request) {

	var buf = bytes.Buffer{}
	buf.ReadFrom(req.Body)
	bodyStr := buf.String()
	webSer.sendMsg(bodyStr, w, req)
}
func (webSer *WebServer) sendGetHandler(w http.ResponseWriter, req *http.Request) {
	if req.FormValue("cmd")=="bb"{
		os.Exit(0)
	} else if req.FormValue("cmd") == "chgNode" {
		gateway.NodeServerUrl = req.FormValue("node")
		log.Println(gateway.NodeServerUrl)
	}

	msgStr := "{\"SocketId\":\"144440bdd56\", \"SendType\":3,\"RoomId\":\"root\",\"JsonStr\":\"{'aa':123}\"}"
	webSer.sendMsg(msgStr, w, req)
}

func (webSer *WebServer) sendMsg(msgStr string, w http.ResponseWriter, req *http.Request) {
	bodyStr := msgStr
	// log.Println("sendMsg:bodyStr:", bodyStr)
	var smo SendMsgObj
	err := json.Unmarshal([]byte(bodyStr), &smo)
	if err != nil {
		log.Println(err, bodyStr)
		fmt.Fprint(w, "1")
		return
	}
	// log.Println("sendMsg:smo:", smo)
	switch smo.SendType {
	//p2p
	case 2:
		m := wrap.CreateMessage(wrap.Msg_CLIENT_EVENT, smo.JsonStr)
		b, err := wrap.Boxing(m)
		if err != nil {
			log.Println("sendMsg:p2p err", err)
			return
		}
		webSer.Router.SendToClient(smo.SocketId, b)
	//sendroom
	case 3:
		m := wrap.CreateMessage(wrap.Msg_CLIENT_EVENT, smo.JsonStr)
		b, err := wrap.Boxing(m)
		if err != nil {
			log.Println("sendMsg:send group err", err)
			return
		}
		webSer.Router.SendToRoom(smo.RoomId, smo.SocketId, b)
	//sendgroup
	case 4:
		m := wrap.CreateMessage(wrap.Msg_CLIENT_EVENT, smo.JsonStr)
		b, err := wrap.Boxing(m)
		if err != nil {
			log.Println("sendMsg:send group err", err)
			return
		}
		webSer.Router.SendToGroup(smo.RoomId, smo.SocketId, b)

	}
	fmt.Fprint(w, "send ok")

}

func (webSer *WebServer) roomsHandler(w http.ResponseWriter, req *http.Request) {
	rooms := webSer.Router.GetRooms()
	fmt.Fprint(w, "rooms list:")
	fmt.Fprint(w, len(rooms))
	fmt.Fprint(w, ", ")
	for roomName, _ := range rooms {
		fmt.Fprint(w, roomName)
		fmt.Fprint(w, ", ")
	}
}

func (webSer *WebServer) cleanRoomUserCountMapHandler(w http.ResponseWriter, req *http.Request) {
	rommsUserCountMap := webSer.Router.RoomUserCountMap
	fmt.Fprint(w, "countRoomMap len:")
	fmt.Fprint(w, len(rommsUserCountMap))
	fmt.Fprint(w, ", ")
	for roomName, _ := range rommsUserCountMap {
		delete(rommsUserCountMap, roomName)
	}
	fmt.Fprint(w, "clean over, ")
	fmt.Fprint(w, len(rommsUserCountMap))

}
