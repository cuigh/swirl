package api

import (
	"net/http"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
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

	return func(ctx web.Context) (err error) {
		var (
			args     = &Args{}
			services []*biz.ServiceBase
			total    int
		)

		if err = ctx.Bind(args); err == nil {
			services, total, err = b.Search(args.Name, args.Mode, args.PageIndex, args.PageSize)
		}

		if err != nil {
			return
		}

		return success(ctx, data.Map{
			"items": services,
			"total": total,
		})
	}
}

func serviceFind(b biz.ServiceBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		name := ctx.Query("name")
		status := ctx.Query("status") == "true"
		service, raw, err := b.Find(name, status)
		if err != nil {
			return err
		} else if service == nil {
			return web.NewError(http.StatusNotFound)
		}
		return success(ctx, data.Map{"service": service, "raw": raw})
	}
}

func serviceDelete(b biz.ServiceBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Delete(args.Name, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func serviceRestart(b biz.ServiceBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Restart(args.Name, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func serviceRollback(b biz.ServiceBiz) web.HandlerFunc {
	type Args struct {
		Name string `json:"name"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Rollback(args.Name, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func serviceScale(b biz.ServiceBiz) web.HandlerFunc {
	type Args struct {
		Name    string `json:"name"`
		Count   uint64 `json:"count"`
		Version uint64 `json:"version"`
	}
	return func(ctx web.Context) (err error) {
		args := &Args{}
		if err = ctx.Bind(args); err == nil {
			err = b.Scale(args.Name, args.Count, args.Version, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func serviceSave(b biz.ServiceBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		s := &biz.Service{}
		err := ctx.Bind(s, true)
		if err == nil {
			if s.ID == "" {
				err = b.Create(s, ctx.User())
			} else {
				err = b.Update(s, ctx.User())
			}
		}
		return ajax(ctx, err)
	}
}

func serviceDeploy(b biz.ServiceBiz) web.HandlerFunc {
	return func(ctx web.Context) error {
		//mode := ctx.Query("mode") // update/replace
		service := &biz.Service{}
		err := ctx.Bind(service, true)
		if err != nil {
			return err
		}

		var s *biz.Service
		if s, _, err = b.Find(service.Name, false); err != nil {
			return err
		}

		if s == nil {
			err = b.Create(service, ctx.User())
		} else {
			err = b.Update(service, ctx.User())
		}
		return ajax(ctx, err)
	}
}

func serviceFetchLogs(b biz.ServiceBiz) web.HandlerFunc {
	type Args struct {
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
			stdout, stderr, err = b.FetchLogs(args.ID, args.Lines, args.Timestamps)
		}
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"stdout": stdout, "stderr": stderr})
	}
}
