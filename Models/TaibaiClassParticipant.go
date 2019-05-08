package Models

import (
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
	User *TaibaiUser
	Classroom *TaibaiClassroom
	Role TaibaiClassRole
	Online bool
	Conn *websocket.Conn

	operateMutex sync.Mutex
}

func NewTaibaiClassParticipant(classroom *TaibaiClassroom, user *TaibaiUser, role TaibaiClassRole) *TaibaiClassParticipant  {
	return &TaibaiClassParticipant{
		Classroom:classroom,
		User: user,
		Role: role,
	}
}

func (this *TaibaiClassParticipant) SetConn(conn *websocket.Conn)  {
	this.operateMutex.Lock()
	defer this.operateMutex.Unlock()

	if this.Conn!=nil {
		err :=this.Conn.Close()
		if err != nil{
			println("close old conn error")
		}
	}


	if  conn == nil{
		this.Conn = nil
		this.Online = false
		return
	}

	this.Conn = conn
	this.Online = true

	go this.ReadLoop()
}

func (this *TaibaiClassParticipant) ReadLoop()  {
	for {
		_, message, err := this.Conn.ReadMessage()
		if err != nil {

			// 有异常的话 肯定要
			log.Println("read:", err)
			wsEvent := TaibaiUserWsEvent{
				ClassroomId: this.Classroom.ClassroomId,
				UserId:this.User.UserId,
				Conn:nil,
			}
			TaibaiClassroomManagerInstance.LeavingWsChan <- wsEvent
			break
		}
		log.Printf("recv: %s", message)
	}
}

func (this *TaibaiClassParticipant) SendMessage(message string)  {
	defer func() {recover()}()
	if this.Conn!=nil {
		this.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}