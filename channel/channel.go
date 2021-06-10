package channel

import (
	"namespace.com/dupfind/file"
)

// Direct contents of two input channels to a single output channel.
// Instead of joining the retrieved and saved channels we could instead have
// both pipes publish to the same channel.  This avoids the need for the
// join.  However, it then requires more complex logic on how to close the
// channel.  If writer A and writer B both publish to the same channel, then
// who closes the channel?  If A closes it, how do they know if B is done?
// If B closes the channel, how does it know that A is done?  Perhaps this
// could be done with a wait group.
func Join(
	in1 <-chan *file.File,
	in2 <-chan *file.File,
	ctl *Control) chan *file.File {
	out := make(chan *file.File)
	go join(in1, in2, out, ctl)
	return out
}

func join(
	in1 <-chan *file.File,
	in2 <-chan *file.File,
	out chan<- *file.File,
	ctl *Control) {

	for in1 != nil || in2 != nil {
		select {
		case <-ctl.Halt:
			in1 = nil
			in2 = nil
		case i, ok := <-in1:
			if !ok {
				in1 = nil
			} else {
				out <- i
			}
		case i, ok := <-in2:
			if !ok {
				in2 = nil
			} else {
				out <- i
			}
		}
	}
	close(out)
}

func Drain(ch <-chan *file.File) {
	for range ch {
	}
}

func DrainAll(chans ...<-chan *file.File) {
	for _, ch := range chans {
		Drain(ch)
	}
}
