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
	LeavingWsChan chan TaibaiUserWsEvent
}

func NewTaibaiClassroomManager() *TaibaiClassroomManager {
	M := &TaibaiClassroomManager{
		ClassroomMap:  make(map[int]*TaibaiClassroom),
		PendingWsChan: make(chan TaibaiUserWsEvent, 3),
		LeavingWsChan: make(chan TaibaiUserWsEvent, 3),
	}

	go M.PendingNewWs()
	go M.LeavingOldWs()

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

func (this *TaibaiClassroomManager) LeavingOldWs() {
	for ws := range this.LeavingWsChan {
		log.Println("leavingOldWs")

		// 没找到教室 直接退出
		classroom, ok := this.ClassroomMap[ws.ClassroomId]
		if !ok {
			return
		}

		// 让用户离开教室
		classroom.onParticipantOffline(ws)
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
