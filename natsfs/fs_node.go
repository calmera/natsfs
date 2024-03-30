package natsfs

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/nats-io/nats.go/jetstream"
	"syscall"
)

type FsNode struct {
	obs  jetstream.ObjectStore
	meta *jetstream.ObjectInfo

	fs.Inode
}

func (n *FsNode) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	// -- load the data from the object store
	b, err := n.obs.GetBytes(ctx, n.meta.Name)
	if err != nil {
		return nil, 0, syscall.EIO
	}

	fh = &FileHandle{
		obs:     n.obs,
		subject: n.meta.Name,
		b:       b,
	}

	return fh, 0, syscall.F_OK
}

func (n *FsNode) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if n.meta == nil {
		out.Mode = fuse.S_IFDIR | 0755
		return syscall.F_OK
	} else {
		out.Mode = fuse.S_IFREG | 0644
		out.Size = n.meta.Size
		return syscall.F_OK
	}
}
