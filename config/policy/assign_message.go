package policy

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type AssignMessagePolicy struct {
	XMLName                   string `xml:"AssignMessage" hcl:"-"`
	Policy                    `hcl:",squash"`
	DisplayName               string          `xml:",omitempty" hcl:"display_name"`
	Add                       *add            `xml:",omitempty" hcl:"add"`
	Copy                      *copy           `xml:",omitempty" hcl:"copy"`
	Remove                    *remove         `xml:",omitempty" hcl:"remove"`
	Set                       *set            `xml:",omitempty" hcl:"set"`
	AssignVariable            *assignVariable `xml:",omitempty" hcl:"assign_variable"`
	AssignTo                  *assignTo       `xml:",omitempty" hcl:"assign_to"`
	IgnoreUnresolvedVariables bool            `xml:",omitempty" hcl:"ignore_unresolved_variables"`
}

type add struct {
	XMLName     string        `xml:"Add" hcl:"-"`
	Headers     []*header     `xml:"Headers>Header" hcl:"header"`
	QueryParams []*queryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams  []*formParam  `xml:"FormParams>FormParam" hcl:"form_param"`
}

type copy struct {
	XMLName      string        `xml:"Copy" hcl:"-"`
	Source       string        `xml:"source,attr,omitempty" hcl:"-"`
	Headers      []*header     `xml:"Headers>Header" hcl:"header"`
	QueryParams  []*queryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams   []*formParam  `xml:"FormParams>FormParam" hcl:"form_param"`
	Payload      bool          `xml:",omitempty" hcl:"payload"`
	Version      bool          `xml:",omitempty" hcl:"version"`
	Verb         bool          `xml:",omitempty" hcl:"verb"`
	Path         bool          `xml:",omitempty" hcl:"path"`
	StatusCode   bool          `xml:",omitempty" hcl:"status_code"`
	ReasonPhrase bool          `xml:",omitempty" hcl:"reason_phrase"`
}

type remove struct {
	XMLName     string        `xml:"Remove" hcl:"-"`
	Headers     []*header     `xml:"Headers>Header" hcl:"header"`
	QueryParams []*queryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams  []*formParam  `xml:"FormParams>FormParam" hcl:"form_param"`
	Payload     bool          `xml:",omitempty" hcl:"payload"`
}

type set struct {
	XMLName      string        `xml:"Set" hcl:"-"`
	Headers      []*header     `xml:"Headers>Header" hcl:"header"`
	QueryParams  []*queryParam `xml:"QueryParams>QueryParam" hcl:"query_param"`
	FormParams   []*formParam  `xml:"FormParams>FormParam" hcl:"form_param"`
	Payload      payload       `xml:",omitempty" hcl:"payload"`
	Version      string        `xml:",omitempty" hcl:"version"`
	Verb         string        `xml:",omitempty" hcl:"verb"`
	Path         string        `xml:",omitempty" hcl:"path"`
	StatusCode   int           `xml:",omitempty" hcl:"status_code"`
	ReasonPhrase string        `xml:",omitempty" hcl:"reason_phrase"`
}

type assignVariable struct {
	Name  string `hcl:"name"`
	Ref   string `hcl:"ref"`
	Value string `hcl:"value"`
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
	Name    string `xml:"name,attr" hcl"-"`
	Value   string `xml:",chardata" hcl:"value"`
}

type queryParam struct {
	XMLName string `xml:"QueryParam" hcl:"-"`
	Name    string `xml:"name,attr" hcl"-"`
	Value   string `xml:",chardata" hcl:"value"`
}

type formParam struct {
	XMLName string `xml:"FormParam" hcl:"-"`
	Name    string `xml:"name,attr" hcl"-"`
	Value   string `xml:",chardata" hcl:"value"`
}

type assignTo struct {
	CreateNew bool   `xml:"createNew,attr" hcl:"create_new"`
	Transport string `xml:"transport,attr,omitempty" hcl:"transport"`
	Type      string `xml:"type,attr,omitempty" hcl:"type"`
	Value     string `xml:",chardata" hcl:"value"`
}

func LoadAssignMessageHCL(item *ast.ObjectItem) (interface{}, error) {
	var p AssignMessagePolicy

	if err := LoadCommonPolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("assign message policy not an object")
	}

	if addList := listVal.Filter("add"); len(addList.Items) > 0 {
		item := addList.Items[0]
		a, err := loadAssignMessageAddHCL(item)
		if err != nil {
			return nil, err
		}

		p.Add = a
	} else {
		p.Add = nil
	}

	if copyList := listVal.Filter("copy"); len(copyList.Items) > 0 {
		item := copyList.Items[0]
		a, err := loadAssignMessageCopyHCL(item)
		if err != nil {
			return nil, err
		}

		p.Copy = a
	} else {
		p.Copy = nil
	}

	if removeList := listVal.Filter("remove"); len(removeList.Items) > 0 {
		item := removeList.Items[0]
		a, err := loadAssignMessageRemoveHCL(item)
		if err != nil {
			return nil, err
		}

		p.Remove = a
	} else {
		p.Remove = nil
	}

	if setList := listVal.Filter("set"); len(setList.Items) > 0 {
		item := setList.Items[0]
		a, err := loadAssignMessageSetHCL(item)
		if err != nil {
			return nil, err
		}

		p.Set = a
	} else {
		p.Set = nil
	}

	return &p, nil
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

func loadHeadersHCL(listVal *ast.ObjectList) ([]*header, error) {
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
	}

	return headers, nil
}

func loadQueryParamsHCL(listVal *ast.ObjectList) ([]*queryParam, error) {
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
	}

	return qparams, nil
}

func loadFormParamsHCL(listVal *ast.ObjectList) ([]*formParam, error) {
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
	}

	return fparams, nil
}
