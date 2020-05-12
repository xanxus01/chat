package msg

import "time"

type untiyMsg struct {
	msgType    string
	msgText    string
	sendTime   time.Time
	accpetTime time.Time
}

func NewUntiyMsg(msgType string, msgText string, sendTime time.Time) *untiyMsg {
	return &untiyMsg{
		msgType:    msgType,
		msgText:    msgText,
		sendTime:   sendTime,
		accpetTime: time.Now(),
	}
}
