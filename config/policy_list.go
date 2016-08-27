package config

import (
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apg-hcl/config/policy"
)

//type PolicyLoaderFunc func(*ast.ObjectItem) (interface{}, error)

//type HCLLoader interface {
//LoadHCL(item *ast.ObjectItem) (interface{}, error)
//}

var PolicyList = map[string]func(*ast.ObjectItem) (interface{}, error){
	"assign_message": policy.LoadAssignMessageHCL,
}
