package assignmessage

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// Add includes elements to add to an HTTP message
type Add struct {
	XMLName     string         `xml:"Add" hcl:"-"`
	Headers     *[]*Header     `xml:"Headers>Header" hcl:"header"`
	QueryParams *[]*QueryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams  *[]*FormParam  `xml:"FormParams>FormParam" hcl:"form_param"`
}

// Copy includes elements to copy to an HTTP message
type Copy struct {
	XMLName      string         `xml:"Copy" hcl:"-"`
	Source       string         `xml:"source,attr,omitempty" hcl:"-"`
	Headers      *[]*Header     `xml:"Headers>Header" hcl:"header"`
	QueryParams  *[]*QueryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams   *[]*FormParam  `xml:"FormParams>FormParam" hcl:"form_param"`
	Payload      bool           `xml:",omitempty" hcl:"payload"`
	Version      bool           `xml:",omitempty" hcl:"version"`
	Verb         bool           `xml:",omitempty" hcl:"verb"`
	Path         bool           `xml:",omitempty" hcl:"path"`
	StatusCode   bool           `xml:",omitempty" hcl:"status_code"`
	ReasonPhrase bool           `xml:",omitempty" hcl:"reason_phrase"`
}

// Remove includes elements to remove from an HTTP message
type Remove struct {
	XMLName     string         `xml:"Remove" hcl:"-"`
	Headers     *[]*Header     `xml:"Headers>Header" hcl:"header"`
	QueryParams *[]*QueryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams  *[]*FormParam  `xml:"FormParams>FormParam" hcl:"form_param"`
	Payload     bool           `xml:",omitempty" hcl:"payload"`
}

// Set includes elements to set on an HTTP message
type Set struct {
	XMLName      string         `xml:"Set" hcl:"-"`
	Headers      *[]*Header     `xml:"Headers>Header" hcl:"header"`
	QueryParams  *[]*QueryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams   *[]*FormParam  `xml:"FormParams>FormParam" hcl:"form_param"`
	Payload      *Payload       `xml:",omitempty" hcl:"payload"`
	Version      string         `xml:",omitempty" hcl:"version"`
	Verb         string         `xml:",omitempty" hcl:"verb"`
	Path         string         `xml:",omitempty" hcl:"path"`
	StatusCode   int            `xml:",omitempty" hcl:"status_code"`
	ReasonPhrase string         `xml:",omitempty" hcl:"reason_phrase"`
}

// Payload includes content to include an HTTP message body.
type Payload struct {
	XMLName        string `xml:"Payload" hcl:"-"`
	ContentType    string `xml:"contentType,attr,omitempty" hcl:"content_type"`
	VariablePrefix string `xml:"variablePrefix,attr,omitempty" hcl:"variable_prefix"`
	VariableSuffix string `xml:"variableSuffix,attr,omitempty" hcl:"variable_suffix"`
	Value          string `xml:",chardata" hcl:"value"`
}

// Header represents an HTTP header
type Header struct {
	XMLName string `xml:"Header" hcl:"-"`
	Name    string `xml:"name,attr" hcl:"-"`
	Value   string `xml:",chardata" hcl:"value"`
}

// QueryParam represents a URL query parameter.
type QueryParam struct {
	XMLName string `xml:"QueryParam" hcl:"-"`
	Name    string `xml:"name,attr" hcl:"-"`
	Value   string `xml:",chardata" hcl:"value"`
}

// FormParam represents an HTTP message form parameter.
type FormParam struct {
	XMLName string `xml:"FormParam" hcl:"-"`
	Name    string `xml:"name,attr" hcl:"-"`
	Value   string `xml:",chardata" hcl:"value"`
}

// DecodeAddHCL converts HCL into an Add struct.
func DecodeAddHCL(item *ast.ObjectItem) (*Add, error) {
	result := new(Add)

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("add not an object")
	}

	headers, err := DecodeHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	qparams, err := DecodeQueryParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.QueryParams = qparams

	fparams, err := DecodeFormParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.FormParams = fparams
	return result, nil
}

// DecodeCopyHCL converts HCL into a Copy struct.
func DecodeCopyHCL(item *ast.ObjectItem) (*Copy, error) {
	var result Copy
	if err := hcl.DecodeObject(&result, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("copy not an object")
	}

	headers, err := DecodeHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	qparams, err := DecodeQueryParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.QueryParams = qparams

	fparams, err := DecodeFormParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.FormParams = fparams
	return &result, nil
}

// DecodeRemoveHCL converts HCL into a Remove struct
func DecodeRemoveHCL(item *ast.ObjectItem) (*Remove, error) {
	var result Remove
	if err := hcl.DecodeObject(&result, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("remove not an object")
	}

	headers, err := DecodeHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	qparams, err := DecodeQueryParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.QueryParams = qparams

	fparams, err := DecodeFormParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.FormParams = fparams
	return &result, nil
}

// DecodeSetHCL converts HCL into a Set struct.
func DecodeSetHCL(item *ast.ObjectItem) (*Set, error) {
	var result Set
	if err := hcl.DecodeObject(&result, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("set not an object")
	}

	headers, err := DecodeHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	qparams, err := DecodeQueryParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.QueryParams = qparams

	fparams, err := DecodeFormParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.FormParams = fparams
	return &result, nil
}

// DecodeHeadersHCL converts HCL headers into a Header array.
func DecodeHeadersHCL(listVal *ast.ObjectList) (*[]*Header, error) {
	var headers []*Header
	if headerList := listVal.Filter("header"); len(headerList.Items) > 0 {
		for _, h := range headerList.Items {
			var hdr Header
			if err := hcl.DecodeObject(&hdr, h.Val.(*ast.ObjectType)); err != nil {
				return nil, err
			}
			hdr.Name = h.Keys[0].Token.Value().(string)

			headers = append(headers, &hdr)
		}
	} else {
		return nil, nil
	}

	return &headers, nil
}

// DecodeQueryParamsHCL converts HCL into a QueryParam array.
func DecodeQueryParamsHCL(listVal *ast.ObjectList) (*[]*QueryParam, error) {
	var qparams []*QueryParam
	if qparamList := listVal.Filter("query_param"); len(qparamList.Items) > 0 {
		for _, q := range qparamList.Items {
			var qparam QueryParam
			if err := hcl.DecodeObject(&qparam, q.Val.(*ast.ObjectType)); err != nil {
				return nil, err
			}
			qparam.Name = q.Keys[0].Token.Value().(string)

			qparams = append(qparams, &qparam)
		}
	} else {
		return nil, nil
	}

	return &qparams, nil
}

// DecodeFormParamsHCL converts HCL into a FormParams array.
func DecodeFormParamsHCL(listVal *ast.ObjectList) (*[]*FormParam, error) {
	var fparams []*FormParam
	if fparamList := listVal.Filter("form_param"); len(fparamList.Items) > 0 {
		for _, f := range fparamList.Items {
			var fparam FormParam
			if err := hcl.DecodeObject(&fparam, f.Val.(*ast.ObjectType)); err != nil {
				return nil, err
			}
			fparam.Name = f.Keys[0].Token.Value().(string)

			fparams = append(fparams, &fparam)
		}
	} else {
		return nil, nil
	}

	return &fparams, nil
}
