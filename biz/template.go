package biz

import (
	"encoding/json"
	"time"

	"github.com/cuigh/auxo/data/guid"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

// Template return a service template biz instance.
var Template = &templateBiz{}

type templateBiz struct {
}

func (b *templateBiz) List(args *model.TemplateListArgs) (tpls []*model.Template, count int, err error) {
	do(func(d dao.Interface) {
		tpls, count, err = d.TemplateList(args)
	})
	return
}

func (b *templateBiz) Create(tpl *model.Template, user web.User) (err error) {
	do(func(d dao.Interface) {
		tpl.ID = guid.New()
		err = d.TemplateCreate(tpl)
		if err == nil {
			Event.CreateServiceTemplate(model.EventActionCreate, tpl.ID, tpl.Name, user)
		}
	})
	return
}

func (b *templateBiz) Get(id string) (tpl *model.Template, err error) {
	do(func(d dao.Interface) {
		tpl, err = d.TemplateGet(id)
	})
	return
}

func (b *templateBiz) FillInfo(id string, si *model.ServiceInfo) (err error) {
	do(func(d dao.Interface) {
		var (
			tpl      *model.Template
			registry *model.Registry
		)

		tpl, err = d.TemplateGet(id)
		if err != nil || tpl == nil {
			return
		}

		err = json.Unmarshal([]byte(tpl.Content), si)
		if err != nil {
			return
		}

		if si.Registry != "" {
			registry, err = Registry.Get(si.Registry)
			if err != nil {
				return
			}
			si.RegistryURL = registry.URL
		}
	})
	return
}

func (b *templateBiz) Delete(id string, user web.User) (err error) {
	do(func(d dao.Interface) {
		var tpl *model.Template
		tpl, err = d.TemplateGet(id)
		if err != nil {
			return
		}

		err = d.TemplateDelete(id)
		if err == nil {
			Event.CreateServiceTemplate(model.EventActionDelete, id, tpl.Name, user)
		}
	})
	return
}

func (b *templateBiz) Update(tpl *model.Template, user web.User) (err error) {
	do(func(d dao.Interface) {
		tpl.UpdatedBy = user.ID()
		tpl.UpdatedAt = time.Now()
		err = d.TemplateUpdate(tpl)
		if err == nil {
			Event.CreateServiceTemplate(model.EventActionUpdate, tpl.ID, tpl.Name, user)
		}
	})
	return
}
