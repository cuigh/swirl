package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/app/flag"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/data/valid"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/net/web/filter"
	"github.com/cuigh/auxo/util/run"
	_ "github.com/cuigh/swirl/api"
	"github.com/cuigh/swirl/biz"
	_ "github.com/cuigh/swirl/dao/bolt"
	_ "github.com/cuigh/swirl/dao/mongo"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/scaler"
)

var (
	//go:embed ui/dist
	webFS embed.FS
)

func main() {
	app.Name = "Swirl"
	app.Version = "1.0.0rc2"
	app.Desc = "A web management UI for Docker, focused on swarm cluster"
	app.Action = func(ctx *app.Context) error {
		return run.Pipeline(misc.LoadOptions, initSystem, scaler.Start, startServer)
	}
	app.Flags.Register(flag.All)
	app.Start()
}

func startServer() (err error) {
	s := web.Auto()
	s.Validator = &valid.Validator{}
	s.ErrorHandler.Default = handleError
	s.Use(filter.NewRecover())
	s.Static("/", http.FS(loadWebFS()), "index.html")

	const prefix = "api."
	g := s.Group("/api", findFilters("identifier", "authorizer")...)
	container.Range(func(name string, service interface{}) bool {
		if strings.HasPrefix(name, prefix) {
			g.Handle("/"+name[len(prefix):], service)
		}
		return true
	})

	app.Run(s)
	return
}

func loadWebFS() fs.FS {
	sub, err := fs.Sub(webFS, "ui/dist")
	if err != nil {
		panic(err)
	}
	return sub
}

func handleError(ctx web.Context, err error) {
	var (
		status       = http.StatusInternalServerError
		code   int32 = 1
	)

	if e, ok := err.(*web.Error); ok {
		status = e.Status()
	}
	if e, ok := err.(*errors.CodedError); ok {
		code = e.Code
	}

	err = ctx.Status(status).Result(code, err.Error(), nil)
	if err != nil {
		ctx.Logger().Error(err)
	}
}

func findFilters(names ...string) []web.Filter {
	var filters []web.Filter
	for _, name := range names {
		filters = append(filters, container.Find(name).(web.Filter))
	}
	return filters
}

func initSystem() error {
	return container.Call(func(b biz.SystemBiz) error {
		ctx, cancel := misc.Context(time.Minute)
		defer cancel()

		return b.Init(ctx)
	})
}

func loadSetting(sb biz.SettingBiz) *misc.Setting {
	var (
		err  error
		opts data.Map
		b    []byte
		s    = &misc.Setting{}
	)

	ctx, cancel := misc.Context(30 * time.Second)
	defer cancel()

	if opts, err = sb.Load(ctx); err == nil {
		if b, err = json.Marshal(opts); err == nil {
			err = json.Unmarshal(b, s)
		}
	}
	if err != nil {
		log.Get("misc").Error("failed to load setting: ", err)
	}
	return s
}

func init() {
	container.Put(loadSetting)
}
