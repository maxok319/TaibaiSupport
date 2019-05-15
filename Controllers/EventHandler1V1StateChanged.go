package Controllers

import (
	"TaiBaiSupport/Models"
	"TaiBaiSupport/TaibaiDBHelper"
	"encoding/json"
	"github.com/bitly/go-simplejson"
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


	// 设置状态
	eventContent := simplejson.New()
	eventContent.SetPath([]string{}, event.EventContent)
	state1v1State, _ := eventContent.Get("1v1").Bool()

	userIdList := TaibaiDBHelper.GetUserList(event.EventClassroomId)
	perWidth := 1440 / len(userIdList)

	for _, userId := range userIdList {
		userStatus := TaibaiDBHelper.GetUserStatus(userId)
		userIndex := userStatus["index"].(int)
		rect := Models.TaibaiRect{}
		if state1v1State {
			rect.X = perWidth * userIndex + 10
			rect.Y = 40
			rect.Width = perWidth - 20
			rect.Height = 810 - 80 - 220
		} else {
			rect.X = 20 + (20+200) * userIndex
			rect.Y = 810 - 200
			rect.Width = 200
			rect.Height = 200
		}
		rectBytes,_ := json.Marshal(&rect)
		TaibaiDBHelper.UpdateUserRect(userId, string(rectBytes))
	}

	// 广播教室消息
	TaibaiClassroomManagerInstance.BroadcastClassroomStatus(event.EventClassroomId)

}