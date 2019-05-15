package Controllers

import (
	"TaiBaiSupport/Models"
	"encoding/json"
	"log"
)

/*
	{
		"1v1":true
	}
*/

func HandleEvent1V1StateChanged(event *Models.TaibaiClassroomEvent) {
	if event.EventType != Models.EventType_1V1StateChanged {
		return
	}

	e, _ := json.Marshal(event.EventContent)
	log.Println(string(e))
}