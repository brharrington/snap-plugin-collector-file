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

	//log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
)

const (
	name       = "file"
	version    = 1
	pluginType = plugin.CollectorPluginType
)

var _ plugin.CollectorPlugin = (*fileCollector)(nil)

type fileCollector struct {
}

func NewFileCollector() *fileCollector {
	return &fileCollector{}
}

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		name,
		version,
		pluginType,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType})
}

func (f *fileCollector) CollectMetrics(metrics []plugin.MetricType) ([]plugin.MetricType, error) {
	return nil, nil
}

func (f *fileCollector) GetMetricTypes(config plugin.ConfigType) ([]plugin.MetricType, error) {
	metricTypes := []plugin.MetricType{}
	return metricTypes, nil

}

func (f *fileCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {

	r1, err := cpolicy.NewStringRule("file", true)
	handleErr(err)
	r1.Description = "Main configuration file for the plugin."

	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()
	config.Add(r1)
	cp.Add([]string{""}, config)
	return cp, nil
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
