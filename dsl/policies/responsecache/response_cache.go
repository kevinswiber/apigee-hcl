package responsecache

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/policy"
)

// ResponseCache represents a <ResponseCache/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/response-cache-policy
type ResponseCache struct {
	XMLName                     string `xml:"ResponseCache" hcl:"-"`
	policy.Policy               `hcl:",squash"`
	Type                        string          `xml:"type,attr,omitempty" hcl:"type"`
	DisplayName                 string          `xml:",omitempty" hcl:"display_name"`
	CacheKey                    *cacheKey       `xml:"CacheKey" hcl:"cache_key"`
	Scope                       string          `xml:",omitempty" hcl:"scope"`
	ExpirySettings              *expirySettings `xml:"ExpirySettings" hcl:"expiry_settings"`
	CacheResource               string          `xml:",omitempty" hcl:"cache_resource"`
	CacheLookupTimeoutInSeconds int             `xml:",omitempty" hcl:"lookup_timeout"`
	ExcludeErrorResponse        bool            `xml:",omitempty" hcl:"exclude_error_response"`
	SkipCacheLookup             string          `xml:",omitempty" hcl:"skip_cache_lookup"`
	SkipCachePopulation         string          `xml:",omitempty" hcl:"skip_cache_population"`
	UseAcceptHeader             bool            `xml:",omitempty" hcl:"use_accept_header"`
	UseResponseCacheHeaders     bool            `xml:",omitempty" hcl:"use_response_cache_headers"`
}

type cacheKey struct {
	XMLName     string              `xml:"CacheKey" hcl:"-"`
	Prefix      string              `xml:",omitempty" hcl:"prefix"`
	KeyFragment []*cacheKeyFragment `xml:",omitempty" hcl:"key_fragment"`
}

type cacheKeyFragment struct {
	XMLName string `xml:"KeyFragment" hcl:"-"`
	Ref     string `xml:"ref,attr,omitempty" hcl:"ref"`
	Value   string `xml:",chardata" hcl:"value"`
}

type expirySettings struct {
	XMLName      string              `xml:"ExpirySettings" hcl:"-"`
	TimeOfDay    *expiryTimeOfDay    `xml:",omitempty" hcl:"time_of_day"`
	TimeoutInSec *expiryTimeoutInSec `xml:",omitempty" hcl:"timeout_in_sec"`
	ExpiryDate   *expiryDate         `xml:",omitempty" hcl:"expiry_date"`
}

type expiryTimeOfDay struct {
	XMLName string `xml:"TimeOfDay" hcl:"-"`
	Ref     string `xml:"ref,attr,omitempty" hcl:"ref"`
	Value   string `xml:",chardata" hcl:"value"`
}

type expiryTimeoutInSec struct {
	XMLName string `xml:"TimeoutInSec" hcl:"-"`
	Ref     string `xml:"ref,attr,omitempty" hcl:"ref"`
	Value   string `xml:",chardata" hcl:"value"`
}

type expiryDate struct {
	XMLName string `xml:"ExpiryDate" hcl:"-"`
	Ref     string `xml:"ref,attr,omitempty" hcl:"ref"`
	Value   string `xml:",chardata" hcl:"value"`
}

// DecodeHCL converts an HCL ast.ObjectItem into a ResponseCache
func DecodeHCL(item *ast.ObjectItem) (interface{}, error) {
	var p ResponseCache

	if err := policy.DecodeHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	return &p, nil
}
