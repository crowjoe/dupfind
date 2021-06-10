package file

import (
	"fmt"
)

type File struct {
	Name     string
	Size     int64
	Checksum *string
}

func (f File) String() string {
	return fmt.Sprintf(
		"Result{Name: %s, Size: %d, Checksum: %s}",
		f.Name,
		f.Size,
		f.checksumString())
}

func (f File) checksumString() string {
	if f.Checksum == nil {
		return "nil"
	} else {
		return *f.Checksum
	}
}

func CloserFn(chans ...chan<- *File) func() {
	return func() {
		for _, ch := range chans {
			close(ch)
		}
	}
}

func FromChanFn(ch <-chan *File) func() *File {
	return func() *File {
		if next, ok := <-ch; ok {
			return next
		}
		return nil
	}
}

func ToChanFn(ch chan<- *File) func(*File) {
	return func(f *File) {
		ch <- f
	}
}

func ForAllFn(predicate func(*File) bool) func([]*File) bool {
	return func(files []*File) bool {
		for _, f := range files {
			if !predicate(f) {
				return false
			}
		}
		return true
	}
}
