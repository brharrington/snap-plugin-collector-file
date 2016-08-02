/*
 * Copyright 2016 Netflix, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package file

import (
	"fmt"
	"math"
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParser(t *testing.T) {

	Convey("parseValue", t, func() {
		So(parseValue("42"),           ShouldResemble, 42.0)
		So(parseValue("42.0"),         ShouldResemble, 42.0)
		So(parseValue("4.2e1"),        ShouldResemble, 42.0)
		So(parseValue("420e-1"),       ShouldResemble, 42.0)
		So(parseValue("42 \n"),        ShouldResemble, 42.0)
		So(parseValue("     \t42 \n"), ShouldResemble, 42.0)
		So(parseValue("1 2 3\n"),      ShouldResemble, "1 2 3")
		So(parseValue("foo bar\n"),    ShouldResemble, "foo bar")

		v := parseValue("NaN")
		So(math.IsNaN(v.(float64)), ShouldBeTrue)
	})

	Convey("parseKeyValue", t, func() {
		So(parseKeyValue("foo 42", ""), ShouldResemble, map[string]interface{}{
			"foo": 42.0,
		})

		So(parseKeyValue("  foo 42", ""), ShouldResemble, map[string]interface{}{
			"foo": 42.0,
		})

		So(parseKeyValue("  foo: 42", ""), ShouldResemble, map[string]interface{}{
			"foo": 42.0,
		})

		So(parseKeyValue("  foo   : 42", ""), ShouldResemble, map[string]interface{}{
			"foo": "",
		})

		So(parseKeyValue("  foo   : 42", ":"), ShouldResemble, map[string]interface{}{
			"foo": 42.0,
		})

		So(parseKeyValue("foo: 42 kB", ""), ShouldResemble, map[string]interface{}{
			"foo": 42.0,
		})

		So(parseKeyValue("foo 42 kB\nbar\t\t32   \n", ""), ShouldResemble, map[string]interface{}{
			"foo": 42.0,
			"bar": 32.0,
		})

		// Windows line feeds
		So(parseKeyValue("foo 42 kB\r\nbar\t\t32   \r\n", ""), ShouldResemble, map[string]interface{}{
			"foo": 42.0,
			"bar": 32.0,
		})
	})

	Convey("parseKeyValueList", t, func() {
		rows := parseKeyValueList("foo : 0\nbar : abc\n\nfoo : 1\nbar : 22", "\n\n", ":")
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"foo": 0.0,
				"bar": "abc",
			},
			map[string]interface{}{
				"foo": 1.0,
				"bar": 22.0,
			},
		})
	})

	Convey("parseKeyValueList", t, func() {
		rows := parseKeyValueList("foo 0.0\nbar abc\n\nfoo 1.0\nbar 22.0", "\n\n", " ")
		fmt.Println(rows)
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"foo": 0.0,
				"bar": "abc",
			},
			map[string]interface{}{
				"foo": 1.0,
				"bar": 22.0,
			},
		})
	})

	Convey("parseKeyRow", t, func() {
		rows, err := parseKeyRow("foo bar\nfoo 42")
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"id": "foo",
				"bar": 42.0,
			},
		})

		rows, err = parseKeyRow("foo: bar\nfoo: 42")
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"id": "foo",
				"bar": 42.0,
			},
		})

		rows, err = parseKeyRow("foo: a b\tc d\nfoo: 1 2 3 /bar/baz")
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"id": "foo",
				"a": 1.0,
				"b": 2.0,
				"c": 3.0,
				"d": "/bar/baz",
			},
		})

		rows, err = parseKeyRow("foo: a\nfoo: 42\nbar: d e f\nbar: 1 2 3\n")
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"id": "foo",
				"a": 42.0,
			},
			map[string]interface{}{
				"id": "bar",
				"d": 1.0,
				"e": 2.0,
				"f": 3.0,
			},
		})

		_, err = parseKeyRow("foo: a b\tc d\n")
		So(err.Error(), ShouldResemble, "key-row format requires even number of lines")

		_, err = parseKeyRow("foo: a b\tc d\nbar: 1 2 3 4")
		So(err.Error(), ShouldResemble, "line 0, rows ids do not match: 'foo' != 'bar'")

		_, err = parseKeyRow("foo: a b\tc d\nfoo: 1 2 3")
		So(err.Error(), ShouldResemble, "line 0, different number of columns: '5' != '4'")
	})

	Convey("parseTable", t, func() {
		rows, err := parseTable("foo bar\n1 2\n", []string{}, 0)
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"foo": 1.0,
				"bar": 2.0,
			},
		})

		rows, err = parseTable("foo bar\n1 2\n3 4\n5 6", []string{}, 0)
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"foo": 1.0,
				"bar": 2.0,
			},
			map[string]interface{}{
				"foo": 3.0,
				"bar": 4.0,
			},
			map[string]interface{}{
				"foo": 5.0,
				"bar": 6.0,
			},
		})

		rows, err = parseTable("1 2\n", []string{"foo", "bar"}, 0)
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"foo": 1.0,
				"bar": 2.0,
			},
		})

		rows, err = parseTable("a b\n1 2\n", []string{"foo", "bar"}, 0)
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"foo": "a",
				"bar": "b",
			},
			map[string]interface{}{
				"foo": 1.0,
				"bar": 2.0,
			},
		})

		rows, err = parseTable("a b\n1 2\n", []string{"foo", "bar"}, 1)
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"foo": 1.0,
				"bar": 2.0,
			},
		})

		rows, err = parseTable("a b\n1 2\n", []string{"foo", "bar"}, 2)
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
		})

		rows, err = parseTable("a b\n1 2\n", []string{"foo", "bar"}, 7)
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
		})
	})

	Convey("parseRegexp", t, func() {
		re := regexp.MustCompile("(\\d)\\s\\d (\\d)")
		rows, err := parseRegexp("foo bar baz\n1 2 3\na b c", "\n", []string{"a", "b"}, re)
		So(err, ShouldBeNil)
		So(rows, ShouldResemble, []map[string]interface{}{
			map[string]interface{}{
				"a": 1.0,
				"b": 3.0,
			},
		})
	})

	Convey("parse /proc/cpuinfo", t, func() {
		p := newParser(newKeyValueConfig("\n\n", ":"))
		rows, err := p.parseFile("testdata/cpuinfo")
		So(err, ShouldBeNil)
		So(len(rows), ShouldEqual, 2)
		for i, row := range rows {
			So(row["processor"], ShouldResemble, float64(i))
		}
	})

	Convey("parse /proc/loadavg", t, func() {
		p := newParser(newTableConfig([]string{"1m", "5m", "15m", "running/total", "last_pid"}, 0))
		rows, err := p.parseFile("testdata/loadavg")
		So(err, ShouldBeNil)
		So(len(rows), ShouldEqual, 1)
		So(rows[0]["1m"], ShouldResemble, 0.01)
		So(rows[0]["5m"], ShouldResemble, 0.05)
		So(rows[0]["15m"], ShouldResemble, 0.05)
		So(rows[0]["running/total"], ShouldResemble, "1/461")
		So(rows[0]["last_pid"], ShouldResemble, 13282.0)
	})

	Convey("parse /proc/net/netstat", t, func() {
		p := newParser(newKeyRowConfig())
		rows, err := p.parseFile("testdata/netstat")
		So(err, ShouldBeNil)
		So(len(rows), ShouldEqual, 2)

		So(rows[0]["id"], ShouldResemble, "TcpExt")
		So(rows[0]["EmbryonicRsts"], ShouldResemble, 3.0)

		So(rows[1]["id"], ShouldResemble, "IpExt")
		So(rows[1]["InNoRoutes"], ShouldResemble, 0.0)
	})

	Convey("parse /proc/net/dev", t, func() {
		columns := []string{
			"interface",
			"recv_bytes",
			"recv_packets",
			"recv_errs",
			"recv_drop",
			"recv_fifo",
			"recv_frame",
			"recv_compressed",
			"recv_multicast",
			"send_bytes",
			"send_packets",
			"send_errs",
			"send_drop",
			"send_fifo",
			"send_colls",
			"send_carrier",
			"send_compressed",
		}
		p := newParser(newTableConfig(columns, 2))
		rows, err := p.parseFile("testdata/net_dev")
		So(err, ShouldBeNil)
		So(len(rows), ShouldEqual, 11)

		So(rows[0]["interface"], ShouldResemble, "veth929074f")
		So(rows[9]["interface"], ShouldResemble, "lo")

		So(rows[8]["recv_packets"], ShouldResemble, 10032.0)
	})

	Convey("parse /proc/stat cpu", t, func() {
		columns := []string{
			"label",
			"user",
			"nice",
			"system",
			"idle",
			"iowait",
			"irq",
			"softirq",
			"steal",
			"guest",
			"guest_nice",
		}
		pattern := "(cpu\\d*)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)"
		p := newParser(newRegexpConfig(columns, pattern))
		rows, err := p.parseFile("testdata/stat")
		So(err, ShouldBeNil)
		So(len(rows), ShouldEqual, 3)

		So(rows[1]["label"], ShouldResemble, "cpu0")
		So(rows[2]["irq"], ShouldResemble, 122.0)
	})
}
