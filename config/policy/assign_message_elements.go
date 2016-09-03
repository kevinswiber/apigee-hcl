package policy

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type add struct {
	XMLName     string         `xml:"Add" hcl:"-"`
	Headers     *[]*header     `xml:"Headers>Header" hcl:"header"`
	QueryParams *[]*queryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams  *[]*formParam  `xml:"FormParams>FormParam" hcl:"form_param"`
}

type copy struct {
	XMLName      string         `xml:"Copy" hcl:"-"`
	Source       string         `xml:"source,attr,omitempty" hcl:"-"`
	Headers      *[]*header     `xml:"Headers>Header" hcl:"header"`
	QueryParams  *[]*queryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams   *[]*formParam  `xml:"FormParams>FormParam" hcl:"form_param"`
	Payload      bool           `xml:",omitempty" hcl:"payload"`
	Version      bool           `xml:",omitempty" hcl:"version"`
	Verb         bool           `xml:",omitempty" hcl:"verb"`
	Path         bool           `xml:",omitempty" hcl:"path"`
	StatusCode   bool           `xml:",omitempty" hcl:"status_code"`
	ReasonPhrase bool           `xml:",omitempty" hcl:"reason_phrase"`
}

type remove struct {
	XMLName     string         `xml:"Remove" hcl:"-"`
	Headers     *[]*header     `xml:"Headers>Header" hcl:"header"`
	QueryParams *[]*queryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams  *[]*formParam  `xml:"FormParams>FormParam" hcl:"form_param"`
	Payload     bool           `xml:",omitempty" hcl:"payload"`
}

type set struct {
	XMLName      string         `xml:"Set" hcl:"-"`
	Headers      *[]*header     `xml:"Headers>Header" hcl:"header"`
	QueryParams  *[]*queryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams   *[]*formParam  `xml:"FormParams>FormParam" hcl:"form_param"`
	Payload      *payload       `xml:",omitempty" hcl:"payload"`
	Version      string         `xml:",omitempty" hcl:"version"`
	Verb         string         `xml:",omitempty" hcl:"verb"`
	Path         string         `xml:",omitempty" hcl:"path"`
	StatusCode   int            `xml:",omitempty" hcl:"status_code"`
	ReasonPhrase string         `xml:",omitempty" hcl:"reason_phrase"`
}

type payload struct {
	XMLName        string `xml:"Payload" hcl:"-"`
	ContentType    string `xml:"contentType,attr,omitempty" hcl:"content_type"`
	VariablePrefix string `xml:"variablePrefix,attr,omitempty" hcl:"variable_prefix"`
	VariableSuffix string `xml:"variableSuffix,attr,omitempty" hcl:"variable_suffix"`
	Value          string `xml:",chardata" hcl:"value"`
}

type header struct {
	XMLName string `xml:"Header" hcl:"-"`
	Name    string `xml:"name,attr" hcl:"-"`
	Value   string `xml:",chardata" hcl:"value"`
}

type queryParam struct {
	XMLName string `xml:"QueryParam" hcl:"-"`
	Name    string `xml:"name,attr" hcl:"-"`
	Value   string `xml:",chardata" hcl:"value"`
}

type formParam struct {
	XMLName string `xml:"FormParam" hcl:"-"`
	Name    string `xml:"name,attr" hcl:"-"`
	Value   string `xml:",chardata" hcl:"value"`
}

func loadAssignMessageAddHCL(item *ast.ObjectItem) (*add, error) {
	result := new(add)

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("add not an object")
	}

	headers, err := loadHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	qparams, err := loadQueryParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.QueryParams = qparams

	fparams, err := loadFormParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.FormParams = fparams
	return result, nil
}

func loadAssignMessageCopyHCL(item *ast.ObjectItem) (*copy, error) {
	var result copy
	if err := hcl.DecodeObject(&result, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("copy not an object")
	}

	headers, err := loadHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	qparams, err := loadQueryParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.QueryParams = qparams

	fparams, err := loadFormParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.FormParams = fparams
	return &result, nil
}

func loadAssignMessageRemoveHCL(item *ast.ObjectItem) (*remove, error) {
	var result remove
	if err := hcl.DecodeObject(&result, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("remove not an object")
	}

	headers, err := loadHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	qparams, err := loadQueryParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.QueryParams = qparams

	fparams, err := loadFormParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.FormParams = fparams
	return &result, nil
}

func loadAssignMessageSetHCL(item *ast.ObjectItem) (*set, error) {
	var result set
	if err := hcl.DecodeObject(&result, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("set not an object")
	}

	headers, err := loadHeadersHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.Headers = headers

	qparams, err := loadQueryParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.QueryParams = qparams

	fparams, err := loadFormParamsHCL(listVal)
	if err != nil {
		return nil, err
	}

	result.FormParams = fparams
	return &result, nil
}

func loadHeadersHCL(listVal *ast.ObjectList) (*[]*header, error) {
	var headers []*header
	if headerList := listVal.Filter("header"); len(headerList.Items) > 0 {
		for _, h := range headerList.Items {
			var hdr header
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

func loadQueryParamsHCL(listVal *ast.ObjectList) (*[]*queryParam, error) {
	var qparams []*queryParam
	if qparamList := listVal.Filter("query_param"); len(qparamList.Items) > 0 {
		for _, q := range qparamList.Items {
			var qparam queryParam
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

func loadFormParamsHCL(listVal *ast.ObjectList) (*[]*formParam, error) {
	var fparams []*formParam
	if fparamList := listVal.Filter("form_param"); len(fparamList.Items) > 0 {
		for _, f := range fparamList.Items {
			var fparam formParam
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
