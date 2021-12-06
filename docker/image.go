package docker

import (
	"context"

	"github.com/cuigh/swirl/misc"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// ImageList return images on the host.
func (d *Docker) ImageList(ctx context.Context, name string, pageIndex, pageSize int) (images []types.ImageSummary, total int, err error) {
	c, err := d.client()
	if err != nil {
		return nil, 0, err
	}

	opts := types.ImageListOptions{}
	if name != "" {
		opts.Filters = filters.NewArgs()
		opts.Filters.Add("reference", name)
	}
	images, err = c.ImageList(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	total = len(images)
	start, end := misc.Page(total, pageIndex, pageSize)
	images = images[start:end]
	return
}

// ImageInspect returns image information.
func (d *Docker) ImageInspect(ctx context.Context, id string) (image types.ImageInspect, raw []byte, err error) {
	var c *client.Client
	if c, err = d.client(); err == nil {
		return c.ImageInspectWithRaw(ctx, id)
	}
	return
}

// ImageHistory returns the changes in an image in history format.
func (d *Docker) ImageHistory(ctx context.Context, id string) (histories []image.HistoryResponseItem, err error) {
	var c *client.Client
	if c, err = d.client(); err == nil {
		return c.ImageHistory(ctx, id)
	}
	return
}

// ImageRemove remove a image.
func (d *Docker) ImageRemove(ctx context.Context, id string) error {
	return d.call(func(c *client.Client) (err error) {
		_, err = c.ImageRemove(ctx, id, types.ImageRemoveOptions{})
		return
	})
}
