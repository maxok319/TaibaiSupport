package Models

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type TaibaiClassRole int

const (
	TeacherRole TaibaiClassRole = iota
	StudentRole
	ObserverRole
)

type TaibaiClassParticipant struct {
	User      *TaibaiUser
	Classroom *TaibaiClassroom
	Role      TaibaiClassRole
	Online    bool
	Index     int
	Rect      TaibaiRect

	Conn     *websocket.Conn
	ConnCtx  context.Context
	ConnStop context.CancelFunc

	operateMutex sync.Mutex
}

func NewTaibaiClassParticipant(classroom *TaibaiClassroom, user *TaibaiUser, role TaibaiClassRole) *TaibaiClassParticipant {
	p := &TaibaiClassParticipant{
		Classroom: classroom,
		User:      user,
		Role:      role,
	}
	p.ConnCtx, p.ConnStop = context.WithCancel(context.Background())
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
					UserId:      this.User.UserId,
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

	this.Classroom.onParticipantReceivedEvent(this, message)
}

func (this *TaibaiClassParticipant) SendMessage(message string) {
	defer func() { recover() }()
	if this.Conn != nil {
		this.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}