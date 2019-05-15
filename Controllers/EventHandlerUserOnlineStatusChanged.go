package Controllers

import (
	"TaiBaiSupport/Models"
	"encoding/json"
	"log"
)

/*
	{
		"UserOnline": true,
		"ClassroomId": classroomId,
		"UserId":userId,
	}
*/

func HandleEventUserOnlineStatusChanged(event *Models.TaibaiClassroomEvent) {

	if event.EventType != Models.EventType_UserOnlineStatusChangd {
		return
	}

	e, _ := json.Marshal(event.EventContent)
	log.Println(string(e))
}


