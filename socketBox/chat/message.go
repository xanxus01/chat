package chat

import (
	"awesomeProject/socketBox"
	"encoding/json"
)

type message struct {
	from string
	to string
	msgType string
	content string
}

func (m *message)sendTo(c *socketBox.Connection)  {
	data,_ :=json.Marshal(m)
	c.WriteMsg(data)
}