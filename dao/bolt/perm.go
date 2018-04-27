package bolt

import (
	"github.com/cuigh/swirl/model"
)

func (d *Dao) PermGet(resType, resID string) (p *model.Perm, err error) {
	key := resType + "." + resID
	var v Value
	v, err = d.get("perm", key)
	if err == nil {
		if v != nil {
			p = &model.Perm{}
			err = v.Unmarshal(p)
		}
	}
	return
}

func (d *Dao) PermUpdate(perm *model.Perm) (err error) {
	key := perm.ResType + "." + perm.ResID
	return d.update("perm", key, perm)
}

func (d *Dao) PermDelete(resType, resID string) (err error) {
	key := resType + "." + resID
	return d.delete("perm", key)
}
