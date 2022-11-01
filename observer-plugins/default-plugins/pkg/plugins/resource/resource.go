/*
Copyright 2022 The Arbiter Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resource

import (
	"context"
	"sync"

	"k8s.io/klog/v2"

	obi "github.com/kube-arbiter/arbiter/pkg/proto/lib/observer"
)

var (
	once         sync.Once
	mustRegister map[string]Observer
)

func Register(instance Observer) {
	once.Do(func() {
		mustRegister = make(map[string]Observer)
	})

	name := instance.Name()
	if _, ok := mustRegister[name]; ok {
		klog.Warningf("Observer %s already exists", name)
	}

	mustRegister[name] = instance
}

type Observer interface {
	Name() string
	FetchData(context.Context, *obi.GetMetricsRequest) (*obi.GetMetricsResponse, error)
	Capabilities() map[string]*obi.CapabilityInfo
}

func GetRegisters(name string) (Observer, bool) {
	v, ok := mustRegister[name]
	return v, ok
}

func Resources() []string {
	pluginNames := make([]string, len(mustRegister))
	i := 0

	for k := range mustRegister {
		pluginNames[i] = k
		i++
	}
	return pluginNames
}
