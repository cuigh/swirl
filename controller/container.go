package controller

import (
	"io"
	"strings"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/swirl/biz/docker"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// ContainerController is a controller of docker container
type ContainerController struct {
	List      web.HandlerFunc `path:"/" name:"container.list" authorize:"!" desc:"container list page"`
	Detail    web.HandlerFunc `path:"/:id/detail" name:"container.detail" authorize:"!" desc:"container detail page"`
	Raw       web.HandlerFunc `path:"/:id/raw" name:"container.raw" authorize:"!" desc:"container raw page"`
	Logs      web.HandlerFunc `path:"/:id/logs" name:"container.logs" authorize:"!" desc:"container logs page"`
	FetchLogs web.HandlerFunc `path:"/:id/fetch_logs" name:"container.fetch_logs" authorize:"?" desc:"fetch container logs"`
	Delete    web.HandlerFunc `path:"/delete" method:"post" name:"container.delete" authorize:"!" desc:"delete container"`
	Exec      web.HandlerFunc `path:"/:id/exec" name:"container.exec" authorize:"!" desc:"run a command in a running container"`
	Connect   web.HandlerFunc `path:"/:id/connect" name:"container.connect" authorize:"!" desc:"connect to a running container"`
}

// Container creates an instance of ContainerController
func Container() (c *ContainerController) {
	return &ContainerController{
		List:      containerList,
		Detail:    containerDetail,
		Raw:       containerRaw,
		Logs:      containerLogs,
		FetchLogs: containerFetchLogs,
		Delete:    containerDelete,
		Exec:      containerExec,
		Connect:   containerConnect,
	}
}

func containerList(ctx web.Context) error {
	args := &model.ContainerListArgs{}
	err := ctx.Bind(args)
	if err != nil {
		return err
	}
	args.PageSize = model.PageSize
	if args.PageIndex == 0 {
		args.PageIndex = 1
	}

	containers, totalCount, err := docker.ContainerList(args)
	if err != nil {
		return err
	}

	m := newPagerModel(ctx, totalCount, model.PageSize, args.PageIndex).
		Set("Name", args.Name).
		Set("Filter", args.Filter).
		Set("Containers", containers)
	return ctx.Render("container/list", m)
}

func containerDetail(ctx web.Context) error {
	id := ctx.P("id")
	container, err := docker.ContainerInspect(id)
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Container", container)
	return ctx.Render("container/detail", m)
}

func containerRaw(ctx web.Context) error {
	id := ctx.P("id")
	container, raw, err := docker.ContainerInspectRaw(id)
	if err != nil {
		return err
	}

	j, err := misc.JSONIndent(raw)
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Container", container).Set("Raw", j)
	return ctx.Render("container/raw", m)
}

func containerLogs(ctx web.Context) error {
	id := ctx.P("id")
	container, _, err := docker.ContainerInspectRaw(id)
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Container", container)
	return ctx.Render("container/logs", m)
}

func containerFetchLogs(ctx web.Context) error {
	id := ctx.P("id")
	line := cast.ToInt(ctx.Q("line"), 500)
	timestamps := cast.ToBool(ctx.Q("timestamps"), false)
	stdout, stderr, err := docker.ContainerLogs(id, line, timestamps)
	if err != nil {
		return ajaxResult(ctx, err)
	}

	return ctx.JSON(data.Map{
		"stdout": stdout.String(),
		"stderr": stderr.String(),
	})
}

func containerDelete(ctx web.Context) error {
	ids := strings.Split(ctx.F("ids"), ",")
	for _, id := range ids {
		if err := docker.ContainerRemove(id); err != nil {
			return ajaxResult(ctx, err)
		}
	}
	return ajaxSuccess(ctx, nil)
}

func containerExec(ctx web.Context) error {
	id := ctx.P("id")
	container, _, err := docker.ContainerInspectRaw(id)
	if err != nil {
		return err
	}

	m := newModel(ctx).Set("Container", container)
	return ctx.Render("container/exec", m)
}

func containerConnect(ctx web.Context) error {
	id := ctx.P("id")
	_, _, err := docker.ContainerInspectRaw(id)
	if err != nil {
		return err
	}

	conn, _, _, err := ws.UpgradeHTTP(ctx.Request(), ctx.Response())
	if err != nil {
		return err
	}

	cmd := ctx.Q("cmd")
	idResp, err := docker.ContainerExecCreate(id, cmd)
	if err != nil {
		return err
	}

	resp, err := docker.ContainerExecAttach(idResp.ID)
	if err != nil {
		return err
	}

	err = docker.ContainerExecStart(idResp.ID)
	if err != nil {
		return err
	}

	var (
		closed   = false
		logger   = log.Get("exec")
		disposer = func() {
			if !closed {
				closed = true
				conn.Close()
				resp.Close()
			}
		}
	)

	// input
	go func() {
		defer disposer()

		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				if !closed {
					logger.Error("Failed to read data from client: ", err)
				}
				break
			}

			if op == ws.OpClose {
				break
			}

			_, err = resp.Conn.Write(msg)
			if err != nil {
				logger.Error("Failed to write data to container: ", err)
				break
			}
		}
	}()

	// output
	go func() {
		defer disposer()

		buf := make([]byte, 1024)
		for {
			n, err := resp.Reader.Read(buf)
			if err == io.EOF {
				break
			} else if err != nil {
				logger.Error("Failed to read data from container: ", err)
				break
			}

			err = wsutil.WriteServerMessage(conn, ws.OpText, buf[:n])
			if err != nil {
				logger.Error("Failed to write data to client: ", err)
				break
			}
		}
	}()
	return nil
}
