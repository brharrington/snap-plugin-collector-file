package file

import (
	"encoding/json"
	"io/ioutil"
	"github.com/intelsdi-x/snap/control/plugin"
	"fmt"
	"github.com/intelsdi-x/snap/core"
	"strings"
	"errors"
	"path/filepath"
	"strconv"
	"time"
)

type fileConfig struct {
	File    string                  `json:"file"`

	Metrics map[string]string    `json:"metrics"`

	Parser  parserConfig          `json:"parser"`
}

func fromJson(data []byte) (*[]fileConfig, error) {
	value := &[]fileConfig{}
	err := json.Unmarshal(data, value)
	return value, err
}

func fromJsonFile(file string) (*[]fileConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	} else {
		return fromJson(data)
	}
}

func isSubstitution(part string) bool {
	length := len(part)
	return length > 0 && part[0] == '{' && part[length - 1] == '}'
}

func filterEmpty(parts []string) []string {
	result := []string{}
	for _, p := range parts {
		p = strings.Trim(p, " \t\n")
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func toNamespace(pattern string) (*core.Namespace, error) {
	if len(pattern) == 0 {
		return nil, errors.New("namespace pattern cannot be empty")
	}

	if pattern[0] != '/' {
		msg := fmt.Sprintf("namespace pattern must begin with /: '%v'", pattern)
		return nil, errors.New(msg)
	}

	parts := filterEmpty(strings.Split(pattern[1:], "/"))
	numParts := len(parts)
	if numParts == 0 {
		msg := fmt.Sprintf("namespace pattern has no elements: '%v'", pattern)
		return nil, errors.New(msg)
	}

	if isSubstitution(parts[0]) {
		msg := fmt.Sprintf("namespace pattern must begin with static element: '%v'", pattern)
		return nil, errors.New(msg)
	}

	if isSubstitution(parts[numParts - 1]) {
		msg := fmt.Sprintf("namespace pattern must end with static element: '%v'", pattern)
		return nil, errors.New(msg)
	}

	ns := core.NewNamespace(parts[0])
	for _, p := range parts[1:] {
		if isSubstitution(p) {
			subst := p[1:len(p) - 1]
			vs := strings.Split(subst, ":")
			ns = ns.AddDynamicElement(vs[0], subst)
		} else {
			ns = ns.AddStaticElement(p)
		}
	}
	return &ns, nil
}

func (c fileConfig) getMetricTypes() ([]plugin.MetricType, error) {

	ms := []plugin.MetricType{}
	for k, _ := range c.Metrics {
		ns, err := toNamespace(k)
		if err != nil {
			return nil, err
		}
		ms = append(ms, plugin.MetricType{Namespace_: *ns})
	}

	return ms, nil
}

func createMetric(file string, vars map[string]interface{}, ns core.Namespace, valueExpr string) (*plugin.MetricType, error) {
	for i := 0; i < len(ns); i++ {
		if ns[i].IsDynamic() {
			parts := strings.Split(ns[i].Description, ":")
			if len(parts) == 3 && parts[1] == "path" {
				idx, err := strconv.Atoi(parts[2])
				if err != nil {
					return nil, err
				}

				path := strings.Split(file, "/")
				pos := idx
				if pos < 0 {
					pos = len(path) + pos
				}

				if pos < 0 || pos >= len(path) {
					msg := fmt.Sprintf("index '%d' out of bounds: %v", idx, path)
					return nil, errors.New(msg)
				}

				ns[i].Value = path[pos]
			} else {
				if value, ok := vars[ns[i].Name]; ok {
					ns[i].Value = fmt.Sprintf("%v", value)
				} else {
					msg := fmt.Sprintf("no value for dynamic element '%v'", ns[i].Name)
					return nil, errors.New(msg)
				}
			}

		}
	}

	value, err := eval(vars, valueExpr)
	if err != nil {
		return nil, err
	}

	m := plugin.NewMetricType(
		ns,
		time.Now(),
		map[string]string{},
		"",
		value,
	)
	return m, nil
}

func (c fileConfig) collectMetrics(queries []plugin.MetricType) ([]plugin.MetricType, error) {
	data := []plugin.MetricType{}
	files, err := filepath.Glob(c.File)
	if err != nil {
		return nil, err
	}
	fmt.Printf("loading files matching pattern [%s]: %v", c.File, files)

	for _, file := range files {
		fmt.Printf("loading file %s", file)
		parser := newParser(c.Parser)
		records, err := parser.parseFile(file)
		if err != nil {
			return nil, err
		}

		for _, record := range records {
			for k, v := range c.Metrics {
				ns, _ := toNamespace(k)
				m, err := createMetric(file, record, *ns, v)
				if err != nil {
					return nil, err
				}
				data = append(data, *m)
			}
		}
	}

	return data, nil
}
