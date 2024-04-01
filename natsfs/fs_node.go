package natsfs

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
	"syscall"
)

type FsNode struct {
	obs  jetstream.ObjectStore
	meta *jetstream.ObjectInfo

	fs.Inode
}

func (n *FsNode) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	out.Mode = in.Mode
	out.Size = in.Size
	return syscall.F_OK
}

func (n *FsNode) Unlink(ctx context.Context, name string) syscall.Errno {
	// -- check if the target node is a directory
	child, fnd := n.Children()[name]
	if !fnd {
		return syscall.ENOENT
	}

	if child.Mode()&syscall.S_IFDIR != 0 {
		return syscall.F_OK
	}

	subject := n.meta.Name + "/" + name
	if err := n.obs.Delete(ctx, subject); err != nil {
		log.Error().Err(err).Msg("unable to delete file")
		return syscall.EIO
	}

	return syscall.F_OK
}

func (n *FsNode) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (node *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	if mode%syscall.S_IFREG == 0 {
		subject := n.meta.Name + "/" + name

		fsNode := &FsNode{obs: n.obs, meta: nil}
		ch := n.NewPersistentInode(ctx, fsNode, fs.StableAttr{Mode: fuse.S_IFREG})

		// -- create an empty file in the object store
		om := jetstream.ObjectMeta{
			Name: subject,
			Metadata: map[string]string{
				"inode": fmt.Sprintf("%d", ch.StableAttr().Ino),
			},
		}

		oi, err := n.obs.Put(ctx, om, &bytes.Buffer{})
		if err != nil {
			log.Error().Err(err).Msg("unable to create file")
			return nil, nil, 0, syscall.EIO
		}
		fsNode.meta = oi

		fh := &FileHandle{
			obs:     n.obs,
			subject: subject,
		}

		n.AddChild(name, ch, true)

		return ch, fh, 0, syscall.F_OK
	} else if mode%syscall.S_IFDIR == 0 {
		ch := n.NewPersistentInode(ctx, &fs.Inode{}, fs.StableAttr{Mode: mode})

		n.AddChild(name, ch, true)

		return ch, nil, 0, syscall.F_OK
	} else {
		return nil, nil, 0, syscall.ENOTSUP
	}
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
