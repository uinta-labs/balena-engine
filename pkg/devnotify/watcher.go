package devnotify

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/docker/docker/container"
	"github.com/docker/docker/layer"
)

// Watcher is responsible for mirroring and syncing the host devfs to the
// location passed to `Prepare`.
type Watcher interface {
	Prepare(layer.RWLayer) error
	Start() error
	Stop() error
}

func NewWatcher(c *container.Container) Watcher {
	return &watcher{containerID: c.ID}
}

type watcher struct {
	ctx         context.Context
	cancelfn    context.CancelFunc
	containerID string
	path        string
}

// Prepare sets up a mirror of the current devfs at `path`. This does not start
// syncing any changes.
func (w *watcher) Prepare(rwLayer layer.RWLayer) error {
	fs, err := rwLayer.Mount("")
	if err != nil {
		return err
	}

	w.path = fs.Path()
	logger := logrus.WithFields(logrus.Fields{
		"balenaext": "devfs-watcher",
		"container": w.containerID,
	})
	err = CloneTree(logger, w.path)
	if err != nil {
		return err
	}

	return rwLayer.Unmount()
}

// Start syncs changes from `/dev` to the destination that was passed to
// `Prepare`.
//
// It doesn't block and exits right away.
func (w *watcher) Start() error {
	logger := logrus.WithFields(logrus.Fields{
		"balenaext": "devfs-watcher",
		"container": w.containerID,
	})

	if _, err := os.Stat(w.path); err != nil {
		return err
	}

	ctx := context.TODO()
	w.ctx, w.cancelfn = context.WithCancel(ctx)

	logger.WithField("target-dir", w.path).Warn("Watching host devfs")
	go SyncTree(w.ctx, logger, w.path)

	return nil
}

// Stop cancels the running watcher.
//
// It doesn't take care of cleaning up the sync destination.
func (w *watcher) Stop() error {
	w.cancelfn()
	return w.ctx.Err()
}
