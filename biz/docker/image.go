package docker

import (
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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
