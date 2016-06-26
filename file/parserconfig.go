package file

import (
	"regexp"
)

type parserConfig struct {
	format string

	columns []string

	recordSep string

	fieldSep string

	// Regular expression used to extract values for the 'regex' format.
	pattern string

	skip uint32
}

func defaultKeyValueConfig() parserConfig {
	return newKeyValueConfig("\n\n", "")
}

func newKeyValueConfig(recordSep string, fieldSep string) parserConfig {
	return parserConfig{
		"key-value",
		[]string{},
		recordSep,
		fieldSep,
		"",
		0,
	}
}

func newKeyRowConfig() parserConfig {
	return parserConfig{
		"key-row",
		[]string{},
		"",
		"",
		"",
		0,
	}
}

func newTableConfig(columns []string, skip uint32) parserConfig {
	return parserConfig{
		"table",
		columns,
		"",
		"",
		"",
		skip,
	}
}

func newRegexpConfig(columns []string, pattern string) parserConfig {
	return parserConfig{
		"regexp",
		columns,
		"\n",
		"",
		pattern,
		0,
	}
}

func (c parserConfig) regexp() *regexp.Regexp {
	return regexp.MustCompile(c.pattern)
}
