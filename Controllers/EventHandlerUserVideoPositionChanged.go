package Controllers

import (
	"TaiBaiSupport/Models"
	"encoding/json"
	"log"
)

func HandleEventUserVideoPositionChanged(event *Models.TaibaiClassroomEvent) {
	if event.EventType != Models.EventType_UserVideoPositionChanged {
		return
	}

	e, _ := json.Marshal(event.EventContent)
	log.Println(string(e))
}