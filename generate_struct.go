package json_to_struct

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed go_struct_tmpl.txt
var structTemplate string

func ConvertYaml(name, rawData string) (map[string]*ConvertedStruct, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(rawData), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall data: %w", err)
	}

	v, err := convertInterface(name, data)
	if err != nil {
		return nil, err
	}
	Flatten(v)
	output := map[string]*ConvertedStruct{}
	for _, i := range v {
		output[i.StructName] = i
	}
	return output, nil
}

func ConvertJson(name, rawData string) (map[string]*ConvertedStruct, error) {
	var data interface{}
	err := json.Unmarshal([]byte(rawData), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall data: %w", err)
	}
	v, err := convertInterface(name, data)
	if err != nil {
		return nil, err
	}
	Flatten(v)
	output := map[string]*ConvertedStruct{}
	for _, i := range v {
		output[i.StructName] = i
	}
	return output, nil
}

func convertInterface(name string, data interface{}) ([]*ConvertedStruct, error) {
	var output []*ConvertedStruct
	switch x := data.(type) {
	case []interface{}:
		for _, i := range x {
			o, err := convertInterface(getVarName(name), i)
			if err != nil {
				return nil, err
			}
			output = append(output, o...)
			break
		}
	case map[string]interface{}:
		currentStruct := &ConvertedStruct{
			StructName: getVarName(name),
			Values:     []*StructType{},
		}
		for k, v := range x {
			switch sv := v.(type) {
			case int, int8, int16, int32, int64:
				currentStruct.Values = append(currentStruct.Values, &StructType{
					Type:     "int",
					Name:     getVarName(k),
					TagValue: k,
				})
			case string:
				currentStruct.Values = append(currentStruct.Values, &StructType{
					Type:     "string",
					Name:     getVarName(k),
					TagValue: k,
				})
			case float32, float64:
				currentStruct.Values = append(currentStruct.Values, &StructType{
					Type:     "float64",
					Name:     getVarName(k),
					TagValue: k,
				})
			case bool:
				currentStruct.Values = append(currentStruct.Values, &StructType{
					Type:     "bool",
					Name:     getVarName(k),
					TagValue: k,
				})
			case []interface{}:
				s := getStructType(k, sv)
				if s != nil {
					currentStruct.Values = append(currentStruct.Values, s)
				} else {
					o, err := convertInterface(getVarName(k), sv)
					if err != nil {
						return nil, err
					}
					output = append(output, o...)
					var t string
					if len(o) != 0 {
						t = fmt.Sprintf("*%s", getVarName(k))
					} else {
						t = "interface{}"
					}

					if isArray(sv) {
						t = fmt.Sprintf("[]%s", t)
					}

					currentStruct.Values = append(currentStruct.Values, &StructType{
						Type:     t,
						Name:     getVarName(k),
						TagValue: k,
					})

				}
			default:
				if sv != nil {
					o, err := convertInterface(getVarName(k), sv)
					if err != nil {
						return nil, err
					}
					output = append(output, o...)
				}
				t := fmt.Sprintf("*%s", getVarName(k))
				if isArray(sv) {
					t = fmt.Sprintf("[]%s", t)
				} else if sv == nil {
					t = "interface{}"
				}
				currentStruct.Values = append(currentStruct.Values, &StructType{
					Type:     t,
					Name:     getVarName(k),
					TagValue: k,
				})
			}
		}
		output = append(output, currentStruct)

	default:
		return nil, fmt.Errorf("Unsupported type: %T\n", x)
	}
	return output, nil
}
