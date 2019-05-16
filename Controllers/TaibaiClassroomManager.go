package Controllers

import (
	"TaiBaiSupport/Models"
	"TaiBaiSupport/TaibaiDBHelper"
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
)

type TaibaiClassroomManager struct {
	OperationRWMux sync.RWMutex
	WSConns        map[int]map[int]*Models.TaibaiWSConn
}

func NewTaibaiClassroomManager() *TaibaiClassroomManager {
	M := &TaibaiClassroomManager{
		WSConns: make(map[int]map[int]*Models.TaibaiWSConn),
	}
	return M
}

func (this *TaibaiClassroomManager) RegisterTaibaiWSConn(classroomId, userId int, conn *websocket.Conn) *Models.TaibaiWSConn {

	this.OperationRWMux.Lock()
	defer this.OperationRWMux.Unlock()

	// redis存入一节课
	TaibaiDBHelper.AddClassroom(classroomId)

	// 先注册教室
	classroomContainer, classroomOk:= this.WSConns[classroomId]
	if !classroomOk{
		classroomContainer = make(map[int]*Models.TaibaiWSConn)
		this.WSConns[classroomId] = classroomContainer
	}

	// redis存入一个学生
	TaibaiDBHelper.AddUserIntoClassroom(classroomId, userId)

	// 先停掉之前的
	taibaiWSConn, wsConnOk := classroomContainer[userId]
	if wsConnOk{
		if taibaiWSConn.Conn != nil{
			_ = taibaiWSConn.Conn.Close()
		}
	}

	// 重新开启一个
	taibaiWSConn = Models.NewTaibaiWSConn(classroomId, userId, conn)
	this.WSConns[classroomId][userId] = taibaiWSConn
	return taibaiWSConn
}

func (this *TaibaiClassroomManager)BroadcastClassroomStatus(classroomId int)  {
	this.OperationRWMux.Lock()
	defer this.OperationRWMux.Unlock()

	classroomStatus :=  TaibaiDBHelper.GetClassroomStatus(classroomId)
	classroomStatus["classroomId"] = classroomId

	participantList := []interface{}{}
	for _, userId := range TaibaiDBHelper.GetUserList(classroomId) {
		participantStatus := TaibaiDBHelper.GetUserStatus(userId)
		participantStatus["userId"] = userId
		participantList = append(participantList, participantStatus)
	}
	classroomStatus["participantList"] = participantList

	message := Models.NewClassroomMessage(Models.MessageType_UpdateClassroomStatus, 0, []int{})
	message.MessageContent = classroomStatus
	clsStatus, _:= json.Marshal(message)
	userWSContainer := this.WSConns[classroomId]
	for _, conn := range userWSContainer {
		conn.SendMessage(clsStatus)
	}
}

var TaibaiClassroomManagerInstance *TaibaiClassroomManager

func init() {
	TaibaiClassroomManagerInstance = NewTaibaiClassroomManager()
}
