package misc

import (
	"fmt"
	"reflect"
)

var Funcs = map[string]interface{}{
	"limit": func(s string, length int) string {
		if len(s) > length {
			return s[:length] + "..."
		}
		return s
	},
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
	"slice": func(values ...interface{}) interface{} {
		return values
	},
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
