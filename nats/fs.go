package nats

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"syscall"
)

func NewFs(nc *nats.Conn, js jetstream.JetStream) (*Node, error) {
	return &Node{}, nil
}

type Node struct {
	obs jetstream.ObjectStore

	fs.Inode
}

func (n *Node) Release(ctx context.Context, f fs.FileHandle) syscall.Errno {
	//TODO implement me
	panic("implement me")
}

func (n *Node) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (node *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) Read(ctx context.Context, f fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	n.obs.GetBytes(ctx)
	//TODO implement me
	panic("implement me")
}

func (n *Node) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	n.obs.GetInfo()
	//TODO implement me
	panic("implement me")
}

func (n *Node) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	//TODO implement me
	panic("implement me")
}

func (n *Node) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.ENOENT
}
