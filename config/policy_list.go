package config

import (
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/config/policy"
)

// PolicyList is a map of HCL policy types to HCL loader functions.
var PolicyList = map[string]func(*ast.ObjectItem) (interface{}, error){
	"assign_message":       policy.LoadAssignMessageHCL,
	"extract_variables":    policy.LoadExtractVariablesHCL,
	"javascript":           policy.LoadJavaScriptHCL,
	"quota":                policy.LoadQuotaHCL,
	"raise_fault":          policy.LoadRaiseFaultHCL,
	"response_cache":       policy.LoadResponseCacheHCL,
	"script":               policy.LoadScriptHCL,
	"spike_arrest":         policy.LoadSpikeArrestHCL,
	"statistics_collector": policy.LoadStatisticsCollectorHCL,
	"verify_api_key":       policy.LoadVerifyAPIKeyHCL,
}
