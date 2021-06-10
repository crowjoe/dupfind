package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"namespace.com/dupfind/channel"
	"namespace.com/dupfind/database"
	"namespace.com/dupfind/dup"
)

func main() {
	var start = time.Now()
	flags := log.Ldate | log.LUTC | log.Ltime | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(flags)

	var batchSizeArg = flag.Int("batch-size", 100, "batch size per db transaction")
	var dbArg = flag.String("db", "", "database cache file (required)")
	var dirArg = flag.String("dir", ".", "directory to scan")
	var workersArg = flag.Int("workers", -1, "max worker threads, -1 to match #cpus")

	flag.Parse()

	if *workersArg <= 0 {
		*workersArg = runtime.GOMAXPROCS(-1)
	}

	var fileSizeCtl = channel.InitControl(*workersArg)
	var checksumCtl = channel.InitControl(*workersArg)

	go func() {
		os.Stdin.Read(make([]byte, 1))
		fileSizeCtl.Stop()
		checksumCtl.Stop()
	}()

	db := dbConnect(*dbArg)
	defer db.Close()

	database.Create(db) // Create the DB if it doesn't exist

	fileSizes := dup.FileSizes(*dirArg, fileSizeCtl)
	checksums := dup.Checksums(fileSizes, db, *batchSizeArg, checksumCtl)

	log.Printf("duplicate checksums:\n%s", checksums.String())
	log.Printf("elapsed time %s", time.Since(start))
}

func dbConnect(dbFile string) *sql.DB {
	if dbFile == "" {
		fmt.Fprintln(os.Stderr, "missing db arg")
		flag.Usage()
		os.Exit(1)
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	return db
}
