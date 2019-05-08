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
		this.ParticipantOnline(ws.ClassroomId, ws.UserId, ws.Conn)
	}
}

func (this *TaibaiClassroomManager) LeavingOldWs() {
	for ws := range this.LeavingWsChan {
		log.Println("leavingOldWs")

		// 没找到教室 直接退出
		if _, ok := this.GetClassroom(ws.ClassroomId); !ok {
			return
		}

		// 让用户下线
		this.ParticipantOffline(ws.ClassroomId, ws.UserId)
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

// 用户上线
func (this *TaibaiClassroomManager) ParticipantOnline(classroomId, userId int, conn *websocket.Conn) {
	this.OperationRWMux.Lock()
	defer this.OperationRWMux.Unlock()

	classroom, _ := this.ClassroomMap[classroomId]
	participant := classroom.addParticipant(userId)
	participant.SetConn(conn)

	// 通知教室里其他在线的人 有人上线了
	message := fmt.Sprintf("%d is online", userId)
	classroom.broadcastMessage(message)
}

// 用户下线
func (this *TaibaiClassroomManager) ParticipantOffline(classroomId, userId int) {
	this.OperationRWMux.Lock()
	defer this.OperationRWMux.Unlock()

	classroom, _ := this.ClassroomMap[classroomId]

	// 通知教室里其他在线的人 有人掉线了
	message := fmt.Sprintf("%d is offline", userId)
	classroom.broadcastMessage(message)
}

var TaibaiClassroomManagerInstance *TaibaiClassroomManager

func init() {
	TaibaiClassroomManagerInstance = NewTaibaiClassroomManager()
}
