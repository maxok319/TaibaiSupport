package Controllers

import (
	"TaiBaiSupport/Models"
	"TaiBaiSupport/TaibaiDBHelper"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"log"
)

// {"rect":{"Height":200,"Width":200,"X":297,"Y":294},"userId":222}

func HandleEventUserVideoPositionChanged(event *Models.TaibaiClassroomEvent) {
	if event.EventType != Models.EventType_UserVideoPositionChanged {
		return
	}

	e, _ := json.Marshal(event.EventContent)
	log.Println("HandleEventUserVideoPositionChanged: ", string(e))


	// 设置状态
	eventContent := simplejson.New()
	eventContent.SetPath([]string{}, event.EventContent)
	userId, _ := eventContent.Get("userId").Int()
	rect, _ := eventContent.Get("rect").MarshalJSON()
	TaibaiDBHelper.UpdateUserRect(userId, string(rect))

	// 广播教室消息
	TaibaiClassroomManagerInstance.BroadcastClassroomStatus(event.EventClassroomId)
}