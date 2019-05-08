package Models

import (
	"fmt"
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
		if _, ok := this.GetClassroom(ws.ClassroomId); !ok {
			this.RegisterClassroom(NewTaibaiClassroom(ws.ClassroomId))
		}

		// 让用户上线
		classroom, _ := this.ClassroomMap[ws.ClassroomId]
		participant := classroom.addParticipant(ws.UserId)
		participant.SetConn(ws.Conn)

		// 通知教室里其他在线的人 有人上线了
		message := fmt.Sprintf("%d is online", ws.UserId)
		classroom.broadcastMessage(message)

	}
}

func (this *TaibaiClassroomManager) LeavingOldWs() {
	for ws := range this.LeavingWsChan {
		log.Println("leavingOldWs")

		// 没找到教室 直接退出
		if _, ok := this.GetClassroom(ws.ClassroomId); !ok {
			return
		}

		// 通知教室里其他在线的人 有人掉线了
		classroom, _ := this.ClassroomMap[ws.ClassroomId]
		message := fmt.Sprintf("%d is offline", ws.UserId)
		classroom.broadcastMessage(message)
	}
}

// 查询教室
func (this *TaibaiClassroomManager) GetClassroom(classroomId int) (classroom *TaibaiClassroom, ok bool) {
	this.OperationRWMux.RLock()
	defer this.OperationRWMux.RUnlock()

	classroom, ok = this.ClassroomMap[classroomId]
	return
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
