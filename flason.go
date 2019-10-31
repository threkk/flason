package flason

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type JsonPair struct {
	Path  string
	Value string
}

func FlattenJson(str string) (pairs []JsonPair, err error) {
	// We recover in case we panic during execution.
	defer func() {
		if r := recover(); r != nil {
			// TODO: Add debugging flag with the value.
			err = errors.New("Unknown value type during recursion.")
			pairs = nil
		}
	}()

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return nil, err
	}

	// TODO: Make it asynchronous using go routines.
	var reduce func(prev string, curr map[string]interface{})
	reduce = func(prev string, curr map[string]interface{}) {

		// Recover from the panic and return nicely.
		for k, v := range curr {
			switch v.(type) {
			case string:
				path := strings.Join([]string{prev, k}, "")
				pairs = append(pairs, JsonPair{
					Path:  path,
					Value: v.(string),
				})
			case float64:
				path := strings.Join([]string{prev, k}, "")
				value := strconv.FormatFloat(v.(float64), 'g', -1, 64)
				pairs = append(pairs, JsonPair{
					Path:  path,
					Value: value,
				})
			case bool:
				path := strings.Join([]string{prev, k}, "")
				value := strconv.FormatBool(v.(bool))
				pairs = append(pairs, JsonPair{
					Path:  path,
					Value: value,
				})
			case nil:
				path := strings.Join([]string{prev, k}, "")
				pairs = append(pairs, JsonPair{
					Path:  path,
					Value: "null",
				})
			case []interface{}:
				path := strings.Join([]string{prev, k}, "")

				sliceToMap := make(map[string]interface{})
				for index, w := range v.([]interface{}) {
					key := strings.Join([]string{"[", strconv.Itoa(index), "]"}, "")
					sliceToMap[key] = w
				}

				reduce(path, sliceToMap)
			case map[string]interface{}:
				path := strings.Join([]string{prev, k, "."}, "")
				reduce(path, v.(map[string]interface{}))

			default:
				panic(v)
			}
		}
	}

	reduce(".", data)
	return pairs, nil
}
