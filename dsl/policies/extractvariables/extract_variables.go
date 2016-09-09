package extractvariables

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	hclEncoding "github.com/kevinswiber/apigee-hcl/dsl/encoding/hcl"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/policy"
)

// ExtractVariables represents an <ExtractVariables/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/extract-variables-policy
type ExtractVariables struct {
	XMLName                   string `xml:"ExtractVariables" hcl:"-"`
	policy.Policy             `hcl:",squash"`
	DisplayName               string          `xml:",omitempty" hcl:"display_name"`
	Source                    *evSource       `xml:",omitempty" hcl:"source"`
	VariablePrefix            string          `xml:",omitempty" hcl:"variable_prefix"`
	IgnoreUnresolvedVariables bool            `xml:",omitempty" hcl:"ignore_unresolved_variables"`
	URIPaths                  []*evURIPath    `xml:"URIPath,omitempty" hcl:"uri_path"`
	QueryParams               []*evQueryParam `xml:"QueryParam,omitempty" hcl:"query_param"`
	Headers                   []*evHeader     `xml:"Header,omitempty" hcl:"header"`
	FormParams                []*evFormParam  `xml:"FormParam,omitempty" hcl:"form_param"`
	Variables                 []*evVariable   `xml:"Variable,omitempty" hcl:"variable"`
	JSONPayload               *evJSONPayload  `xml:",omitempty" hcl:"json_payload"`
	XMLPayload                *evXMLPayload   `xml:",omitempty" hcl:"xml_payload"`
}

type evSource struct {
	XMLName      string `xml:"Source" hcl:"-"`
	ClearPayload bool   `xml:"clearPayload,attr,omitempty" hcl:"clear_payload"`
	Value        string `xml:",chardata" hcl:"value"`
}

type evURIPath struct {
	XMLName  string       `xml:"URIPath" hcl:"-"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evQueryParam struct {
	XMLName  string       `xml:"QueryParam" hcl:"-"`
	Name     string       `xml:"name,attr" hcl:"-"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evHeader struct {
	XMLName  string       `xml:"Header" hcl:"-"`
	Name     string       `xml:"name,attr" hcl:"-"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evFormParam struct {
	XMLName  string       `xml:"FormParam" hcl:"-"`
	Name     string       `xml:"name,attr" hcl:"-"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evVariable struct {
	XMLName  string       `xml:"Variable" hcl:"-"`
	Name     string       `xml:"name,attr" hcl:"-"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evJSONPayload struct {
	XMLName   string                   `xml:"JSONPayload" hcl:"-"`
	Variables []*evJSONPayloadVariable `xml:"Variable" hcl:"variable"`
}

type evJSONPayloadVariable struct {
	XMLName  string `xml:"Variable" hcl:"-"`
	Name     string `xml:"name,attr" hcl:"-"`
	Type     string `xml:"type,attr,omitempty" hcl:"type"`
	JSONPath string `hcl:"json_path"`
}

type evXMLPayload struct {
	XMLName               string                   `xml:"XMLPayload" hcl:"-"`
	StopPayloadProcessing bool                     `xml:"stopPayloadProcessing,attr,omitempty" hcl:"stop_payload_processing"`
	Namespaces            []*evXMLPayloadNamespace `xml:"Namespaces>Namespace,omitempty" hcl:"namespace"`
	Variables             []*evXMLPayloadVariable  `xml:"Variable" hcl:"variable"`
}

type evXMLPayloadNamespace struct {
	Prefix string `xml:"prefix,attr,omitempty" hcl:"-"`
	Value  string `xml:",chardata" hcl:"value"`
}

type evXMLPayloadVariable struct {
	XMLName string `xml:"Variable" hcl:"-"`
	Name    string `xml:"name,attr" hcl:"-"`
	Type    string `xml:"type,attr,omitempty" hcl:"type"`
	XPath   string `hcl:"xpath"`
}

type evPattern struct {
	XMLName    string `xml:"Pattern" hcl:"-"`
	IgnoreCase bool   `xml:"ignoreCase,attr,omitempty" hcl:"ignore_case"`
	Value      string `xml:",chardata" hcl:"value"`
}

// DecodeHCL converts an HCL ast.ObjectItem into an ExtractVariables object.
func DecodeHCL(item *ast.ObjectItem) (interface{}, error) {
	var errors *multierror.Error
	var p ExtractVariables

	if err := policy.DecodeHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		pos := item.Val.Pos()
		newError := hclEncoding.PosError{
			Pos: pos,
			Err: fmt.Errorf("extract variables policy not an object"),
		}
		return nil, &newError
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	if uriPathList := listVal.Filter("uri_path"); len(uriPathList.Items) > 0 {
		uriPaths, err := decodeURIPathsHCL(uriPathList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.URIPaths = uriPaths
		}
	}

	if queryParamList := listVal.Filter("query_param"); len(queryParamList.Items) > 0 {
		queryParams, err := decodeQueryParamsHCL(queryParamList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.QueryParams = queryParams
		}
	}

	if headerList := listVal.Filter("header"); len(headerList.Items) > 0 {
		headers, err := decodeHeadersHCL(headerList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.Headers = headers
		}
	}

	if formParamList := listVal.Filter("form_param"); len(formParamList.Items) > 0 {
		formParams, err := decodeFormParamsHCL(formParamList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.FormParams = formParams
		}
	}

	if variableList := listVal.Filter("variable"); len(variableList.Items) > 0 {
		variables, err := decodeVariablesHCL(variableList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.Variables = variables
		}
	}

	if jsonPayloadList := listVal.Filter("json_payload"); len(jsonPayloadList.Items) > 0 {
		jsonPayload, err := decodeJSONPayloadHCL(jsonPayloadList.Items[0])
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.JSONPayload = jsonPayload
		}
	}

	if xmlPayloadList := listVal.Filter("xml_payload"); len(xmlPayloadList.Items) > 0 {
		xmlPayload, err := decodeXMLPayloadHCL(xmlPayloadList.Items[0])
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.XMLPayload = xmlPayload
		}
	}

	if len(p.URIPaths) == 0 && len(p.QueryParams) == 0 &&
		len(p.FormParams) == 0 && len(p.Variables) == 0 &&
		p.JSONPayload == nil && p.XMLPayload == nil {
		pos := item.Val.Pos()
		newError := hclEncoding.PosError{
			Pos: pos,
			Err: fmt.Errorf("extract variables requires one of uri_path, query_param, " +
				"header, form_param, json_payload, or xml_payload"),
		}
		errors = multierror.Append(errors, &newError)

	}

	if errors != nil {
		return nil, errors
	}

	return &p, nil
}

func decodeURIPathsHCL(items []*ast.ObjectItem) ([]*evURIPath, error) {
	var uriPaths []*evURIPath
	for _, item := range items {
		var up evURIPath

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("uri_path not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&up, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if patternList := listVal.Filter("pattern"); len(patternList.Items) > 0 {
			patterns, err := decodePatternsHCL(patternList.Items)
			if err != nil {
				return nil, err
			}

			up.Patterns = patterns
		}

		uriPaths = append(uriPaths, &up)
	}
	return uriPaths, nil
}

func decodeQueryParamsHCL(items []*ast.ObjectItem) ([]*evQueryParam, error) {
	var queryParams []*evQueryParam
	for _, item := range items {
		var qp evQueryParam

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("query_param not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&qp, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("query_param requires a name"),
			}
			return nil, &newError
		}

		qp.Name = item.Keys[0].Token.Value().(string)

		if patternList := listVal.Filter("pattern"); len(patternList.Items) > 0 {
			patterns, err := decodePatternsHCL(patternList.Items)
			if err != nil {
				return nil, err
			}

			qp.Patterns = patterns
		}

		queryParams = append(queryParams, &qp)
	}
	return queryParams, nil
}

func decodeHeadersHCL(items []*ast.ObjectItem) ([]*evHeader, error) {
	var headers []*evHeader
	for _, item := range items {
		var hdr evHeader

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("header not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&hdr, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("header requires a name"),
			}
			return nil, &newError
		}

		hdr.Name = item.Keys[0].Token.Value().(string)

		if patternList := listVal.Filter("pattern"); len(patternList.Items) > 0 {
			patterns, err := decodePatternsHCL(patternList.Items)
			if err != nil {
				return nil, err
			}

			hdr.Patterns = patterns
		}

		headers = append(headers, &hdr)
	}
	return headers, nil
}

func decodeFormParamsHCL(items []*ast.ObjectItem) ([]*evFormParam, error) {
	var formParams []*evFormParam
	for _, item := range items {
		var fp evFormParam

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("form_param not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&fp, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("form_param requires a name"),
			}
			return nil, &newError
		}

		fp.Name = item.Keys[0].Token.Value().(string)

		if patternList := listVal.Filter("pattern"); len(patternList.Items) > 0 {
			patterns, err := decodePatternsHCL(patternList.Items)
			if err != nil {
				return nil, err
			}

			fp.Patterns = patterns
		}

		formParams = append(formParams, &fp)
	}
	return formParams, nil
}

func decodeVariablesHCL(items []*ast.ObjectItem) ([]*evVariable, error) {
	var variables []*evVariable
	for _, item := range items {
		var v evVariable

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("variable not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&v, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("variable requires a name"),
			}
			return nil, &newError
		}

		v.Name = item.Keys[0].Token.Value().(string)

		if patternList := listVal.Filter("pattern"); len(patternList.Items) > 0 {
			patterns, err := decodePatternsHCL(patternList.Items)
			if err != nil {
				return nil, err
			}

			v.Patterns = patterns
		}

		variables = append(variables, &v)
	}
	return variables, nil
}

func decodeJSONPayloadHCL(item *ast.ObjectItem) (*evJSONPayload, error) {
	var p evJSONPayload

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		pos := item.Val.Pos()
		newError := hclEncoding.PosError{
			Pos: pos,
			Err: fmt.Errorf("json_payload not an object"),
		}
		return nil, &newError
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	if variableList := listVal.Filter("variable"); len(variableList.Items) > 0 {
		variables, err := decodeJSONPayloadVariablesHCL(variableList.Items)
		if err != nil {
			return nil, err
		}

		p.Variables = variables
	}

	return &p, nil
}

func decodeJSONPayloadVariablesHCL(items []*ast.ObjectItem) ([]*evJSONPayloadVariable, error) {
	var variables []*evJSONPayloadVariable
	for _, item := range items {
		var v evJSONPayloadVariable

		if _, ok := item.Val.(*ast.ObjectType); !ok {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("variable not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&v, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("variable requires a name"),
			}
			return nil, &newError
		}

		v.Name = item.Keys[0].Token.Value().(string)

		variables = append(variables, &v)
	}

	return variables, nil
}

func decodeXMLPayloadHCL(item *ast.ObjectItem) (*evXMLPayload, error) {
	var p evXMLPayload

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		pos := item.Val.Pos()
		newError := hclEncoding.PosError{
			Pos: pos,
			Err: fmt.Errorf("xml_payload not an object"),
		}
		return nil, &newError
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	if namespaceList := listVal.Filter("namespace"); len(namespaceList.Items) > 0 {
		namespaces, err := decodeXMLPayloadNamespacesHCL(namespaceList.Items)
		if err != nil {
			return nil, err
		}

		p.Namespaces = namespaces
	}

	if variableList := listVal.Filter("variable"); len(variableList.Items) > 0 {
		variables, err := decodeXMLPayloadVariablesHCL(variableList.Items)
		if err != nil {
			return nil, err
		}

		p.Variables = variables
	}

	return &p, nil
}

func decodeXMLPayloadVariablesHCL(items []*ast.ObjectItem) ([]*evXMLPayloadVariable, error) {
	var variables []*evXMLPayloadVariable
	for _, item := range items {
		var v evXMLPayloadVariable

		if _, ok := item.Val.(*ast.ObjectType); !ok {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("variable not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&v, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("variable requires a name"),
			}
			return nil, &newError
		}

		v.Name = item.Keys[0].Token.Value().(string)

		variables = append(variables, &v)
	}

	return variables, nil
}

func decodeXMLPayloadNamespacesHCL(items []*ast.ObjectItem) ([]*evXMLPayloadNamespace, error) {
	var namespaces []*evXMLPayloadNamespace
	for _, item := range items {
		var v evXMLPayloadNamespace

		if _, ok := item.Val.(*ast.ObjectType); !ok {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("namespace not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&v, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("namespace requires a name"),
			}
			return nil, &newError
		}

		v.Prefix = item.Keys[0].Token.Value().(string)

		namespaces = append(namespaces, &v)
	}

	return namespaces, nil
}

func decodePatternsHCL(items []*ast.ObjectItem) ([]*evPattern, error) {
	var patterns []*evPattern
	for _, item := range items {
		var pat evPattern

		if _, ok := item.Val.(*ast.ObjectType); !ok {
			pos := item.Val.Pos()
			newError := hclEncoding.PosError{
				Pos: pos,
				Err: fmt.Errorf("pattern not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&pat, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		patterns = append(patterns, &pat)
	}

	return patterns, nil
}
