package security

import (
	"strings"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/swirl/biz"
	"github.com/cuigh/swirl/model"
)

// Checker check permission of user
func Checker(user web.User, h web.HandlerInfo) bool {
	if au, ok := user.(*model.AuthUser); ok {
		return au.IsAllowed(h.Name())
	}
	return false
}

// Permiter is a middleware for validate data permission.
func Permiter(next web.HandlerFunc) web.HandlerFunc {
	return func(ctx web.Context) error {
		opt := ctx.Handler().Option("perm")
		if opt != "" {
			array := strings.Split(opt, ",")
			err := biz.Perm.Check(ctx.User(), array[0], array[1], ctx.P(array[2]))
			if err != nil {
				return err
			}
		}
		return next(ctx)
	}
}

// Perm holds permission key and description.
type Perm struct {
	Key  string
	Text string
}

// PermGroup holds information of a perm group.
type PermGroup struct {
	Name  string
	Perms []Perm
}

// Perms holds all valid perm groups.
var Perms = []PermGroup{
	{
		Name: "Registry",
		Perms: []Perm{
			{Key: "registry.list", Text: "View list"},
			{Key: "registry.create", Text: "Create"},
			{Key: "registry.delete", Text: "Delete"},
			{Key: "registry.update", Text: "Update"},
		},
	},
	{
		Name: "Node",
		Perms: []Perm{
			{Key: "node.list", Text: "View list"},
			{Key: "node.detail", Text: "View detail"},
			{Key: "node.raw", Text: "View raw"},
			{Key: "node.edit", Text: "View edit"},
			{Key: "node.update", Text: "Update"},
			{Key: "node.delete", Text: "Delete"},
		},
	},
	{
		Name: "Network",
		Perms: []Perm{
			{Key: "network.list", Text: "View list"},
			{Key: "network.new", Text: "View new"},
			{Key: "network.detail", Text: "View detail"},
			{Key: "network.raw", Text: "View raw"},
			{Key: "network.create", Text: "Create"},
			{Key: "network.delete", Text: "Delete"},
			{Key: "network.disconnect", Text: "Disconnect container"},
		},
	},
	{
		Name: "Service",
		Perms: []Perm{
			{Key: "service.list", Text: "View list"},
			{Key: "service.new", Text: "View new"},
			{Key: "service.detail", Text: "View detail"},
			{Key: "service.raw", Text: "View raw"},
			{Key: "service.logs", Text: "View logs"},
			{Key: "service.stats", Text: "View stats"},
			{Key: "service.edit", Text: "View edit"},
			{Key: "service.create", Text: "Create"},
			{Key: "service.delete", Text: "Delete"},
			{Key: "service.update", Text: "Update"},
			{Key: "service.scale", Text: "Scale"},
			{Key: "service.rollback", Text: "Rollback"},
		},
	},
	{
		Name: "Template",
		Perms: []Perm{
			{Key: "template.list", Text: "View list"},
			{Key: "template.new", Text: "View new"},
			{Key: "template.edit", Text: "View edit"},
			{Key: "template.create", Text: "Create"},
			{Key: "template.delete", Text: "Delete"},
			{Key: "template.update", Text: "Update"},
		},
	},
	{
		Name: "Stack",
		Perms: []Perm{
			{Key: "stack.task.list", Text: "View task list"},
			{Key: "stack.task.delete", Text: "Delete task"},
			{Key: "stack.archive.list", Text: "View archive list"},
			{Key: "stack.archive.new", Text: "View archive new"},
			{Key: "stack.archive.detail", Text: "View archive detail"},
			{Key: "stack.archive.edit", Text: "View archive edit"},
			{Key: "stack.archive.delete", Text: "Delete archive"},
			{Key: "stack.archive.create", Text: "Create archive"},
			{Key: "stack.archive.update", Text: "Update archive"},
			{Key: "stack.archive.deploy", Text: "Deploy archive"},
		},
	},
	{
		Name: "Task",
		Perms: []Perm{
			{Key: "task.list", Text: "View list"},
			{Key: "task.detail", Text: "View detail"},
			{Key: "task.raw", Text: "View raw"},
			{Key: "task.logs", Text: "View logs"},
		},
	},
	{
		Name: "Image",
		Perms: []Perm{
			{Key: "image.list", Text: "View list"},
			{Key: "image.detail", Text: "View detail"},
			{Key: "image.raw", Text: "View raw"},
			{Key: "image.delete", Text: "Delete"},
		},
	},
	{
		Name: "Container",
		Perms: []Perm{
			{Key: "container.list", Text: "View list"},
			{Key: "container.detail", Text: "View detail"},
			{Key: "container.raw", Text: "View raw"},
			{Key: "container.logs", Text: "View logs"},
			{Key: "container.delete", Text: "Delete"},
		},
	},
	{
		Name: "Volume",
		Perms: []Perm{
			{Key: "volume.list", Text: "View list"},
			{Key: "volume.new", Text: "View new"},
			{Key: "volume.detail", Text: "View detail"},
			{Key: "volume.raw", Text: "View raw"},
			{Key: "volume.create", Text: "Create"},
			{Key: "volume.delete", Text: "Delete"},
			{Key: "volume.prune", Text: "Prune"},
		},
	},
	{
		Name: "Secret",
		Perms: []Perm{
			{Key: "secret.list", Text: "View list"},
			{Key: "secret.new", Text: "View new"},
			{Key: "secret.edit", Text: "View edit"},
			{Key: "secret.create", Text: "Create"},
			{Key: "secret.delete", Text: "Delete"},
			{Key: "secret.update", Text: "Update"},
		},
	},
	{
		Name: "Config",
		Perms: []Perm{
			{Key: "config.list", Text: "View list"},
			{Key: "config.new", Text: "View new"},
			{Key: "config.edit", Text: "View edit"},
			{Key: "config.create", Text: "Create"},
			{Key: "config.delete", Text: "Delete"},
			{Key: "config.update", Text: "Update"},
		},
	},
	{
		Name: "Role",
		Perms: []Perm{
			{Key: "role.list", Text: "View list"},
			{Key: "role.new", Text: "View new"},
			{Key: "role.detail", Text: "View detail"},
			{Key: "role.edit", Text: "View edit"},
			{Key: "role.create", Text: "Create"},
			{Key: "role.delete", Text: "Delete"},
			{Key: "role.update", Text: "Update"},
		},
	},
	{
		Name: "User",
		Perms: []Perm{
			{Key: "user.list", Text: "View list"},
			{Key: "user.new", Text: "View new"},
			{Key: "user.detail", Text: "View detail"},
			{Key: "user.edit", Text: "View edit"},
			{Key: "user.create", Text: "Create"},
			{Key: "user.delete", Text: "Delete"},
			{Key: "user.update", Text: "Update"},
			{Key: "user.block", Text: "Block"},
			{Key: "user.unblock", Text: "Unblock"},
		},
	},
	{
		Name: "Setting",
		Perms: []Perm{
			{Key: "setting.edit", Text: "View edit"},
			{Key: "setting.update", Text: "Update"},
		},
	},
	{
		Name: "Event",
		Perms: []Perm{
			{Key: "event.list", Text: "View list"},
		},
	},
}
