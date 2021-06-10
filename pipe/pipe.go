package pipe

import (
	"fmt"
	"sync"

	"namespace.com/dupfind/channel"
	"namespace.com/dupfind/file"
)

type Pipe struct {
	NextFn    func() *file.File
	TryFn     func([]*file.File) bool
	OkFn      func(*file.File)
	NilFn     func(*file.File)
	DoneFn    func()
	Control   *channel.Control
	Parallel  bool
	BatchSize int
}

func (p *Pipe) Par() *Pipe {
	p.Parallel = true
	return p
}

func (p *Pipe) Batch(size int) *Pipe {
	p.BatchSize = size
	return p
}

func (p *Pipe) batchSize() int {
	if p.BatchSize < 0 {
		panic(fmt.Errorf("negative batch size %d", p.BatchSize))
	} else if p.BatchSize == 0 {
		return 1
	} else {
		return p.BatchSize
	}
}

func (p *Pipe) Run() {
	go p.run()
}

func (p *Pipe) run() {
	var pending sync.WaitGroup

	if p.Parallel {
		pending.Add(1)
	}

	buffer := make([]*file.File, p.batchSize())
	i := 0

	for next := p.NextFn(); next != nil; next = p.NextFn() {
		buffer[i] = next
		i += 1

		if i >= p.batchSize() {
			if p.Parallel {
				pending.Add(1)
				go p.tryBatch(buffer[:i], &pending)
				// parallel requires a new buffer for each batch
				buffer = make([]*file.File, p.batchSize())
			} else {
				p.tryBatch(buffer[:i], nil)
			}
			i = 0
		}

	}

	// check if there is a final partial batch
	if i > 0 {
		if p.Parallel {
			pending.Add(1)
			go p.tryBatch(buffer[:i], &pending)
		} else {
			p.tryBatch(buffer[:i], nil)
		}
	}

	if p.Parallel {
		pending.Done()
		go func() {
			pending.Wait()
			p.DoneFn()
		}()
	} else {
		p.DoneFn()
	}
}

func (p *Pipe) tryBatch(files []*file.File, pending *sync.WaitGroup) {
	if p.Parallel {
		if pending != nil {
			defer pending.Done()
		}
		if p.Control.Stopped() {
			return // skipping due to early cancellation
		}
		if !p.Control.Lock() {
			return // stopped
		}
		defer p.Control.Release()
	} else {
		if p.Control.Stopped() {
			return // skipping due to early cancellation
		}
	}

	if p.TryFn == nil || p.TryFn(files) {
		for _, f := range files {
			p.OkFn(f)
		}
	} else {
		for _, f := range files {
			p.NilFn(f)
		}
	}
}
