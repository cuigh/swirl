package misc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/cuigh/auxo/util/i18n"
)

var Funcs = map[string]interface{}{
	//"time": func(t time.Time) string {
	//	return t.Local().Format("2006-01-02 15:04:05")
	//},
	"eq": func(v1, v2 interface{}) bool {
		return fmt.Sprint(v1) == fmt.Sprint(v2)
	},
	"elem": func(i interface{}) interface{} {
		v := reflect.ValueOf(i)
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		return v.Interface()
	},
	"trimZero": func(v interface{}) interface{} {
		s := fmt.Sprint(v)
		if s == "0" {
			return ""
		}
		return s
	},
	"trimPrefix": func(s, prefix string) string {
		return strings.TrimPrefix(s, prefix)
	},
}

func Message(lang string) func(key string, args ...interface{}) string {
	t, err := i18n.Find(lang, "en")
	if err != nil {
		panic(err)
	}
	if t == nil {
		panic("can't find language files")
	}

	return func(key string, args ...interface{}) string {
		if s := t.Format(key, args...); s != "" {
			return s
		}
		return "[" + key + "]"
	}
}

func FormatTime(offset int32) func(t time.Time) string {
	const layout = "2006-01-02 15:04:05"

	var loc *time.Location
	if offset == 0 {
		loc = time.Local
	} else {
		loc = time.FixedZone("", int(offset))
	}

	return func(t time.Time) string {
		return t.In(loc).Format(layout)
	}
}

func Page(count, pageIndex, pageSize int) (start, end int) {
	start = pageSize * (pageIndex - 1)
	end = pageSize * pageIndex
	if count < start {
		start, end = 0, 0
	} else if count < end {
		end = count
	}
	return
}

func JSONIndent(raw []byte) (s string, err error) {
	buf := &bytes.Buffer{}
	err = json.Indent(buf, raw, "", "    ")
	if err == nil {
		s = buf.String()
	}
	return
}
