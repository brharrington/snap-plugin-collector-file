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
	"fmt"
)

func TestFileConfig(t *testing.T) {

	Convey("fromJson", t, func() {

		configs, _ := fromJsonFile("testdata/fileconfig.json")

		for _, c := range *configs {
			ms, err := c.getMetricTypes()
			fmt.Println(ms)
			fmt.Println(err)

			mts, err := c.collectMetrics(nil)
			fmt.Println(err)
			for _, m := range mts {
				fmt.Println(m)
			}
		}

		//fmt.Println(config)
		//fmt.Println(err)
		//So(err, ShouldBeNil)
		//So(value, ShouldResemble, 1.0)
	})

}
