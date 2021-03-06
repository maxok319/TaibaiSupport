package Controllers

import (
	"TaibaiSupport/Models"
	"TaibaiSupport/TaibaiDBHelper"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"log"
)

/*
	{
		"Online": true,
		"ClassroomId": classroomId,
		"UserId":userId,
	}
*/

func HandleEventUserOnlineStatusChanged(event *Models.TaibaiClassroomEvent) {

	if event.EventType != Models.EventType_UserOnlineStatusChangd {
		return
	}
	e, _ := json.Marshal(event.EventContent)
	log.Println("HandleEventUserOnlineStatusChanged: ", string(e))


	// 设置状态
	eventContent := simplejson.New()
	eventContent.SetPath([]string{}, event.EventContent)
	userOnline, _ := eventContent.Get("online").Bool()
	classroomId, _:= eventContent.Get("classroomId").Int()
	userId, _ := eventContent.Get("userId").Int()

	// 可能在别的机器上登陆的
	TaibaiDBHelper.AddUserIntoClassroom(event.EventClassroomId, userId)

	// 修改状态
	TaibaiDBHelper.UpdateUserOnlineStatus(userId, userOnline)

	// 广播教室消息
	TaibaiClassroomManagerInstance.BroadcastClassroomStatus(classroomId)
}
