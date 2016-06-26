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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExprEval(t *testing.T) {

	Convey("eval", t, func() {
		vars := map[string]interface{}{
			"a": 1.0,
			"b": 42.0,
			"c": "foo",
		}

		value, err := eval(vars, "1.0")
		So(err, ShouldBeNil)
		So(value, ShouldResemble, 1.0)

		value, err = eval(vars, "1.0,2.0,:add")
		So(err, ShouldBeNil)
		So(value, ShouldResemble, 3.0)

		value, err = eval(vars, "1.0,2.0,:sub")
		So(err, ShouldBeNil)
		So(value, ShouldResemble, -1.0)

		value, err = eval(vars, "1.0,2.0,:mul")
		So(err, ShouldBeNil)
		So(value, ShouldResemble, 2.0)

		value, err = eval(vars, "1.0,2.0,:div")
		So(err, ShouldBeNil)
		So(value, ShouldResemble, 0.5)

		value, err = eval(vars, "{a},2.0,:add")
		So(err, ShouldBeNil)
		So(value, ShouldResemble, 3.0)

		value, err = eval(vars, "{b},2.0,:add")
		So(err, ShouldBeNil)
		So(value, ShouldResemble, 44.0)

		value, err = eval(vars, "{a},2.0,:add,{b},:mul")
		So(err, ShouldBeNil)
		So(value, ShouldResemble, 126.0)

		value, err = eval(vars, "foo")
		So(err, ShouldBeNil)
		So(value, ShouldResemble, "foo")

		value, err = eval(vars, "{c}")
		So(err, ShouldBeNil)
		So(value, ShouldResemble, "foo")

		value, err = eval(vars, "{c},4,:add")
		So(err.Error(), ShouldResemble, "not a number: 'foo' string")

		value, err = eval(vars, "4,:add")
		So(err.Error(), ShouldResemble, "need at least two arguments on the stack: [4]")

		value, err = eval(vars, "{c},4")
		So(err.Error(), ShouldResemble, "stack should have single item, found: [foo 4]")
	})

}
