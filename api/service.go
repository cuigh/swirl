package api

import (
	"net/http"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/misc"
)

// ServiceHandler encapsulates service related handlers.
type ServiceHandler struct {
	Search    web.HandlerFunc `path:"/search" auth:"service.view" desc:"search services"`
	Find      web.HandlerFunc `path:"/find" auth:"service.view" desc:"find service by name"`
	Delete    web.HandlerFunc `path:"/delete" method:"post" auth:"service.delete" desc:"delete service"`
	Restart   web.HandlerFunc `path:"/restart" method:"post" auth:"service.restart" desc:"restart service"`
	Rollback  web.HandlerFunc `path:"/rollback" method:"post" auth:"service.rollback" desc:"rollback service"`
	Scale     web.HandlerFunc `path:"/scale" method:"post" auth:"service.edit" desc:"scale service"`
	Save      web.HandlerFunc `path:"/save" method:"post" auth:"service.edit" desc:"create or update service"`
	Deploy    web.HandlerFunc `path:"/deploy" method:"post" auth:"service.deploy" desc:"deploy service"`
	FetchLogs web.HandlerFunc `path:"/fetch-logs" auth:"service.logs" desc:"fetch logs of service"`
}

// NewService creates an instance of ServiceHandler
func NewService(b biz.ServiceBiz) *ServiceHandler {
	return &ServiceHandler{
		Search:    serviceSearch(b),
		Find:      serviceFind(b),
		Delete:    serviceDelete(b),
		Restart:   serviceRestart(b),
		Rollback:  serviceRollback(b),
		Scale:     serviceScale(b),
		Save:      serviceSave(b),
		Deploy:    serviceDeploy(b),
		FetchLogs: serviceFetchLogs(b),
	}
}

func serviceSearch(b biz.ServiceBiz) web.HandlerFunc {
	type Args struct {
		Name      string `json:"name" bind:"name"`
		Mode      string `json:"mode" bind:"mode"`
		PageIndex int    `json:"pageIndex" bind:"pageIndex"`
		PageSize  int    `json:"pageSize" bind:"pageSize"`
	}

	return func(c web.Context) (err error) {
		var (
			args     = &Args{}
			services []*biz.ServiceBase
			total    int
		)

		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			services, total, err = b.Search(ctx, args.Name, args.Mode, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(c, data.Map{
			"items": services,
			"total": total,
		})
	}
}

func serviceFind(b biz.ServiceBiz) web.HandlerFunc {
	return func(c web.Context) error {
		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		name := c.Query("name")
		status := c.Query("status") == "true"
		service, raw, err := b.Find(ctx, name, status)
		if err != nil {
			return err
		} else if service == nil {
			return web.NewError(http.StatusNotFound)
		}
		return success(c, data.Map{"service": service, "raw": raw})
	}
}

func serviceDelete(b biz.ServiceBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Delete(ctx, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func serviceRestart(b biz.ServiceBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Restart(ctx, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func serviceRollback(b biz.ServiceBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Rollback(ctx, args.Name, c.User())
		}
		return ajax(c, err)
	}
}

func serviceScale(b biz.ServiceBiz) web.HandlerFunc {
	type Args struct {
		Name    string `json:"name"`
		Count   uint64 `json:"count"`
		Version uint64 `json:"version"`
	}
	return func(c web.Context) (err error) {
		args := &Args{}
		if err = c.Bind(args); err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			err = b.Scale(ctx, args.Name, args.Count, args.Version, c.User())
		}
		return ajax(c, err)
	}
}

func serviceSave(b biz.ServiceBiz) web.HandlerFunc {
	return func(c web.Context) error {
		s := &biz.Service{}
		err := c.Bind(s, true)
		if err == nil {
			ctx, cancel := misc.Context(defaultTimeout)
			defer cancel()

			if s.ID == "" {
				err = b.Create(ctx, s, c.User())
			} else {
				err = b.Update(ctx, s, c.User())
			}
		}
		return ajax(c, err)
	}
}

func serviceDeploy(b biz.ServiceBiz) web.HandlerFunc {
	return func(c web.Context) error {
		//mode := c.Query("mode") // update/replace
		service := &biz.Service{}
		err := c.Bind(service, true)
		if err != nil {
			return err
		}

		ctx, cancel := misc.Context(defaultTimeout)
		defer cancel()

		var s *biz.Service
		if s, _, err = b.Find(ctx, service.Name, false); err != nil {
			return err
		}

		if s == nil {
			err = b.Create(ctx, service, c.User())
		} else {
			// Only the image field is allowed to be changed when updating.
			s.Image = service.Image
			err = b.Update(ctx, s, c.User())
		}
		return ajax(c, err)
	}
}

func serviceFetchLogs(b biz.ServiceBiz) web.HandlerFunc {
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
