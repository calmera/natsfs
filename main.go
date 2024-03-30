package main

import (
	"flag"
	"fmt"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/zipfs"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Scans the arg list and sets up flags
	debug := flag.Bool("debug", false, "debug on")
	flag.Parse()
	if flag.NArg() < 1 {
		_, prog := filepath.Split(os.Args[0])
		fmt.Printf("usage: %s MOUNTPOINT\n", prog)
		os.Exit(2)
	}

	root := &zipfs.MultiZipFs{}
	sec := time.Second
	opts := fs.Options{
		EntryTimeout: &sec,
		AttrTimeout:  &sec,
	}
	opts.Debug = *debug
	server, err := fs.Mount(flag.Arg(0), root, &opts)
	if err != nil {
		fmt.Printf("Mount fail: %v\n", err)
		os.Exit(1)
	}

	server.Serve()
}
