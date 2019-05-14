package Models

import (
	"TaiBaiSupport/TaibaiDBHelper"
	"TaiBaiSupport/TaibaiJson"
	"TaiBaiSupport/TaibaiUtils"
	"encoding/json"
	"github.com/bitly/go-simplejson"
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


func (this *TaibaiClassroom) sendEventMQ(event TaibaiClassroomEvent){

}


func (this *TaibaiClassroom) saveActionIntoRedis (action interface{}) {
	listKey := "actionlist:" + strconv.Itoa(this.ClassroomId)
	TaibaiDBHelper.GetRedisClient().RPush(listKey, action)
}

func (this *TaibaiClassroom) onParticipantReceivedEvent(participant *TaibaiClassParticipant, message []byte ){
	event := &TaibaiClassroomEvent{}
	err := json.Unmarshal(message, event)
	if err != nil {
		log.Println(err)
		return
	}

	eventType := event.EventType
	switch eventType {
	case EventType_UserVideoPositionChanged:
		this.onUserVideoPositionChanged(event)
	case EventType_1V1StateChanged:
		this.on1v1StateChanged(event)
	}
}


// 0. client给server一个action
// 1. server转发action给别的server
// 2. server合成message给其clients
// 3. server将message保存至redis

func (this *TaibaiClassroom) onParticipantOnline(ws TaibaiUserWsEvent) {
	participant := this.addParticipant(ws.UserId)
	participant.SetConn(ws.Conn)

	log.Printf("%d is online", ws.UserId)
	this.broadcastClassroomStatus()
}

func (this *TaibaiClassroom) onParticipantOffline(ws TaibaiUserWsEvent) {
	log.Printf("%d is offline", ws.UserId)
	// 通知教室里其他在线的人 有人上线了
	this.broadcastClassroomStatus()
}

func (this *TaibaiClassroom) onUserVideoPositionChanged(event *TaibaiClassroomEvent) {
	/*
		{
		    "eventTime": 1557489041,
		    "eventType": 1,
		    "eventProducer": 0,
		    "eventContent": {
		        "userId": 111,
		        "rect": {
		            "X": 189.0,
		            "Y": 506.99999999999994,
		            "Width": 200.0,
		            "Height": 200.0
		        }
		    }
		}
	*/
	eventContent := event.EventContent
	eventContentObject := simplejson.New()
	eventContentObject.SetPath([]string{}, eventContent)

	userId := eventContentObject.Get("userId").MustInt()
	rect := TaibaiRect{}
	TaibaiUtils.SimpleJsonToStruct(eventContentObject.Get("rect"), &rect)
	if participant, ok:= this.Participants[userId]; ok{
		participant.Rect = rect
	}

	message := NewClassroomMessage(MessageType_UpdateUserVideoPosition, 0, []int{})
	// event里只设置了一个人的位置 但可能造成了多人的位置改动 1V1模式等
	messageContent := TaibaiJson.JsonArray{}
	messageContent = append(messageContent, eventContent)
	message.MessageContent = messageContent
	message.MessageOriginEvent = *event
	this.sendClassroomMessage(message)
}

func (this *TaibaiClassroom) on1v1StateChanged(event *TaibaiClassroomEvent) {
	/*
		{
		    "eventTime": 1557598302,
		    "eventType": 2,
		    "eventProducer": 0,
		    "eventContent": {
		        "1v1": true
		    }
		}
	*/
	eventContent := event.EventContent
	eventContentObject := simplejson.New()
	eventContentObject.SetPath([]string{}, eventContent)

	state1v1State := eventContentObject.Get("1v1").MustBool()
	message := NewClassroomMessage(MessageType_UpdateUserVideoPosition, 0, []int{})
	// event里只设置了一个人的位置 但可能造成了多人的位置改动 1V1模式等
	messageContent := TaibaiJson.JsonArray{}

	perWidth := 1440 / len(this.Participants)
	for _, particpant := range this.Participants {
		rect := TaibaiRect{}

		if state1v1State {
			rect.X = perWidth * particpant.Index + 10
			rect.Y = 40
			rect.Width = perWidth - 20
			rect.Height = 810 - 80
		} else {
			rect.X = 20 + (20+200)*particpant.Index
			rect.Y = 810 - 200
			rect.Width = 200
			rect.Height = 200
		}

		userId := particpant.User.UserId
		if participant, ok:= this.Participants[userId]; ok{
			participant.Rect = rect
		}

		particpantPosition := TaibaiJson.JsonObject{
			"userId": userId,
			"rect": rect,
		}
		messageContent = append(messageContent, particpantPosition)
	}


	message.MessageContent = messageContent
	message.MessageOriginEvent = *event
	this.sendClassroomMessage(message)
}
