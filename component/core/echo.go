package core

// Echo component pongs an arbitrary ping.
type Echo struct {
	Ping <-chan interface{}
	Pong chan<- interface{}
}

// Run reads from ping and writes to pong.
func (e *Echo) Run() {
	for p := range e.Ping {
		e.Pong <- p
	}
}
