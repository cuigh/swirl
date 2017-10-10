package docker

import (
	"context"

	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// ImageList return images on the host.
func ImageList(name string, pageIndex, pageSize int) (images []*model.ImageListInfo, totalCount int, err error) {
	ctx, cli, err := mgr.Client()
	if err != nil {
		return nil, 0, err
	}

	opts := types.ImageListOptions{}
	if name != "" {
		opts.Filters = filters.NewArgs()
		opts.Filters.Add("reference", name)
	}
	summaries, err := cli.ImageList(ctx, opts)
	if err != nil {
		return nil, 0, err
	}
	//sort.Slice(images, func(i, j int) bool {
	//	return images[i].ID < images[j].ID
	//})

	totalCount = len(summaries)
	start, end := misc.Page(totalCount, pageIndex, pageSize)
	summaries = summaries[start:end]
	if length := len(summaries); length > 0 {
		images = make([]*model.ImageListInfo, length)
		for i, summary := range summaries {
			images[i] = model.NewImageListInfo(summary)
		}
	}
	return
}

// ImageInspect returns image information.
func ImageInspect(id string) (image types.ImageInspect, raw []byte, err error) {
	var (
		ctx context.Context
		cli *client.Client
	)
	if ctx, cli, err = mgr.Client(); err == nil {
		return cli.ImageInspectWithRaw(ctx, id)
	}
	return
}

// ImageHistory returns the changes in an image in history format.
func ImageHistory(id string) (histories []image.HistoryResponseItem, err error) {
	var (
		ctx context.Context
		cli *client.Client
	)
	if ctx, cli, err = mgr.Client(); err == nil {
		return cli.ImageHistory(ctx, id)
	}
	return
}

// ImageRemove remove a image.
func ImageRemove(id string) error {
	return mgr.Do(func(ctx context.Context, cli *client.Client) (err error) {
		opts := types.ImageRemoveOptions{}
		_, err = cli.ImageRemove(ctx, id, opts)
		return
	})
}
