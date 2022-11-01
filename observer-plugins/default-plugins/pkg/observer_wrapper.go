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

package pkg

import (
	"context"

	"k8s.io/klog/v2"

	"github.com/kube-arbiter/arbiter-plugins/observer-plugins/default-plugins/pkg/plugins/resource"
	obi "github.com/kube-arbiter/arbiter/pkg/proto/lib/observer"
)

type server struct {
	obi.UnimplementedServerServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) GetPluginNames(ctx context.Context, req *obi.GetPluginNameRequest) (*obi.GetPluginNameResponse, error) {
	return &obi.GetPluginNameResponse{
		Names: resource.Resources(),
	}, nil
}

func (s *server) PluginCapabilities(ctx context.Context, req *obi.PluginCapabilitiesRequest) (*obi.PluginCapabilitiesResponse, error) {
	plugins := resource.Resources()
	result := &obi.PluginCapabilitiesResponse{
		Capabilities: map[string]*obi.PluginCapability{},
	}
	for _, plugin := range plugins {
		if i, ok := resource.GetRegisters(plugin); ok {
			result.Capabilities[plugin] = &obi.PluginCapability{}
			result.Capabilities[plugin].Capability = i.Capabilities()
		}
	}

	return result, nil
}

func (s *server) GetMetrics(ctx context.Context, req *obi.GetMetricsRequest) (*obi.GetMetricsResponse, error) {
	klog.Infof("GetMetrics with req: %#v\n", req.String())
	if instance, ok := resource.GetRegisters(req.Source); ok {
		response, err := instance.FetchData(ctx, req)
		if err != nil {
			klog.Errorf("GetMetrics fetch data from %s error: %s\n", req.MetricName, err)
		}
		return response, err
	}

	klog.Warningf("GetMetrics request plugin %s, but it don't exists.", req.MetricName)
	return &obi.GetMetricsResponse{}, nil
}
