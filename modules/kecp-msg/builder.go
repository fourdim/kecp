package kecpmsg

func NewListMsg(list []string) *Message {
	return &Message{
		Type:    List,
		Payload: list,
	}
}

func NewJoinMsg(name string, clientKey string) *Message {
	return &Message{
		Type:            Join,
		Payload:         name,
		ExceptClientKey: clientKey,
	}
}

func NewLeaveMsg(name string, clientKey string) *Message {
	return &Message{
		Type:            Leave,
		Payload:         name,
		ExceptClientKey: clientKey,
	}
}

func NewErrorMsg(err error) *Message {
	return &Message{
		Type:    Error,
		Payload: err.Error(),
	}
}
