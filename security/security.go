package security

import (
	"github.com/cuigh/auxo/app/container"
)

const PkgName = "security"

func init() {
	container.Put(NewIdentifier, container.Name("identifier"))
	container.Put(NewAuthorizer, container.Name("authorizer"))
}
