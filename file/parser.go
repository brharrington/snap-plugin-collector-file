package file

import (
	"strconv"
	"strings"
	"regexp"
	"io/ioutil"
	"fmt"
	"errors"
)

type parser struct {
	config parserConfig
}

func newParser(config parserConfig) parser {
	return parser{
		config,
	}
}

func (p parser) parseFile(file string) ([]map[string]interface{}, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return p.parseString(string(data))
}

func (p parser) parseString(data string) ([]map[string]interface{}, error) {
	switch p.config.format {
	case "table":
		return parseTable(data, p.config.columns, p.config.skip)
	case "key-value":
		return parseKeyValueList(data, p.config.recordSep, p.config.fieldSep), nil
	case "key-row":
		return parseKeyRow(data)
	case "regexp":
		return parseRegexp(data, p.config.recordSep, p.config.columns, p.config.regexp())
	default:
		return nil, errors.New(fmt.Sprintf("unknown file format: '%v'", p.config.format))
	}
}

func parseValue(data string) interface{} {
	tmp := strings.Trim(data, ": \t\r\n")
	v, err := strconv.ParseFloat(tmp, 64)
	if err == nil {
		return v
	} else {
		return tmp
	}
}

func parseKeyValue(data string, fieldSep string) map[string]interface{} {
	lines := strings.Split(data, "\n")
	values := map[string]interface{}{}
	for _, line := range lines {
		var fields []string
		if fieldSep == "" {
			fields = strings.Fields(line)
		} else {
			fields = strings.Split(line, fieldSep)
		}

		// Lines with less than two fields are ignored
		if len(fields) >= 2 {
			// If the field name ends with a ':' strip it out
			k := strings.Trim(fields[0], " :")
			values[k] = parseValue(fields[1])
		}
	}
  return values
}

func parseKeyValueList(data string, recordSep string, fieldSep string) []map[string]interface{} {
	items := []map[string]interface{}{}
	records := strings.Split(data, recordSep)
	for _, record := range records {
		item := parseKeyValue(record, fieldSep)
		if len(item) > 0 {
			items = append(items, item)
		}
	}
	return items
}

func parseKeyRow(data string) ([]map[string]interface{}, error) {
	rows := []map[string]interface{}{}
	lines := strings.Split(strings.Trim(data, "\n"), "\n")
	if len(lines) % 2 != 0 {
		return rows, errors.New("key-row format requires even number of lines")
	}

	for i := 0; i < len(lines); i += 2 {
		headers := strings.Fields(lines[i])
		values := strings.Fields(lines[i + 1])

		hlen := len(headers)
		vlen := len(values)

		// Lines cannot be empty. Maybe these should be ignored instead?
		if hlen == 0 {
			return rows, errors.New(fmt.Sprintf("line %v, empty lines are not allowed", i))
		}

		// Number of headers must match number of values
		if hlen != vlen {
			msg := fmt.Sprintf("line %v, different number of columns: '%v' != '%v'", i, hlen, vlen)
			return rows, errors.New(msg)
		}

		// Check that the ids are the same
		hid := strings.Trim(headers[0], ":")
		vid := strings.Trim(values[0], ":")
		if hid != vid {
			msg := fmt.Sprintf("line %v, rows ids do not match: '%v' != '%v'", i, hid, vid)
			return rows, errors.New(msg)
		}

		row := map[string]interface{}{
			"id": hid,
		}

		for j := 1; j < hlen; j++ {
			row[headers[j]] = parseValue(values[j])
		}

		rows = append(rows, row)
	}
	return rows, nil
}

func parseTable(data string, columns []string, skip uint32) ([]map[string]interface{}, error) {
	rows := []map[string]interface{}{}
	lines := strings.Split(strings.Trim(data, "\n"), "\n")
	if int(skip) <= len(lines) {
		lines = lines[skip:]
	} else {
		return rows, nil
	}

	headers := columns
	if len(headers) == 0 {
		headers = strings.Fields(lines[0])
		lines = lines[1:]
	}

	for _, line := range lines {
		values := strings.Fields(line)
		if len(values) == len(headers) {
			row := map[string]interface{}{}
			for i, v := range values {
				row[headers[i]] = parseValue(v)
			}
			rows = append(rows, row)
		}
	}

	return rows, nil
}

func parseRegexp(data string, recordSep string, columns []string, pattern *regexp.Regexp) ([]map[string]interface{}, error) {
	items := []map[string]interface{}{}
	records := strings.Split(data, recordSep)

	clen := len(columns)
	for i, record := range records {
		values := pattern.FindStringSubmatch(record)
		if values == nil {
			continue
		}
		values = values[1:]

		vlen := len(values)
		if clen != vlen {
			msg := fmt.Sprintf("record %v, different number of columns: '%v' != '%v'", i, clen, vlen)
			return items, errors.New(msg)
		}

		item := map[string]interface{}{}
		for j, v := range values {
			item[columns[j]] = parseValue(v)
		}
		items = append(items, item)
	}

	return items, nil
}
