package config

import (
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/config/policy"
)

var PolicyList = map[string]func(*ast.ObjectItem) (interface{}, error){
	"assign_message": policy.LoadAssignMessageHCL,
	"quota":          policy.LoadQuotaHCL,
	"script":         policy.LoadScriptHCL,
	"javascript":     policy.LoadJavaScriptHCL,
	"verify_apiky":   policy.LoadVerifyAPIKeyHCL,
}
