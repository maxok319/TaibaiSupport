package Controllers

import (
	"TaibaiSupport/Models"
	"encoding/json"
	"log"
)

func HandleEventRTCICECandidateChanged(event *Models.TaibaiClassroomEvent) {
	if event.EventType != Models.EventType_RTCICECandidateChanged {
		return
	}

	e, _ := json.Marshal(event.EventContent)
	log.Println("HandleEventRTCICECandidateChanged: ", string(e))

	// 直接广播
	TaibaiClassroomManagerInstance.BroadcastOriginEvent(event.EventClassroomId, *event)
}