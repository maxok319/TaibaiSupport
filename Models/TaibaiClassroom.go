package Models

import (
	"TaiBaiSupport/TaibaiDBHelper"
	"TaiBaiSupport/TaibaiJson"
	"encoding/json"
	"log"
	"strconv"
	"time"
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
		p.Rect.X = 20 + (20+200)*p.Index
		p.Rect.Y = 810 - 200
		p.Rect.Width = 200
		p.Rect.Height = 200
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

func (this *TaibaiClassroom) participantPositionChanged(userId int, rect TaibaiRect){
	if participant, ok:= this.Participants[userId]; ok{
		participant.Rect = rect
	}
}

func (this *TaibaiClassroom) broadcastClassroomStatus() {
	classroomStatus := TaibaiJson.JsonObject{}
	classroomStatus["classroomId"] = this.ClassroomId
	participantList := TaibaiJson.JsonArray{}
	for _, p := range this.Participants {
		participantStatus := TaibaiJson.JsonObject{}
		participantStatus["index"] = p.Index
		participantStatus["online"] = p.Online
		participantStatus["userId"] = p.User.UserId
		participantStatus["rect"] = p.Rect
		participantList = append(participantList, participantStatus)
	}
	classroomStatus["participantList"] = participantList

	message := TaibaiJson.JsonObject{}
	message["messageType"] = "classroomStatus"
	message["messageTime"] = time.Now().Unix()
	message["messageContent"] = classroomStatus

	wspackage, _ := json.Marshal(message)
	this.broadcastMessage(string(wspackage))
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

func (this *TaibaiClassroom) sendClassroomMessage(message *TaibaiClassroomMessage)  {
	messageBytes,_ := json.Marshal(message)
	this.saveActionIntoRedis(messageBytes)

	if len(message.MessageReceiver) == 0{
		this.broadcastMessage(string(messageBytes))
	} else {
		for userId := range message.MessageReceiver{
			this.singleMessage(userId, string(messageBytes))
		}
	}
}

func (this *TaibaiClassroom) saveActionIntoRedis (action interface{}) {
	listKey := "actionlist:" + strconv.Itoa(this.ClassroomId)
	TaibaiDBHelper.GetRedisClient().RPush(listKey, action)
}