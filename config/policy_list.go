package config

import (
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/config/policy"
)

var PolicyList = map[string]func(*ast.ObjectItem) (interface{}, error){
	"assign_message":    policy.LoadAssignMessageHCL,
	"extract_variables": policy.LoadExtractVariablesHCL,
	"javascript":        policy.LoadJavaScriptHCL,
	"quota":             policy.LoadQuotaHCL,
	"raise_fault":       policy.LoadRaiseFaultHCL,
	"response_cache":    policy.LoadResponseCacheHCL,
	"script":            policy.LoadScriptHCL,
	"spike_arrest":      policy.LoadSpikeArrestHCL,
	"verify_api_key":    policy.LoadVerifyAPIKeyHCL,
}
