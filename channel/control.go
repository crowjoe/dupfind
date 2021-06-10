package channel

type Control struct {
	Halt     chan struct{}
	Workers  chan struct{}
	stopping bool
}

func InitControl(numWorkers int) *Control {
	return &Control{
		Halt:    make(chan struct{}),
		Workers: make(chan struct{}, numWorkers)}
}

func (c *Control) Stop() {
	if !c.stopping {
		c.stopping = true
		close(c.Halt)
	}
}

func (c *Control) Lock() bool {
	select {
	case c.Workers <- struct{}{}:
		return true
	case <-c.Halt:
		return false
	}
}

func (c *Control) Release() {
	<-c.Workers
}

func (c *Control) Stopped() bool {
	select {
	case <-c.Halt:
		return true
	default:
		return false
	}
}
