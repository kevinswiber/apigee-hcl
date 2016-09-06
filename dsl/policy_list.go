package dsl

import (
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/assignmessage"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/extractvariables"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/javascript"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/quota"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/raisefault"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/responsecache"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/script"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/servicecallout"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/spikearrest"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/statisticscollector"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/verifyapikey"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/xmltojson"
)

// PolicyList is a map of HCL policy types to policy factory functions.
var PolicyList = map[string]func(*ast.ObjectItem) (interface{}, error){
	"assign_message":       assignmessage.DecodeHCL,
	"extract_variables":    extractvariables.DecodeHCL,
	"javascript":           javascript.DecodeHCL,
	"quota":                quota.DecodeHCL,
	"raise_fault":          raisefault.DecodeHCL,
	"response_cache":       responsecache.DecodeHCL,
	"script":               script.DecodeHCL,
	"service_callout":      servicecallout.DecodeHCL,
	"spike_arrest":         spikearrest.DecodeHCL,
	"statistics_collector": statisticscollector.DecodeHCL,
	"verify_api_key":       verifyapikey.DecodeHCL,
	"xml_to_json":          xmltojson.DecodeHCL,
}
