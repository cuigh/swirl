package api

import (
	"io"
	"net/http"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/misc"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// ContainerHandler encapsulates container related handlers.
type ContainerHandler struct {
	Search    web.HandlerFunc `path:"/search" auth:"container.view" desc:"search containers"`
	Find      web.HandlerFunc `path:"/find" auth:"container.view" desc:"find container by name"`
	Delete    web.HandlerFunc `path:"/delete" method:"post" auth:"container.delete" desc:"delete container"`
	FetchLogs web.HandlerFunc `path:"/fetch-logs" auth:"container.logs" desc:"fetch logs of container"`
	Connect   web.HandlerFunc `path:"/connect" auth:"container.execute" desc:"connect to a running container"`
	Prune     web.HandlerFunc `path:"/prune" method:"post" auth:"container.delete" desc:"delete unused containers"`
}

// NewContainer creates an instance of ContainerHandler
func NewContainer(b biz.ContainerBiz) *ContainerHandler {
	return &ContainerHandler{
		Search:    containerSearch(b),
		Find:      containerFind(b),
		Delete:    containerDelete(b),
		FetchLogs: containerFetchLogs(b),
		Connect:   containerConnect(b),
		Prune:     containerPrune(b),
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

	return func(c web.Context) (err error) {
		var (
			args       = &Args{}
			containers []*biz.Container
			total      int
		)

		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			containers, total, err = b.Search(ctx, args.Node, args.Name, args.Status, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(c, data.Map{
			"items": containers,
			"total": total,
		})
	}
}

func containerFind(b biz.ContainerBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		node := c.Query("node")
		id := c.Query("id")
		container, raw, err := b.Find(ctx, node, id)
		if err != nil {
			return err
		} else if container == nil {
			return web.NewError(http.StatusNotFound)
		}
		return success(c, data.Map{"container": container, "raw": raw})
	}
}

func containerDelete(b biz.ContainerBiz) web.HandlerFunc {
	type Args struct {
		Node string `json:"node"`
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Delete(ctx, args.Node, args.ID, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func containerFetchLogs(b biz.ContainerBiz) web.HandlerFunc {
	type Args struct {
		Node       string `json:"node" bind:"node"`
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

			stdout, stderr, err = b.FetchLogs(ctx, args.Node, args.ID, args.Lines, args.Timestamps)
		}
		if err != nil {
			return err
		}
		return success(c, data.Map{"stdout": stdout, "stderr": stderr})
	}
}

func containerConnect(b biz.ContainerBiz) web.HandlerFunc {
	return func(c web.Context) error {
		var (
			node = c.Query("node")
			id   = c.Query("id")
			cmd  = c.Query("cmd")
		)

		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		container, _, err := b.Find(ctx, node, id)
		if err != nil {
			return err
		} else if container == nil {
			return web.NewError(http.StatusNotFound)
		}

		conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
		if err != nil {
			return err
		}

		idResp, err := b.ExecCreate(ctx, node, id, cmd)
		if err != nil {
			return err
		}

		resp, err := b.ExecAttach(ctx, node, idResp.ID)
		if err != nil {
			return err
		}

		err = b.ExecStart(ctx, node, idResp.ID)
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

func containerPrune(b biz.ContainerBiz) web.HandlerFunc {
	type Args struct {
		Node string `json:"node"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err != nil {
			return err
		}

		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		count, size, err := b.Prune(ctx, args.Node, c.User())
		if err != nil {
			return err
		}

		return success(c, data.Map{
			"count": count,
			"size":  size,
		})
	}
}
