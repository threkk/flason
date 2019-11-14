package flason

import (
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
)

type JsonPair struct {
	Path  string
	Value string
}

type byPath []JsonPair

func (p byPath) Len() int {
	return len(p)
}

func (p byPath) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p byPath) Less(i, j int) bool {
	return p[i].Path < p[j].Path
}

func FlattenJson(str, starter string) (pairs []JsonPair, err error) {
	// We recover in case we panic during execution.
	defer func() {
		if r := recover(); r != nil {
			// TODO: Add debugging flag with the value.
			err = errors.New("Unknown value type during recursion.")
			pairs = nil
		}
	}()

	var data interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return nil, err
	}

	// TODO: Make it asynchronous using go routines.
	var reduce func(prev string, curr interface{})
	reduce = func(prev string, curr interface{}) {

		// These cover all the possible JSON types.
		// https://golang.org/pkg/encoding/json/#Unmarshal
		switch curr.(type) {
		case string:
			pairs = append(pairs, JsonPair{
				Path:  prev,
				Value: curr.(string),
			})
		case float64:
			value := strconv.FormatFloat(curr.(float64), 'g', -1, 64)
			pairs = append(pairs, JsonPair{
				Path:  prev,
				Value: value,
			})
		case bool:
			value := strconv.FormatBool(curr.(bool))
			pairs = append(pairs, JsonPair{
				Path:  prev,
				Value: value,
			})
		case nil:
			pairs = append(pairs, JsonPair{
				Path:  prev,
				Value: "null",
			})
		case []interface{}:
			for index, value := range curr.([]interface{}) {
				path := strings.Join([]string{prev, "[", strconv.Itoa(index), "]"}, "")
				reduce(path, value)
			}

		case map[string]interface{}:
			for key, value := range curr.(map[string]interface{}) {
				path := strings.Join([]string{prev, ".", key}, "")
				reduce(path, value)
			}

		default:
			panic(curr)
		}
	}

	reduce(starter, data)
	sort.Sort(byPath(pairs))

	return pairs, nil
}
