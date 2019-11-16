package flason

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type JsonPair struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

type FlatJson []JsonPair

func (p FlatJson) Len() int {
	return len(p)
}

func (p FlatJson) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p FlatJson) Less(i, j int) bool {
	return p[i].Path < p[j].Path
}

func (p FlatJson) PrintAsJson(f *os.File) error {
	content, err := json.Marshal(p)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	_, err = w.Write(content)
	if err != nil {
		return err
	}

	return w.Flush()
}

func (p FlatJson) PrintAsCsv(f *os.File) error {
	elements := [][]string{
		{"path", "value"},
	}

	for _, pair := range p {
		elements = append(elements, []string{
			pair.Path,
			pair.Value,
		})
	}

	w := csv.NewWriter(f)
	w.WriteAll(elements)

	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func (p FlatJson) PrintAsIni(f *os.File) error {
	w := bufio.NewWriter(f)

	var line string
	for _, pair := range p {
		line = fmt.Sprintf("%s = %s\n", pair.Path, pair.Value)
		_, err := w.WriteString(line)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}

func FlattenJson(str, starter string) (pairs FlatJson, err error) {
	// We recover in case we panic during execution.
	defer func() {
		if r := recover(); r != nil {
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
	sort.Sort(FlatJson(pairs))

	return pairs, nil
}
