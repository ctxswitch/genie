package sinks

type Sink interface {
	Send([]byte)
	Connect()
	Init()
}
