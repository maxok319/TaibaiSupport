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
	taibaiWSConn := TaibaiClassroomManagerInstance.RegisterTaibaiWSConn(classroomId, userId, conn)
	go handleWSEvent(taibaiWSConn)
}

func handleWSEvent(taibaiWSConn *Models.TaibaiWSConn)  {
	// Todo 接MQ 就一股脑的给MQ发送就好了
	for event := range taibaiWSConn.EventChan{
		event.EventClassroomId = taibaiWSConn.ClassroomId
		HandleTaibaiClassroomEvent(event)
	}
}

func HandleTaibaiClassroomEvent(event *Models.TaibaiClassroomEvent)  {
	// Todo 接MQ 一股脑的从MQ拿到消息处理
	switch event.EventType {
	case Models.EventType_UserOnlineStatusChangd:
		HandleEventUserOnlineStatusChanged(event)
	case Models.EventType_UserVideoPositionChanged:
		HandleEventUserVideoPositionChanged(event)
	case Models.EventType_1V1StateChanged:
		HandleEvent1V1StateChanged(event)
	}
}