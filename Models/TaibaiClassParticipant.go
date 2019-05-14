package Models

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type TaibaiClassParticipant struct {
	UserId       int
	Classroom    *TaibaiClassroom
	Conn         *websocket.Conn
	operateMutex sync.Mutex

	Online bool
	Index  int
	Rect   TaibaiRect
}

func NewTaibaiClassParticipant(classroom *TaibaiClassroom, userId int) *TaibaiClassParticipant {
	p := &TaibaiClassParticipant{
		Classroom: classroom,
		UserId:    userId,
	}
	return p
}

func (this *TaibaiClassParticipant) SetConn(conn *websocket.Conn) {
	this.operateMutex.Lock()
	defer this.operateMutex.Unlock()

	// 保存老的websocket
	oldConn := this.Conn

	this.Conn = conn
	this.Online = true
	go this.ReadLoop(this.Conn)

	// 在最后断掉老的
	if oldConn != nil {
		err := oldConn.Close()
		if err != nil {
			println("close old conn error")
		}
	}
}

func (this *TaibaiClassParticipant) ReadLoop(Conn *websocket.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("捕获到的错误：%v\n", r)
		}
	}()
	for {
		_, message, err := Conn.ReadMessage()
		if err != nil {
			log.Println("Read WS:", err)
			if Conn == this.Conn {
				this.Conn = nil
				this.Online = false

				wsEvent := TaibaiUserWsEvent{
					ClassroomId: this.Classroom.ClassroomId,
					UserId:      this.UserId,
					Conn:        nil,
				}
				this.Classroom.onParticipantOffline(wsEvent)
			}
			return

		} else {
			this.ReceiveMessage(message)
		}
	}

}

func (this *TaibaiClassParticipant) ReceiveMessage(message []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Receive Message：%v\n", r)
		}
	}()

	log.Println("receive", string(message))

	this.Classroom.onParticipantReceiveWSMessage(this, message)
}

func (this *TaibaiClassParticipant) SendMessage(message string) {
	defer func() { recover() }()
	if this.Conn != nil {
		this.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}
