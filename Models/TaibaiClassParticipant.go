package Models

import (
	"TaiBaiSupport/TaibaiDBHelper"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
)

type TaibaiClassParticipant struct {
	UserId       int
	Classroom    *TaibaiClassroom
	Conn         *websocket.Conn
	operateMutex sync.Mutex

	Index int
	Rect  TaibaiRect
}

func NewTaibaiClassParticipant(classroom *TaibaiClassroom, userId int) *TaibaiClassParticipant {
	p := &TaibaiClassParticipant{
		Classroom: classroom,
		UserId:    userId,
	}
	return p
}

func (this *TaibaiClassParticipant) GetIndex() int {
	// 去库里查这个index
	index, err := TaibaiDBHelper.GetInstance().HGet(this.GetRedisClassroomAndUserIdKey(), "index").Int()
	if err != nil {
		index = len(this.Classroom.Participants)
		this.SetIndex(index)
	}
	return index
}

func (this *TaibaiClassParticipant) SetIndex(index int) {
	TaibaiDBHelper.GetInstance().HSet(this.GetRedisClassroomAndUserIdKey(), "index", index)
}

func (this *TaibaiClassParticipant) GetRect() TaibaiRect {
	rect := TaibaiRect{}
	rectStr, err := TaibaiDBHelper.GetInstance().HGet(this.GetRedisClassroomAndUserIdKey(), "rect").Result()
	if err != nil {
		rect.X = 20 + (20+200)*this.GetIndex()
		rect.Y = 810 - 200
		rect.Width = 200
		rect.Height = 200
	} else {
		_ = json.Unmarshal([]byte(rectStr), &rect)
	}
	return rect
}

func (this *TaibaiClassParticipant) SetRect(rect TaibaiRect) {
	rectStr, _ := json.Marshal(rect)
	_ = TaibaiDBHelper.GetInstance().HSet(this.GetRedisClassroomAndUserIdKey(), "rect", string(rectStr))
}

func (this *TaibaiClassParticipant) GetRedisClassroomAndUserIdKey() string {
	return strconv.Itoa(this.Classroom.ClassroomId) + strconv.Itoa(this.UserId)
}

func (this *TaibaiClassParticipant) SetConn(conn *websocket.Conn) {
	this.operateMutex.Lock()
	defer this.operateMutex.Unlock()

	// 保存老的websocket
	oldConn := this.Conn

	this.Conn = conn
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

				wsEvent := TaibaiWSConn{
					ClassroomId: this.Classroom.ClassroomId,
					UserId:      this.UserId,
					Conn:        nil,
				}
				this.Classroom.onParticipantOffline(wsEvent)
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

	this.Classroom.onParticipantReceiveWSMessage(this, message)
}

func (this *TaibaiClassParticipant) SendMessage(message string) {
	defer func() { recover() }()
	if this.Conn != nil {
		this.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}
