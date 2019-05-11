package Models

import (
	"TaiBaiSupport/TaibaiJson"
	"TaiBaiSupport/TaibaiUtils"
	"context"
	"encoding/json"
	"github.com/bitly/go-simplejson"
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
	Rect      TaibaiRect

	Conn     *websocket.Conn
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

	// 保存老的websocket
	oldConn := this.Conn

	this.Conn = conn
	this.Online = true
	go this.ReadLoop(this.Conn)

	// 在最后断掉老的
	if oldConn != nil {
		err := oldConn.Close()
		if err != nil {
			println("close old conn error")
		}
	}
}

func (this *TaibaiClassParticipant) ReadLoop(Conn *websocket.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("捕获到的错误：%v\n", r)
		}
	}()
	for {
		_, message, err := Conn.ReadMessage()
		if err != nil {
			log.Println("Read WS:", err)
			if Conn == this.Conn {
				this.Conn = nil
				this.Online = false

				wsEvent := TaibaiUserWsEvent{
					ClassroomId: this.Classroom.ClassroomId,
					UserId:      this.User.UserId,
					Conn:        nil,
				}
				TaibaiClassroomManagerInstance.LeavingWsChan <- wsEvent
			}
			return

		} else {
			this.ReceiveMessage(message)
		}
	}

}

func (this *TaibaiClassParticipant) ReceiveMessage(message []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Receive Message：%v\n", r)
		}
	}()

	log.Println("receive", string(message))

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

func (this *TaibaiClassParticipant) SendMessage(message string) {
	defer func() { recover() }()
	if this.Conn != nil {
		this.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}

func (this *TaibaiClassParticipant) onUserVideoPositionChanged(event *TaibaiClassroomEvent) {
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
	this.Classroom.participantPositionChanged(userId, rect)

	message := NewClassroomMessage(MessageType_UpdateUserVideoPosition, this.User.UserId, []int{})
	// event里只设置了一个人的位置 但可能造成了多人的位置改动 1V1模式等
	messageContent := TaibaiJson.JsonArray{}
	messageContent = append(messageContent, eventContent)
	message.MessageContent = messageContent
	message.MessageOriginEvent = *event
	this.Classroom.sendClassroomMessage(message)
}

func (this *TaibaiClassParticipant) on1v1StateChanged(event *TaibaiClassroomEvent) {
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
	message := NewClassroomMessage(MessageType_UpdateUserVideoPosition, this.User.UserId, []int{})
	// event里只设置了一个人的位置 但可能造成了多人的位置改动 1V1模式等
	messageContent := TaibaiJson.JsonArray{}

	perWidth := 1440 / len(this.Classroom.Participants)
	for _, particpant := range this.Classroom.Participants {
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
		this.Classroom.participantPositionChanged(userId, rect)

		particpantPosition := TaibaiJson.JsonObject{
			"userId": userId,
			"rect": rect,
		}
		messageContent = append(messageContent, particpantPosition)
	}


	message.MessageContent = messageContent
	message.MessageOriginEvent = *event
	this.Classroom.sendClassroomMessage(message)
}
