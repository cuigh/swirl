package api

import (
	"net/http"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/misc"
)

// TaskHandler encapsulates node related handlers.
type TaskHandler struct {
	Search    web.HandlerFunc `path:"/search" auth:"task.view" desc:"search tasks"`
	Find      web.HandlerFunc `path:"/find" auth:"task.view" desc:"find task by id"`
	FetchLogs web.HandlerFunc `path:"/fetch-logs" auth:"task.logs" desc:"fetch logs of task"`
}

// NewTask creates an instance of TaskHandler
func NewTask(b biz.TaskBiz) *TaskHandler {
	return &TaskHandler{
		Search:    taskSearch(b),
		Find:      taskFind(b),
		FetchLogs: taskFetchLogs(b),
	}
}

func taskSearch(b biz.TaskBiz) web.HandlerFunc {
	type Args struct {
		Service   string `json:"service" bind:"service"`
		State     string `json:"state" bind:"state"`
		PageIndex int    `json:"pageIndex" bind:"pageIndex"`
		PageSize  int    `json:"pageSize" bind:"pageSize"`
	}

	return func(c web.Context) (err error) {
		var (
			args  = &Args{}
			tasks []*biz.Task
			total int
		)

		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			tasks, total, err = b.Search(ctx, "", args.Service, args.State, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(c, data.Map{
			"items": tasks,
			"total": total,
		})
	}
}

func taskFind(b biz.TaskBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		id := c.Query("id")
		task, raw, err := b.Find(ctx, id)
		if err != nil {
			return err
		} else if task == nil {
			return web.NewError(http.StatusNotFound)
		}
		return success(c, data.Map{"task": task, "raw": raw})
	}
}

func taskFetchLogs(b biz.TaskBiz) web.HandlerFunc {
	type Args struct {
		ID         string `json:"id" bind:"id"`
		Lines      int    `json:"lines" bind:"lines"`
		Timestamps bool   `json:"timestamps" bind:"timestamps"`
	}

	return func(c web.Context) (err error) {
		var (
			args           = &Args{}
			stdout, stderr string
		)
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			stdout, stderr, err = b.FetchLogs(ctx, args.ID, args.Lines, args.Timestamps)
		}
		if err != nil {
			return err
		}
		return success(c, data.Map{"stdout": stdout, "stderr": stderr})
	}
}
