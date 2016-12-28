package gateway

import (
	"golang.org/x/net/websocket"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ReportSecond    int = 10
	PingSecond      int = 5
	InitConnTimeout int = 5
)

var (
	NodeServerUrl       = "http://221.122.71.51:8103"
	GoServerID          = "51"
	IdTail        int64 = 1
)

type ReportResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Sum  int    `json:"sum"`
}

func CreateRouter(handler Handler) *Router {
	r := &Router{
		handler:          handler,
		chUp:             make(chan *Packet),
		chDown:           make(chan *Packet),
		rooms:            make(map[string]map[string]*Socket),
		chConState:       make(chan *conState),
		RoomUserCountMap: make(map[string]int),
	}
	r.init()
	go r.runUp()
	go r.runDown()
	return r
}

func ReportUserCount(r *Router) {
	reportTimer := time.NewTicker(time.Duration(ReportSecond) * time.Second) //报告人数定时器
	for {
		<-reportTimer.C
		rooms := r.GetRooms()
		for roomName, room := range rooms {
			if roomName != ROOT_ROOM {
				roomUserCount := len(room)
				if roomUserCount > 0 {

					go func(roomName string, roomUserCount int) {

						url := NodeServerUrl + "/live/refresh-online-count?liveClassroomId=" +
							roomName + "&count=" + strconv.Itoa(roomUserCount) + "&goId=" + GoServerID
						resp, err := http.Get(url)

						defer resp.Body.Close()

						bodyBytes, err2 := ioutil.ReadAll(resp.Body)
						if err2 != nil {
							return
						}
						var repRes ReportResult
						err3 := json.Unmarshal(bodyBytes, &repRes)
						if err3 != nil {
							return
						}
						if repRes.Code == 0 {
							r.RoomUserCountMap[roomName] = repRes.Sum
						}

					}(roomName, roomUserCount)

				}
			}
		}
	}
}

func CreateWsHandler(router *Router) http.Handler {
	return websocket.Handler(func(wsConn *websocket.Conn) {
		adapter := &wsAdapter{conn: wsConn}
		socket := <-initConnection(adapter, router)
		socket.read(router.chConState, router.chUp)
	})
}
func CreateTcpListener(addr string, router *Router) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("tcp listener err", err)
		return
	}

	initWithTimeout := func(conn net.Conn) {
		adapter := &tcpAdapter{conn: conn}
		select {

		case <-time.After(time.Duration(InitConnTimeout) * time.Second):
			adapter.Close()

		case socket := <-initConnection(adapter, router):
			socket.read(router.chConState, router.chUp)
		}
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error accepting", err)
				continue
			}
			go initWithTimeout(conn)
		}
	}()
}

func initConnection(adapter ConAdapter, router *Router) <-chan *Socket {

	chSocket := make(chan *Socket)

	socketId := generateId()
	remoteIp := adapter.RemoteIp()
	socket := &Socket{
		SocketId:    socketId,
		adapter:     adapter,
		ticker:      time.NewTicker(time.Duration(PingSecond) * time.Second), //定时器 用ticker.C这个channel到时间就发心跳
		ch:          make(chan []byte, 256),
		ClientRooms: make(map[string]bool),
		router:      router,
		SocketState: true,
		RemoteIp:    remoteIp,
	}
	router.chConState <- &conState{socket, true}
	go socket.write(router.chConState)

	go socket.sendHandshake(chSocket)
	return chSocket
}

func generateId() string {
	IdTail++
	IdTail %= 100
	now := time.Now().UnixNano()
	nowi := now + IdTail
	return GoServerID + strconv.FormatInt(nowi, 16)
}
