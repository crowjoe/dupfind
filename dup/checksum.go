package dup

import (
	"namespace.com/dupfind"
	"namespace.com/dupfind/cache"
	"namespace.com/dupfind/channel"
	"namespace.com/dupfind/file"
	"namespace.com/dupfind/pipe"
)

func FileSizes(dir string, ctl *channel.Control) cache.I64StrCache {
	return makeFileSizeCache(dir, ctl).Duplicates()
}

func makeFileSizeCache(dir string, ctl *channel.Control) cache.I64StrCache {
	sizes := cache.MakeI64Str()
	files := dupfind.Filenames(dir, ctl)

	updater := func(f *file.File) {
		sizes.Update(f.Size, f.Name)
	}

	(&pipe.Pipe{
		NextFn:  file.FromChanFn(files),
		OkFn:    updater,
		DoneFn:  ctl.Stop,
		Control: ctl}).Run()

	<-ctl.Halt

	channel.Drain(files)

	return sizes
}
