package Models

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type TaibaiUserWsEvent struct {
	ClassroomId int
	UserId      int
	Conn        *websocket.Conn
}

type TaibaiClassroomManager struct {
	ClassroomMap   map[int]*TaibaiClassroom
	OperationRWMux sync.RWMutex

	PendingWsChan chan TaibaiUserWsEvent
}

func NewTaibaiClassroomManager() *TaibaiClassroomManager {
	M := &TaibaiClassroomManager{
		ClassroomMap:  make(map[int]*TaibaiClassroom),
		PendingWsChan: make(chan TaibaiUserWsEvent, 3),
	}

	go M.PendingNewWs()

	return M
}

func (this *TaibaiClassroomManager) PendingNewWs() {
	for ws := range this.PendingWsChan {
		log.Println("pendingNewWs")

		// 未注册教室 先注册教室
		classroom, ok := this.ClassroomMap[ws.ClassroomId]
		if !ok {
			classroom = NewTaibaiClassroom(ws.ClassroomId)
			this.RegisterClassroom(classroom)
		}

		// 让用户加入教室
		classroom.onParticipantOnline(ws)
	}
}

// 注册教室
func (this *TaibaiClassroomManager) RegisterClassroom(classroom *TaibaiClassroom) {
	this.OperationRWMux.Lock()
	defer this.OperationRWMux.Unlock()

	this.ClassroomMap[classroom.ClassroomId] = classroom
}

var TaibaiClassroomManagerInstance *TaibaiClassroomManager

func init() {
	TaibaiClassroomManagerInstance = NewTaibaiClassroomManager()
}
