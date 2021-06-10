package dup

import (
	"database/sql"

	"namespace.com/dupfind"
	"namespace.com/dupfind/cache"
	"namespace.com/dupfind/channel"
	"namespace.com/dupfind/database"
	"namespace.com/dupfind/file"
	"namespace.com/dupfind/pipe"
)

func Checksums(
	sizeCache cache.I64StrCache,
	db *sql.DB,
	batchSize int,
	ctl *channel.Control) cache.StrStrCache {

	return makeChecksumCache(sizeCache, db, batchSize, ctl).Duplicates()
}

func makeChecksumCache(
	sizeCache cache.I64StrCache,
	db *sql.DB,
	batchSize int,
	ctl *channel.Control) cache.StrStrCache {

	checksumCache := cache.MakeStrStr()

	retriever := func(f *file.File) bool {
		return database.SelectOne(db, f).Checksum != nil
	}

	// skip files which failed to get a checksum
	skipper := func(f *file.File) {}

	importer := func(files []*file.File) bool {
		return database.Import(db, files)
	}

	updater := func(f *file.File) {
		// This keys off of only checksum.  There is a tiny, non-zero chance
		// that two files can have the same checksum, but different sizes.
		// Keying off of both checksum and size would avoid this.  However,
		// there is still an even smaller chance that two files can have the
		// same checksum and same size, but not be identical.  The only way to
		// absolutely avoid this would be to compare the files directly.  But
		// that is out of scope for now.
		checksumCache.Update(*f.Checksum, f.Name)
	}

	sizes := sizeCache.ToChan(ctl)
	retrieved := make(chan *file.File)
	absent := make(chan *file.File)
	created := make(chan *file.File)
	saved := make(chan *file.File)
	joined := channel.Join(retrieved, saved, ctl)

	(&pipe.Pipe{
		NextFn:  file.FromChanFn(sizes),
		TryFn:   file.ForAllFn(retriever),
		OkFn:    file.ToChanFn(retrieved),
		NilFn:   file.ToChanFn(absent),
		DoneFn:  file.CloserFn(retrieved, absent),
		Control: ctl}).Run()

	(&pipe.Pipe{
		NextFn:  file.FromChanFn(absent),
		TryFn:   file.ForAllFn(dupfind.AddChecksum),
		OkFn:    file.ToChanFn(created),
		NilFn:   skipper,
		DoneFn:  file.CloserFn(created),
		Control: ctl}).Par().Run()

	(&pipe.Pipe{
		NextFn:  file.FromChanFn(created),
		TryFn:   importer,
		OkFn:    file.ToChanFn(saved),
		DoneFn:  file.CloserFn(saved),
		Control: ctl}).Batch(batchSize).Run()

	(&pipe.Pipe{
		NextFn:  file.FromChanFn(joined),
		OkFn:    updater,
		NilFn:   updater,
		DoneFn:  ctl.Stop,
		Control: ctl}).Run()

	<-ctl.Halt

	channel.DrainAll(sizes, absent, retrieved, created, saved)

	return checksumCache
}
