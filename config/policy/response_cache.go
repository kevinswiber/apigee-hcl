package policy

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type ResponseCachePolicy struct {
	XMLName                     string `xml:"ResponseCache" hcl:"-"`
	Policy                      `hcl:",squash"`
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
	ExpiryDate   *expiryDate         `xml:"" hcl:"expiry_date"`
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

func LoadResponseCacheHCL(item *ast.ObjectItem) (interface{}, error) {
	var p ResponseCachePolicy

	if err := LoadCommonPolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	// var listVal *ast.ObjectList
	// if ot, ok := item.Val.(*ast.ObjectType); ok {
	// 	listVal = ot.List
	// } else {
	// 	return nil, fmt.Errorf("response cache policy not an object")
	// }

	// if cacheKeyList := listVal.Filter("cache_key"); len(cacheKeyList.Items) > 0 {
	// 	a, err := loadResponseCacheCacheKeyHCL(cacheKeyList.Items)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	p.CacheKey = a
	// }
	//
	// if expirySettingsList := listVal.Filter("expiry_settings"); len(expirySettingsList.Items) > 0 {
	// 	a, err := loadResponseCacheExpirySettingsHCL(expirySettingsList.Items)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	p.ExpirySettings = a
	// }

	return p, nil
}

// func loadResponseCacheCacheKeyHCL(item *ast.ObjectItem) (*cacheKey, error) {
// 	result := new(cacheKey)
//
// 	var listVal *ast.ObjectList
// 	if ot, ok := item.Val.(*ast.ObjectType); ok {
// 		listVal = ot.List
// 	} else {
// 		return nil, fmt.Errorf("cacheKey not an object")
// 	}
//
// 	if cacheKeyPrefixList := listVal.Filter("prefix"); len(cacheKeyFilterList.Items) > 0 {
// 		a, err := loadResponseCacheKeyPrefixHCL(cacheKeyPrefixList.Items)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		result.Prefix = a
// 	}
//
// 	if cacheKeyFragmentList := listVal.Filter("key_fragment"); len(cacheKeyFragmentList.Items) > 0 {
// 		a, err := loadResponseCacheKeyFragmentHCL(cacheKeyFragmentList.Items)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		result.KeyFragment = a
// 	}
//
// 	return result, nil
// }
//
// func loadResponseCacheExpirySettingsHCL(item *ast.ObjectItem) (*expirySettings, error) {
// 	result := new(expirySettings)
//
// 	var listVal *ast.ObjectList
// 	if ot, ok := item.Val.(*ast.ObjectType); ok {
// 		listVal = ot.List
// 	} else {
// 		return nil, fmt.Errorf("expirySettings not an object")
// 	}
//
// }
