package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/calmera/natsfs/natsfs"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Scans the arg list and sets up flags
	debug := flag.Bool("debug", false, "debug on")
	url := flag.String("url", nats.DefaultURL, "The nats server URL")
	obsBucket := flag.String("bucket", "natsfs", "The jetstream bucket to use")
	jwt := flag.String("jwt", "", "The JWT token to use")
	seed := flag.String("seed", "", "The seed to use")
	creds := flag.String("creds", "", "The credentials file to use")

	flag.Parse()

	if flag.NArg() < 1 {
		_, prog := filepath.Split(os.Args[0])
		fmt.Printf("usage: %s MOUNTPOINT\n", prog)
		os.Exit(2)
	}

	hn, err := os.Hostname()
	if err != nil {
		hn = "unknown"
	}

	nopts := []nats.Option{
		nats.Name(fmt.Sprintf("NATSFS-%s", hn)),
	}

	if creds != nil {
		nopts = append(nopts, nats.UserCredentials(*creds))
	}

	if jwt != nil && seed != nil {
		nopts = append(nopts, nats.UserCredentials(*jwt, *seed))
	}

	nc, err := nats.Connect(*url, nopts...)
	if err != nil {
		log.Panic().Err(err).Msg("failed to connect to nats")
		os.Exit(1)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create jetstream context")
		os.Exit(1)
	}

	obs, err := js.ObjectStore(context.Background(), *obsBucket)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create object store")
		os.Exit(1)
	}

	root := natsfs.NewFs(obs)
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
