package sinks

type Sink interface {
	Send(event []byte)
	Connect()
	Init()
}
