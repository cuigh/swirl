package biz

import (
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/swirl/dao"
)

func do(fn func(d dao.Interface)) {
	d, err := dao.Get()
	if err != nil {
		panic(errors.Wrap("failed to load storage engine", err))
	}

	fn(d)
}
