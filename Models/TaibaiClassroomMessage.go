package Models

import (
	"TaiBaiSupport/TaibaiUtils"
	"time"
)

type TaibaiMessageType int

const (
	MessageType_UpdateClassroomStatus TaibaiMessageType = iota
)

// 0. client给server一个action
// 1. server合成message给其clients
// 2. server将message保存至redis

type TaibaiClassroomMessage struct {
	MessageID          int                  `json:"messageId"`          // 消息ID
	MessageTime        int64                `json:"messageTime"`        // 消息时间(ns)
	MessageType        TaibaiMessageType    `json:"messageType"`        // 消息类型
	MessageSender      int                  `json:"messageSender"`      // 消息主播(0代表系统消息)
	MessageReceiver    []int                `json:"messageReceiver"`    // 消息听众(空代表给所有人)
	MessageContent     interface{}          `json:"messageContent"`     // 消息内容
	MessageOriginEvent TaibaiClassroomEvent `json:"messageOriginEvent"` // 消息的原始事件
}

func NewClassroomMessage(messageType TaibaiMessageType, sender int, receiver []int) *TaibaiClassroomMessage {
	m := TaibaiClassroomMessage{}
	m.MessageID = TaibaiUtils.GenerateMessageId()
	m.MessageTime = time.Now().Unix()
	m.MessageType = messageType
	m.MessageSender = sender
	m.MessageReceiver = receiver
	return &m
}
