package Controllers

import (
	"TaiBaiSupport/Models"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{}


// 处理新的链接
func HandleEventPendingWS(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	v := r.URL.Query()
	classroomId, _ := strconv.Atoi(v.Get("classroomId"))
	userId, _ := strconv.Atoi(v.Get("userId"))

	// 开启协程去取数据
	taibaiWSConn := Models.NewTaibaiWSConn(classroomId, userId, conn)
	go handleWSEvent(taibaiWSConn)
}

func handleWSEvent(taibaiWSConn *Models.TaibaiWSConn)  {
	for event := range taibaiWSConn.EventChan{
		HandleTaibaiClassroomEvent(event)
	}
}

func HandleTaibaiClassroomEvent(event *Models.TaibaiClassroomEvent)  {
	switch event.EventType {
	case Models.EventType_UserOnlineStatusChangd:
		HandleEventUserOnlineStatusChanged(event)
	case Models.EventType_UserVideoPositionChanged:
		HandleEventUserVideoPositionChanged(event)
	case Models.EventType_1V1StateChanged:
		HandleEvent1V1StateChanged(event)
	}
}