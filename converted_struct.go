package json_to_struct

import (
	"bytes"
	"go/format"
	"html/template"
	"strings"
)

type ConvertedStruct struct {
	StructName string
	Values     []*StructType
}

type StructType struct {
	Type     string
	Name     string
	TagValue string
}

func (c *ConvertedStruct) containsField(st *StructType) bool {
	for _, c1 := range c.Values {
		if strings.EqualFold(c1.Name, st.Name) && strings.EqualFold(c1.Type, st.Type) {
			return true
		}
	}
	return false
}

func (c *ConvertedStruct) Contains(v2 *ConvertedStruct) bool {
	for _, c2 := range v2.Values {
		if !c.containsField(c2) {
			return false
		}
	}
	return true
}

func (c *ConvertedStruct) ToString() (string, error) {
	tmpl, err := template.New("to_struct").Parse(structTemplate)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, c)

	p, err := format.Source(tpl.Bytes())
	if err != nil {
		return "", err
	}
	return string(p), err
}
