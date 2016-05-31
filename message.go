package main

type Message struct {
	To   string      `json:"to"`
	From string      `json:"from"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (msg *Message) Clear() {
	msg.To = ""
	msg.From = ""
	msg.Type = ""
	msg.Data = nil
}
