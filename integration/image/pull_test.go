package image

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/internal/test/daemon"

	"gotest.tools/assert"
	"gotest.tools/skip"
)

func TestImagePullPlatformInvalid(t *testing.T) {
	skip.If(t, versions.LessThan(testEnv.DaemonAPIVersion(), "1.40"), "experimental in older versions")
	defer setupTest(t)()
	client := testEnv.APIClient()
	ctx := context.Background()

	_, err := client.ImagePull(ctx, "docker.io/library/hello-world:latest", types.ImagePullOptions{Platform: "foobar"})
	assert.Assert(t, err != nil)
	assert.ErrorContains(t, err, "unknown operating system or architecture")
	assert.Assert(t, errdefs.IsInvalidParameter(err))
}

func TestImagePullComparePullDuration(t *testing.T) {
	skip.If(t, testEnv.IsRemoteDaemon())

	for _, storageDriver := range []string{"aufs", "overlay2"} {
		t.Run(fmt.Sprintf("storageDriver=%s", storageDriver), func(t *testing.T) {

			var (
				durSync, durNoSync time.Duration
				args               = []string{fmt.Sprintf("--storage-driver=%s", storageDriver)}
				testImage          = "balenalib/amd64-debian:build" // should have a few layers so it actually has an impact on pull performance
			)

			d := daemon.New(t)
			client, err := d.NewClient()
			assert.NilError(t, err)

			d.Start(t, append(args, []string{fmt.Sprintf("--storage-opt=%s.sync_diffs=true", storageDriver)}...)...)
			info := d.Info(t)
			assert.Equal(t, info.Driver, storageDriver)

			ctx := context.Background()
			start := time.Now()
			_, err = client.ImagePull(ctx, testImage, types.ImagePullOptions{})
			durSync = time.Now().Sub(start)
			t.Logf("%s/syncDiffs=true took %s", t.Name(), durSync)
			assert.NilError(t, err)

			d.Stop(t)
			d.Start(t, append(args, []string{fmt.Sprintf("--storage-opt=%s.sync_diffs=false", storageDriver)}...)...)
			defer d.Stop(t)

			start = time.Now()
			_, err = client.ImagePull(ctx, testImage, types.ImagePullOptions{})
			durNoSync = time.Now().Sub(start)
			t.Logf("%s/syncDiffs=false took %s", t.Name(), durNoSync)
			assert.NilError(t, err)

			assert.Assert(t, durSync > durNoSync)
		})
	}
}
