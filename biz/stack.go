package biz

import (
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/dao"
	"github.com/cuigh/swirl/model"
)

// Stack return a stack biz instance.
var Archive = &archiveBiz{}

type archiveBiz struct {
}

func (b *archiveBiz) List(args *model.ArchiveListArgs) (archives []*model.Archive, count int, err error) {
	do(func(d dao.Interface) {
		archives, count, err = d.ArchiveList(args)
	})
	return
}

func (b *archiveBiz) Create(archive *model.Archive) (err error) {
	do(func(d dao.Interface) {
		err = d.ArchiveCreate(archive)
		//if err == nil {
		//	Event.CreateStackArchive(model.EventActionCreate, archive.ID, archive.Name, ctx.User())
		//}
	})
	return
}

func (b *archiveBiz) Delete(id string, user web.User) (err error) {
	do(func(d dao.Interface) {
		var archive *model.Archive
		archive, err = d.ArchiveGet(id)
		if err != nil {
			return
		}

		err = d.ArchiveDelete(id)
		if err == nil {
			Event.CreateStackArchive(model.EventActionDelete, id, archive.Name, user)
		}
	})
	return
}

func (b *archiveBiz) Get(id string) (archives *model.Archive, err error) {
	do(func(d dao.Interface) {
		archives, err = d.ArchiveGet(id)
	})
	return
}

func (b *archiveBiz) Update(archive *model.Archive) (err error) {
	do(func(d dao.Interface) {
		err = d.ArchiveUpdate(archive)
	})
	return
}
