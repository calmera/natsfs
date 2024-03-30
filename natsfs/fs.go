package natsfs

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
	"syscall"
)

func NewFs(obs jetstream.ObjectStore) *Fs {
	return &Fs{
		obs:  obs,
		done: make(chan struct{}),
	}
}

type Fs struct {
	obs jetstream.ObjectStore

	fs.Inode
	done chan struct{}
}

func (r *Fs) Release(ctx context.Context, f fs.FileHandle) syscall.Errno {
	close(r.done)
	return syscall.F_OK
}

func (r *Fs) OnAdd(ctx context.Context) {
	go func() {
		if err := r.Watch(ctx); err != nil {
			log.Info().Err(err).Msg("watch failed")
		}
	}()
}

func (r *Fs) Watch(ctx context.Context) error {
	w, err := r.obs.Watch(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return w.Stop()
		case e := <-w.Updates():
			if e == nil {
				continue
			}

			if e.Deleted {
				r.handleDelete(ctx, e)
			} else {
				r.handlePut(ctx, e)
			}
		}
	}
}

func (r *Fs) handlePut(ctx context.Context, e *jetstream.ObjectInfo) {
	path := Path(e.Name)

	dir := path.Parent()
	base := path.LastElement()

	p := &r.Inode
	for _, pe := range dir.Elements() {
		if len(pe) == 0 {
			continue
		}

		ch := p.GetChild(pe)
		if ch == nil {
			ch = p.NewPersistentInode(ctx, &fs.Inode{}, fs.StableAttr{Mode: fuse.S_IFDIR})
			p.AddChild(pe, ch, true)
		}

		p = ch
	}

	ch := p.NewPersistentInode(ctx, &FsNode{obs: r.obs, meta: e}, fs.StableAttr{Mode: fuse.S_IFREG})
	p.AddChild(base, ch, true)
}

func (r *Fs) handleDelete(ctx context.Context, e *jetstream.ObjectInfo) {
	path := Path(e.Name)

	toDel := path.LastElement()
	pn, parent := r.Parent()
	for {
		parent.RmChild(toDel)

		// stop deleting if the parent has children
		if len(parent.Children()) > 0 {
			break
		}

		// -- the parent has no children, so we need to remove it
		toDel = pn
		pn, parent = parent.Parent()
	}
}
