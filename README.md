# dupfind: the duplicate file finder

### Summary

The dupfind utility is designed to recursively scan a given directory and identify all duplicate files.

### Goals

* Gain experience with GoLang by working on a simple, but non-trivial application.  
* Demonstrate knowledge of common GoLang features such as channels and goroutines.  
* Optimize parallel processing performance.  
* Experiment with reducing duplicate code by passing functions as arguments.  
* Attempt to follow established design patterns and coding conventions while still developing personal GoLang style preferences.
* Contribute to the open software community.

### Usage

```
Usage of dupfind:
  -batch-size int
    	batch size per db transaction (default 100)
  -db string
    	database cache file (required)
  -dir string
    	directory to scan (default ".")
  -workers int
    	max worker threads, -1 to match #cpus (default -1)
```

Examples

```
dupfind -db /tmp/mydb.dat

dupfind -db ./db -dir ~/bin

dupfind -db db.dat -batch-size 10

dupfind -db `mktemp` -workers 1

dupfind -db foo -dir ~ -batch-size 9 -workers 2
```

### Design

1. Part 1 - Build a collection of files having the same size.
    1. Traverse the file system under the given directory.  The names of any files that are encountered are stored in a File struct along with the size of the file.  The File record is written to a downstream channel.
    1. The file names and sizes are consumed from the channel and cached in a map.  The file size is used as the key to the map.  The names of files which have the same size are stored together as a single value within the map under the given size key.  Multiple files with the same size are combined as a linked list.  In this case, linked lists have the benefit of low memory overhead, fast traversal, and easy to prepend to.  Duplicates should be rare, but if there are an extremly high number of duplicates having the same size, we can quickly prepend to the map value instead of needing to allocate or copy large buffers.
    1. A cache map of file sizes and file names has now been built.  Create a copy of the cache which *only* contains key-value pairs where the linked value has at least 2 elements.  These items could be duplicates.
1. Part 2 - Build a collection of files having the same checksum.
    1. Traverse the collection of files with the same size and write each File record to a downstream channel.
    1. For each File record, attempt to load its checksum from our database file.
        1. If found, write the updated File record to the retrieved channel.
        1. If not found, write the File record to the absent channel.
    1. Consume missing files from the absent channel.  Calculate their checksum and update the File record.  Write the updated File record to the created channel.
    1. Consume files with newly calculated checksum values from the created channel.  Save them in batches to the database.  Write the File records to the saved channel.
    1. Consume files from both the retrieved channel and saved channel.  Write them both to a joined channel.  (Combining multiple channels into one is sometimes to be avoided and it is recommended to have the upstream producers just write to the same channel in the first place.  But this makes shutdown orchestration more difficult.  It was a calculated choice to have separate channels and then join them like this.)
    1. Consume files from the joined channel and add them to the map cache which is keyed by checksum string.  Files having the same checksum will be added to the same value.  As with the size map, the value is a linked list.  File names added to this list are prepended.
    1. A new map of file checksums and file names has now been built.  We were able to build it faster because we prefiltered for files where there was at least one other file with the same size.  Now we can filter this checksum map for records where there are at least two files having the same checksum.
    1. The filtered checksum map now contains only items where at least two files have the same checksum.  These files are almost certainly duplicates.  Print these files.

### Future enhancements

* Add TTL expirations to DB records
* Support other kinds of checksums (sha, etc)
* Add logic which actually compares files to prove they are identical and not merely having the same checksum.
* Support a streaming scenario where there is a neverending channel of incoming files to check.  (Requires a solution where caches might never be complete.)
* Once GoLang supports generics, update library packages to use them.
* Add more unit tests
