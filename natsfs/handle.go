package natsfs

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/nats-io/nats.go/jetstream"
	"sync"
	"syscall"
)

type FileHandle struct {
	obs     jetstream.ObjectStore
	subject string

	bl sync.RWMutex
	b  []byte
}

func (h *FileHandle) Flush(ctx context.Context) syscall.Errno {
	_, err := h.obs.PutBytes(ctx, h.subject, h.b)
	if err != nil {
		return syscall.EIO
	}

	return syscall.F_OK
}

func (h *FileHandle) Write(ctx context.Context, data []byte, off int64) (written uint32, errno syscall.Errno) {
	h.bl.RLock()
	h.bl.Lock()
	defer h.bl.RUnlock()
	defer h.bl.Unlock()

	if off >= int64(len(h.b)) {
		return 0, syscall.EINVAL
	}

	if off+int64(len(data)) > int64(len(h.b)) {
		h.b = append(h.b, make([]byte, off+int64(len(data))-int64(len(h.b)))...)
	}

	_ = copy(h.b[off:], data)

	return uint32(len(data)), syscall.F_OK
}

func (h *FileHandle) Release(ctx context.Context) syscall.Errno {
	h.bl.Lock()
	defer h.bl.Unlock()

	h.b = nil
	return syscall.F_OK
}

func (h *FileHandle) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	h.bl.RLock()
	defer h.bl.RUnlock()

	if off >= int64(len(h.b)) {
		return nil, syscall.EINVAL
	}

	_ = copy(dest, h.b[off:len(dest)])

	return fuse.ReadResultData(dest), syscall.F_OK
}
