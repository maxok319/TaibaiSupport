package Models

import "C"
import (
	"TaiBaiSupport/TaibaiJson"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type TaibaiWSConn struct {
	Conn        *websocket.Conn
	ClassroomId int
	UserId      int
	EventChan  chan *TaibaiClassroomEvent
}

func NewTaibaiWSConn(classroomId, userId int, conn *websocket.Conn) *TaibaiWSConn{
	taibaiWSConn := &TaibaiWSConn{
		ClassroomId:classroomId,
		UserId:userId,
		Conn:conn,
		EventChan:make(chan *TaibaiClassroomEvent, 10),
	}

	// 模拟一个用户上线事件
	event := NewTaibaiClassroomEvent(EventType_UserOnlineStatusChangd)
	event.EventContent = map[string]interface{}{
		"UserOnline":  true,
		"ClassroomId": classroomId,
		"UserId":      userId,
	}
	taibaiWSConn.EventChan <- event

	// 开启协程去读消息
	go taibaiWSConn.startReadLoop()
	return taibaiWSConn
}

func (this *TaibaiWSConn) startReadLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("捕获到的错误：%v\n", r)
		}
	}()
	for {
		_, message, err := this.Conn.ReadMessage()
		if err != nil {
			log.Println("Read WS:", err)

			// 模拟一个用户掉线事件
			event := NewTaibaiClassroomEvent(EventType_UserOnlineStatusChangd)
			event.EventContent = TaibaiJson.JsonObject{
				"UserOnline":  false,
				"ClassroomId": this.ClassroomId,
				"UserId":      this.UserId,
			}
			this.EventChan <- event
			close(this.EventChan)
			this.Conn = nil
			return

		} else {
			event := TaibaiClassroomEvent{}
			_ = json.Unmarshal(message, &event)
			this.EventChan <- &event
		}
	}
}
