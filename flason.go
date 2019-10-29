package flason

import (
	"encoding/json"
	"strconv"
	"strings"
)

type JsonPair struct {
	Path  string
	Value string
}

func FlattenJson(str string) ([]JsonPair, error) {
	var pairs []JsonPair = make([]JsonPair, 0)

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return nil, err
	}

	// TODO: Change panic for error.
	// TODO: Make it asynchronous using go routines.
	var reduce func(prev string, curr map[string]interface{})
	reduce = func(prev string, curr map[string]interface{}) {
		for k, v := range curr {
			switch v.(type) {
			case string:
				path := strings.Join([]string{prev, ".", k}, "")
				pairs = append(pairs, JsonPair{
					Path:  path,
					Value: v.(string),
				})
			case float64:
				path := strings.Join([]string{prev, ".", k}, "")
				value := strconv.FormatFloat(v.(float64), 'g', -1, 64)
				pairs = append(pairs, JsonPair{
					Path:  path,
					Value: value,
				})
			case bool:
				path := strings.Join([]string{prev, ".", k}, "")
				value := strconv.FormatBool(v.(bool))
				pairs = append(pairs, JsonPair{
					Path:  path,
					Value: value,
				})
			case nil:
				path := strings.Join([]string{prev, ".", k}, "")
				pairs = append(pairs, JsonPair{
					Path: path,
					Value: "null",
				})
			case []interface{}:
				for index, w := range v.([]interface{}) {
					path := strings.Join([]string{ prev, "[", strconv.Itoa(index), "]"}, "")
					reduce(path , w.(map[string]interface{}))
				}
			case map[string]interface{}:
				for newKey, w := range v.(map[string]interface{}) {
					path := strings.Join([]string{prev, ".", newKey}, "")
					reduce(path, w.(map[string]interface{}))
				}
			default:
				panic("Unknown value type")
			}
		}
	}

	reduce("$", data)

	return pairs, nil
}
