package json_to_struct

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

func Flatten(ss []*ConvertedStruct) {
	for ik, iv := range ss {
		for jk, jv := range ss {
			if jk == ik {
				continue
			}
			if iv.Contains(jv) {
				replaceTrimStruct(iv, jv, ss)
				ss = append(ss[:jk], ss[jk+1:]...)
				Flatten(ss)
				return
			}
		}
	}
}
func replaceTrimStruct(new, old *ConvertedStruct, ss []*ConvertedStruct) {
	for i, c1 := range ss {
		for j, _ := range c1.Values {
			ss[i].Values[j].Type = strings.ReplaceAll(ss[i].Values[j].Type, old.StructName, new.StructName)
		}
	}
	return
}

func isArray(i interface{}) bool {
	switch i.(type) {
	case []interface{}:
		return true
	}
	return false
}

func getVarName(name string) string {
	return strcase.ToCamel(name)
}

func getArrayType(i []interface{}) interface{} {
	for _, v := range i {
		return v
	}
	return nil
}

func getStructType(k string, v interface{}) *StructType {
	switch sv := v.(type) {
	case int, int8, int16, int32, int64:
		return &StructType{
			Type:     "int",
			Name:     getVarName(k),
			TagValue: k,
		}
	case string:
		return &StructType{
			Type:     "string",
			Name:     getVarName(k),
			TagValue: k,
		}
	case float32, float64:
		return &StructType{
			Type:     "float64",
			Name:     getVarName(k),
			TagValue: k,
		}
	case bool:
		return &StructType{
			Type:     "bool",
			Name:     getVarName(k),
			TagValue: k,
		}
	case []interface{}:
		iv := getArrayType(sv)
		if iv != nil {
			st := getStructType(k, iv)
			if st != nil {
				st.Type = fmt.Sprintf("[]%s", st.Type)
				return st
			}
		}
	}
	return nil
}
