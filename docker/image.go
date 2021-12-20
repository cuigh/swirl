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
func (d *Docker) ImageList(ctx context.Context, node, name string, pageIndex, pageSize int) (images []types.ImageSummary, total int, err error) {
	c, err := d.agent(node)
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
func (d *Docker) ImageInspect(ctx context.Context, node, id string) (image types.ImageInspect, raw []byte, err error) {
	var c *client.Client
	if c, err = d.agent(node); err == nil {
		return c.ImageInspectWithRaw(ctx, id)
	}
	return
}

// ImageHistory returns the changes in an image in history format.
func (d *Docker) ImageHistory(ctx context.Context, node, id string) (histories []image.HistoryResponseItem, err error) {
	var c *client.Client
	if c, err = d.agent(node); err == nil {
		return c.ImageHistory(ctx, id)
	}
	return
}

// ImageRemove remove a image.
func (d *Docker) ImageRemove(ctx context.Context, node, id string) error {
	c, err := d.agent(node)
	if err == nil {
		_, err = c.ImageRemove(ctx, id, types.ImageRemoveOptions{})
	}
	return err
}
