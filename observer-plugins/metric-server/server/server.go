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

package server

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/net/context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	obi "github.com/kube-arbiter/arbiter/pkg/proto/lib/observer"
)

const (
	podMetricAPI  = `apis/metrics.k8s.io/v1beta1/namespaces/%s/pods/%s`
	nodeMetricAPI = `apis/metrics.k8s.io/v1beta1/nodes/%s`

	PluginName = "metric-server"
	metricOpt  = "time"
	NodeKind   = "Node"
	PodKind    = "Pod"
)

type server struct {
	obi.UnimplementedServerServer
	client *kubernetes.Clientset
}

func NewServer(client *kubernetes.Clientset) *server {
	return &server{client: client}
}

func (s *server) GetPluginName(ctx context.Context, req *obi.GetPluginNameRequest) (*obi.GetPluginNameResponse, error) {
	method := "GetPluginName"
	klog.V(4).Infof("%s %s\n", method, req.String())

	return &obi.GetPluginNameResponse{Name: PluginName}, nil
}

func (s *server) PluginCapabilities(ctx context.Context, req *obi.PluginCapabilitiesRequest) (*obi.PluginCapabilitiesResponse, error) {
	method := "PluginCapabilities"
	klog.V(4).Infof("%s %s\n", method, req.String())

	return &obi.PluginCapabilitiesResponse{
		MetricInfo: map[string]*obi.MetricInfo{
			"cpu": {
				MetricUnit:  "n",
				Description: "cpu usage",
				Aggregation: []string{metricOpt},
			},
			"memory": {
				MetricUnit:  "Ki",
				Description: "memory usage",
				Aggregation: []string{metricOpt},
			},
		},
	}, nil
}

func (s *server) GetMetrics(ctx context.Context, req *obi.GetMetricsRequest) (*obi.GetMetricsResponse, error) {
	method := "GetMetrics"
	klog.V(4).Infof("%s %s", method, req.String())

	returnOjb := &obi.GetMetricsResponse{
		ResourceName: req.ResourceNames[0],
		Namespace:    req.Namespace,
		Unit:         req.Unit,
	}
	if req.Kind != NodeKind && req.Kind != PodKind {
		klog.Errorf("[Error] don't support kind %s\n", req.Kind)
		return &obi.GetMetricsResponse{}, nil
	}
	var calculate ResourceUsage
	switch req.Kind {
	case PodKind:
		queryPath := fmt.Sprintf(podMetricAPI, req.Namespace, req.ResourceNames[0])
		podMetricBytes, err := s.client.RESTClient().Get().AbsPath(queryPath).DoRaw(ctx)
		if err != nil {
			return returnOjb, err
		}
		podMetric := PodMetrics{}
		if err := json.Unmarshal(podMetricBytes, &podMetric); err != nil {
			klog.Errorf("[Error] unmarshal pod-metric error: %s\n", err)
			return returnOjb, err
		}
		calculate = &podMetric
	case NodeKind:
		queryPath := fmt.Sprintf(nodeMetricAPI, req.ResourceNames[0])
		nodeMetricBytes, err := s.client.RESTClient().Get().AbsPath(queryPath).DoRaw(ctx)
		if err != nil {
			klog.Errorf("[Error] failed to get node metric from metric-server: %s\n", err)
			return returnOjb, err
		}
		nodeMetric := NodeMetrics{}
		if err := json.Unmarshal(nodeMetricBytes, &nodeMetric); err != nil {
			klog.Errorf("[Error] unmarshal node-metric error: %s\n", err)
			return returnOjb, err
		}
		calculate = &nodeMetric

	default:
		klog.Errorf("[Error] don't support kind %s\n", req.Kind)
		return &obi.GetMetricsResponse{}, fmt.Errorf("do not support kind %s", req.Kind)
	}

	klog.V(4).Infof("Query: %s\n", req.Query)

	value, unit := calculate.SumResources(req.MetricName)
	if unit != "" {
		returnOjb.Unit = unit
	}
	returnOjb.Records = []*obi.GetMetricsResponseRecord{
		{
			Timestamp: time.Now().UnixMilli(),
			Value:     fmt.Sprintf("%.3f", float64(value)),
		},
	}
	return returnOjb, nil
}
