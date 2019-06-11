package Controllers

import "TaibaiSupport/Models"

func HandleEventRTCSDPChanged(event *Models.TaibaiClassroomEvent) {
	if event.EventType != Models.EventType_RTCSDPChanged {
		return
	}

}