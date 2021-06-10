package dupfind

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"

	"namespace.com/dupfind/file"
)

func AddChecksum(f *file.File) bool {
	if f.Checksum != nil {
		panic(fmt.Errorf(
			"file %s already has a checksum (%s)",
			f.Name,
			*f.Checksum))
	}

	file, err := os.Open(f.Name)
	if err != nil {
		// Could be due to insufficient permission, etc.
		log.Print(err)
		return false
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		panic(err)
	}

	bytes := h.Sum(nil)
	checksum := fmt.Sprintf("%x", bytes)
	f.Checksum = &checksum
	return true
}
