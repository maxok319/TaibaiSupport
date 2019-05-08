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

	Conn *websocket.Conn

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

	// 先断掉原来老的websocket
	oldConn := this.Conn


	// 这次是要来了空的
	if conn == nil {
		log.Println("selt conn nil")
		this.Conn = nil
		this.Online = false
	} else {
		this.Conn = conn
		this.Online = true
		go this.ReadLoop(this.Conn)
	}

	if oldConn != nil {
		err := oldConn.Close()
		if err != nil {
			println("close old conn error")
		}
	}
}

func (this *TaibaiClassParticipant) ReadLoop(Conn* websocket.Conn) {
	defer func() { recover() }()
	for {
		_, message, err := Conn.ReadMessage()
		if err != nil {
			// 有异常的话 肯定要
			log.Println("read:", err)
			wsEvent := TaibaiUserWsEvent{
				ClassroomId: this.Classroom.ClassroomId,
				UserId:      this.User.UserId,
				Conn:        nil,
			}
			if Conn==this.Conn {
				TaibaiClassroomManagerInstance.LeavingWsChan <- wsEvent
			}
			return
		} else {
			log.Printf("recv: %s", message)
		}
	}

}

func (this *TaibaiClassParticipant) SendMessage(message string) {
	defer func() { recover() }()
	if this.Conn != nil {
		this.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}
