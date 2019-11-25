package flason

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func uniq(src []string) []string {
	seen := make(map[string]bool)

	for _, str := range src {
		seen[str] = true
	}

	result := make([]string, len(seen))
	for key := range seen {
		result = append(result, key)
	}

	return result
}

type JSONPair struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

type FlatJSON []JSONPair

func (p FlatJSON) Len() int {
	return len(p)
}

func (p FlatJSON) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p FlatJSON) Less(i, j int) bool {
	return p[i].Path < p[j].Path
}

func (p FlatJSON) PrintAsJSON(f *os.File) error {
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

func (p FlatJSON) PrintAsCSV(f *os.File) error {
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

func (p FlatJSON) PrintAsINI(f *os.File) error {
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

func (p FlatJSON) PrintOnlyPath(f *os.File, unique bool) error {
	w := bufio.NewWriter(f)

	paths := make([]string, len(p))
	for i, pair := range p {
		paths[i] = fmt.Sprintf("%s\n", pair.Path)
		if !unique {
			w.WriteString(paths[i])
		}
	}

	if unique {
		r, err := regexp.Compile(`\[\d+\]`)
		if err != nil {
			return err
		}

		for i, pair := range p {
			paths[i] = fmt.Sprintf("%s\n", r.ReplaceAllString(pair.Path, "[]"))
		}

		paths = uniq(paths)
		sort.Strings(paths)
		for _, p := range paths {
			w.WriteString(p)
		}
	}

	return w.Flush()
}

func FlattenJSON(str, starter string) (pairs FlatJSON, err error) {
	// We recover in case we panic during execution.
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unknown value type during recursion")
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
			pairs = append(pairs, JSONPair{
				Path:  prev,
				Value: curr.(string),
			})
		case float64:
			value := strconv.FormatFloat(curr.(float64), 'g', -1, 64)
			pairs = append(pairs, JSONPair{
				Path:  prev,
				Value: value,
			})
		case bool:
			value := strconv.FormatBool(curr.(bool))
			pairs = append(pairs, JSONPair{
				Path:  prev,
				Value: value,
			})
		case nil:
			pairs = append(pairs, JSONPair{
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
	sort.Sort(FlatJSON(pairs))

	return pairs, nil
}
