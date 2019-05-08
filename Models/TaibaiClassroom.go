package Models

import (
	"TaiBaiSupport/TaibaiJson"
	"encoding/json"
	"log"
)

type TaibaiClassroom struct {
	ClassroomId int
	StartTime   int
	StopTime    int

	Participants map[int]*TaibaiClassParticipant
}

// 这里要读取库才可以 传入一个json串
func NewTaibaiClassroom(classroomId int) *TaibaiClassroom {

	classroom := &TaibaiClassroom{
		ClassroomId:  classroomId,
		Participants: make(map[int]*TaibaiClassParticipant),
	}

	return classroom
}

func (this *TaibaiClassroom) addParticipant(userId int) *TaibaiClassParticipant {
	p, ok := this.Participants[userId]
	if !ok {
		p = NewTaibaiClassParticipant(this, &TaibaiUser{UserId: userId}, StudentRole)
		p.Index = len(this.Participants)
		this.Participants[userId] = p
	}
	return p
}

func (this *TaibaiClassroom) participantOnline(ws TaibaiUserWsEvent) {
	participant := this.addParticipant(ws.UserId)
	participant.SetConn(ws.Conn)

	log.Printf("%d is online", ws.UserId)
	this.broadcastClassroomStatus()
}

func (this *TaibaiClassroom) participantOffline(ws TaibaiUserWsEvent) {
	log.Printf("%d is offline", ws.UserId)
	// 通知教室里其他在线的人 有人上线了
	this.broadcastClassroomStatus()
}

func (this * TaibaiClassroom)broadcastClassroomStatus()  {
	classroomStatus := TaibaiJson.JsonObject{}
	classroomStatus["classroomId"] = this.ClassroomId
	participantList := TaibaiJson.JsonArray{}
	for _, p := range this.Participants {
		participantStatus := TaibaiJson.JsonObject{}
		participantStatus["index"] = p.Index
		participantStatus["online"] = p.Online
		participantStatus["userId"] = p.User.UserId
		participantList = append(participantList, participantStatus)
	}
	classroomStatus["participantList"] = participantList

	message, _ := json.Marshal(classroomStatus)
	this.broadcastMessage(string(message))
}

func (this *TaibaiClassroom) broadcastMessage(message string) {
	for _, p := range this.Participants {
		p.SendMessage(message)
	}
}

func (this *TaibaiClassroom) singleMessage(userId int, message string) {
	if p, ok := this.Participants[userId]; ok {
		p.SendMessage(message)
	}
}
