package dupfind

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"namespace.com/dupfind/channel"
	"namespace.com/dupfind/file"
)

func Filenames(path string, ctl *channel.Control) chan *file.File {
	var pending sync.WaitGroup
	out := make(chan *file.File)

	if info := fileInfo(path, ctl); info != nil {
		pending.Add(1)
		go files(path, info, out, &pending, ctl)

		go func() {
			pending.Wait()
			close(out)
		}()

	} else {
		close(out)
	}

	return out
}

func fileInfo(path string, ctl *channel.Control) os.FileInfo {
	if !ctl.Lock() {
		return nil // stopped
	}
	defer ctl.Release()

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		panic(err)
	}

	return info
}

func files(
	path string,
	info os.FileInfo,
	out chan<- *file.File,
	pending *sync.WaitGroup,
	ctl *channel.Control) {

	defer pending.Done()

	if ctl.Stopped() {
		return // skipping due to early cancellation
	}

	if info.IsDir() {
		entries, err := dirEntries(path, ctl)
		if err != nil {
			log.Print(err)
			return
		}

		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				log.Print(err)
			} else {
				joined := filepath.Join(path, entry.Name())
				pending.Add(1)
				go files(joined, info, out, pending, ctl)
			}
		}

	} else if info.Mode().IsRegular() {
		absFilename, err := filepath.Abs(path)
		if err != nil {
			log.Print(err)
		}
		out <- &file.File{Name: absFilename, Size: info.Size()}
	} // anything else would be a non-regular file and will be skipped
}

func dirEntries(path string, ctl *channel.Control) ([]os.DirEntry, error) {
	if !ctl.Lock() {
		return nil, nil // stopped
	}
	defer ctl.Release()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Supposedly, ReadDir is supposed to be more efficient than Readdir
	return file.ReadDir(0)
}
