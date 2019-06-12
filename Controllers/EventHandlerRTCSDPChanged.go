package Controllers

import (
	"TaibaiSupport/Models"
	"encoding/json"
	"log"
)

func HandleEventRTCSDPChanged(event *Models.TaibaiClassroomEvent) {
	if event.EventType != Models.EventType_RTCSDPChanged {
		return
	}

	e, _ := json.Marshal(event.EventContent)
	log.Println("HandleEventRTCSDPChanged: ", string(e))

	// 直接广播
	TaibaiClassroomManagerInstance.BroadcastOriginEvent(event.EventClassroomId, *event)
}