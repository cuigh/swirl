package biz

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mapToOptions(m map[string]string) (opts data.Options) {
	if len(m) == 0 {
		return
	}

	opts = data.Options{}
	for k, v := range m {
		opts = append(opts, data.Option{Name: k, Value: v})
	}
	sort.Slice(opts, func(i, j int) bool {
		return opts[i].Name < opts[j].Name
	})
	return opts
}

func envToOptions(env []string) (opts data.Options) {
	if len(env) == 0 {
		return
	}

	opts = make(data.Options, len(env))
	for i, e := range env {
		opts[i] = data.ParseOption(e, "=")
	}
	sort.Slice(opts, func(i, j int) bool {
		return opts[i].Name < opts[j].Name
	})
	return opts
}

func toEnv(opts data.Options) (env []string) {
	if len(opts) > 0 {
		env = make([]string, len(opts))
		for i, opt := range opts {
			env[i] = opt.Name + "=" + opt.Value
		}
		sort.Strings(env)
	}
	return
}

func toMap(opts data.Options) (m map[string]string) {
	if len(opts) == 0 {
		return
	}

	m = make(map[string]string)
	for _, opt := range opts {
		m[opt.Name] = opt.Value
	}
	return
}

func parseArgs(args string) []string {
	if args == "" {
		return nil
	}
	return strings.Split(args, " ")
}

func formatTime(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05")
}

// generate 8-chars short id, only suitable for small dataset
func createId() string {
	id := [12]byte(primitive.NewObjectID())
	return fmt.Sprintf("%x", md5.Sum(id[:]))[:8]
}

func normalizeImage(image string) string {
	// remove hash added by docker
	if i := strings.Index(image, "@sha256:"); i > 0 {
		image = image[:i]
	}
	return image
}

func indentJSON(raw []byte) (s string, err error) {
	buf := &bytes.Buffer{}
	err = json.Indent(buf, raw, "", "    ")
	if err == nil {
		s = buf.String()
	}
	return
}

func now() dao.Time {
	return dao.Time(time.Now())
}

func newOperator(user web.User) dao.Operator {
	return dao.Operator{ID: user.ID(), Name: user.Name()}
}

func init() {
	container.Put(NewNetwork)
	container.Put(NewNode)
	container.Put(NewRegistry)
	container.Put(NewService)
	container.Put(NewTask)
	container.Put(NewConfig)
	container.Put(NewSecret)
	container.Put(NewStack)
	container.Put(NewImage)
	container.Put(NewContainer)
	container.Put(NewVolume)
	container.Put(NewUser)
	container.Put(NewRole)
	container.Put(NewEvent)
	container.Put(NewSetting)
	container.Put(NewMetric)
	container.Put(NewChart)
	container.Put(NewSystem)
	container.Put(NewSession)
	container.Put(NewDashboard)
}
