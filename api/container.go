package api

import (
	"io"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// ContainerHandler encapsulates container related handlers.
type ContainerHandler struct {
	Search    web.HandlerFunc `path:"/search" auth:"container.view" desc:"search containers"`
	Find      web.HandlerFunc `path:"/find" auth:"container.view" desc:"find container by name"`
	Delete    web.HandlerFunc `path:"/delete" method:"post" auth:"container.delete" desc:"delete container"`
	FetchLogs web.HandlerFunc `path:"/fetch-logs" auth:"container.logs" desc:"fetch logs of container"`
	Connect   web.HandlerFunc `path:"/connect" auth:"*" desc:"connect to a running container"`
}

// NewContainer creates an instance of ContainerHandler
func NewContainer(b biz.ContainerBiz) *ContainerHandler {
	return &ContainerHandler{
		Search:    containerSearch(b),
		Find:      containerFind(b),
		Delete:    containerDelete(b),
		FetchLogs: containerFetchLogs(b),
		Connect:   containerConnect(b),
	}
}

func containerSearch(b biz.ContainerBiz) web.HandlerFunc {
	type Args struct {
		Node      string `json:"node" bind:"node"`
		Name      string `json:"name" bind:"name"`
		Status    string `json:"status" bind:"status"`
		PageIndex int    `json:"pageIndex" bind:"pageIndex"`
		PageSize  int    `json:"pageSize" bind:"pageSize"`
	}

	return func(ctx web.Context) (err error) {
		var (
			args       = &Args{}
			containers []*biz.Container
			total      int
		)

		if err = ctx.Bind(args); err == nil {
			containers, total, err = b.Search(args.Node, args.Name, args.Status, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(ctx, data.Map{
			"items": containers,
			"total": total,
		})
	}
}

func containerFind(b biz.ContainerBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		node := ctx.Query("node")
		id := ctx.Query("id")
		container, raw, err := b.Find(node, id)
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"container": container, "raw": raw})
	}
}

func containerDelete(b biz.ContainerBiz) web.HandlerFunc {
	type Args struct {
		Node string `json:"node"`
		ID   string `json:"id"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Delete(args.Node, args.ID, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func containerFetchLogs(b biz.ContainerBiz) web.HandlerFunc {
	type Args struct {
		Node       string `json:"node" bind:"node"`
		ID         string `json:"id" bind:"id"`
		Lines      int    `json:"lines" bind:"lines"`
		Timestamps bool   `json:"timestamps" bind:"timestamps"`
	}

	return func(ctx web.Context) (err error) {
		var (
			args           = &Args{}
			stdout, stderr string
		)
		if err = ctx.Bind(args); err == nil {
			stdout, stderr, err = b.FetchLogs(args.Node, args.ID, args.Lines, args.Timestamps)
		}
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"stdout": stdout, "stderr": stderr})
	}
}

func containerConnect(b biz.ContainerBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		var (
			node = ctx.Query("node")
			id   = ctx.Query("id")
			cmd  = ctx.Query("cmd")
		)

		_, _, err := b.Find(node, id)
		if err != nil {
			return err
		}

		conn, _, _, err := ws.UpgradeHTTP(ctx.Request(), ctx.Response())
		if err != nil {
			return err
		}

		idResp, err := b.ExecCreate(node, id, cmd)
		if err != nil {
			return err
		}

		resp, err := b.ExecAttach(node, idResp.ID)
		if err != nil {
			return err
		}

		err = b.ExecStart(node, idResp.ID)
		if err != nil {
			return err
		}

		var (
			closed   = false
			logger   = log.Get("container")
			disposer = func() {
				if !closed {
					closed = true
					_ = conn.Close()
					resp.Close()
				}
			}
		)

		// input
		go func() {
			defer disposer()

			var (
				msg []byte
				op  ws.OpCode
			)

			for {
				msg, op, err = wsutil.ReadClientData(conn)
				if err != nil {
					if !closed {
						logger.Error("failed to read data from client: ", err)
					}
					break
				}

				if op == ws.OpClose {
					break
				}

				_, err = resp.Conn.Write(msg)
				if err != nil {
					logger.Error("failed to write data to container: ", err)
					break
				}
			}
		}()

		// output
		go func() {
			defer disposer()

			var (
				n   int
				buf = make([]byte, 1024)
			)

			for {
				n, err = resp.Reader.Read(buf)
				if err == io.EOF {
					break
				} else if err != nil {
					logger.Error("failed to read data from container: ", err)
					break
				}

				err = wsutil.WriteServerMessage(conn, ws.OpText, buf[:n])
				if err != nil {
					logger.Error("failed to write data to client: ", err)
					break
				}
			}
		}()
		return nil
	}
}
