package Models

type TaibaiEventType int

const (
	EventType_UserOnlineStatusChangd TaibaiEventType = iota
	EventType_UserVideoPositionChanged
	EventType_1V1StateChanged
)

type TaibaiClassroomEvent struct {
	EventID       int             `json:"eventId"`       // 事件ID
	EventTime     int64           `json:"eventTime"`     // 事件时间(ns)
	EventType     TaibaiEventType `json:"eventType"`     // 事件类型
	EventSender   int             `json:"eventSender"`   // 事件主播
	EventContent  interface{}     `json:"eventContent"`  // 事件内容
}
