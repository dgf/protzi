package component

// Echo component pongs an arbitrary ping back.
type Echo struct {
	Ping <-chan interface{}
	Pong chan<- interface{}
}

func (e *Echo) Run() {
	for p := range e.Ping {
		e.Pong <- p
	}
}
