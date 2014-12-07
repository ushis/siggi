package sighub

type Message struct {
	To   string      `json:"to"`
	From string      `json:"from"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var messagePool = make(chan *Message, 100)

func NewMessage() (msg *Message) {
	select {
	case msg = <-messagePool:
		// Got one from the pool
	default:
		msg = new(Message)
	}
	msg.To = ""
	msg.From = ""
	msg.Type = ""
	msg.Data = nil
	return msg
}

func (m *Message) Free() {
	select {
	case messagePool <- m:
		// Stored message in the pool
	default:
		// Pool is full. It's a job for the GC.
	}
}
