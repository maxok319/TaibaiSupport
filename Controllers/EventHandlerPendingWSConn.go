package Controllers

import (
	"TaibaiSupport/Models"
	"encoding/json"
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

// 收到ws消息 都先推给mq
func handleWSEvent(taibaiWSConn *Models.TaibaiWSConn)  {
	for event := range taibaiWSConn.EventChan{
		event.EventClassroomId = taibaiWSConn.ClassroomId
		eventJson,_ := json.Marshal(event)
		log.Println("给mq发送：", string(eventJson))
		RabbitmqEventTobeSendChan <- *event
	}
}

