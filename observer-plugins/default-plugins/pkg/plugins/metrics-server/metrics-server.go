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

package metricsserver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/kube-arbiter/arbiter-plugins/observer-plugins/default-plugins/pkg/flags"
	"github.com/kube-arbiter/arbiter-plugins/observer-plugins/default-plugins/pkg/plugins/resource"
	obi "github.com/kube-arbiter/arbiter/pkg/proto/lib/observer"
)

const (
	podMetricAPI  = `apis/metrics.k8s.io/v1beta1/namespaces/%s/pods/%s`
	nodeMetricAPI = `apis/metrics.k8s.io/v1beta1/nodes/%s`

	PluginName = "metrics-server"
	metricOpt  = "time"
	NodeKind   = "Node"
	PodKind    = "Pod"
)

type metricServer struct {
	client *kubernetes.Clientset
	cfg    *rest.Config
}

// NewMetricServer for register
func NewMetricServer(cfg *rest.Config) *metricServer {
	return &metricServer{
		cfg:    cfg,
		client: kubernetes.NewForConfigOrDie(cfg),
	}
}

func (ms *metricServer) Name() string {
	return PluginName
}

func (ms *metricServer) Capabilities() map[string]*obi.CapabilityInfo {
	return map[string]*obi.CapabilityInfo{
		"cpu": {
			MetricUnit:  "m",
			Description: "request pod or node cpu information from metrics server",
			Aggregation: []string{metricOpt},
		},
		"memory": {
			MetricUnit:  "byte",
			Description: "request pod or node memory information from metrics server",
			Aggregation: []string{metricOpt},
		},
	}
}

func (ms *metricServer) FetchData(ctx context.Context, req *obi.GetMetricsRequest) (*obi.GetMetricsResponse, error) {
	method := "metricServer/FetchData"

	returnObject := &obi.GetMetricsResponse{
		ResourceName: req.ResourceNames[0],
		Namespace:    req.Namespace,
		Unit:         req.Unit,
		Source:       PluginName,
	}

	if req.Kind != NodeKind && req.Kind != PodKind {
		klog.Warningf("[Error] %s don't support kind %s\n", method, req.Kind)
		return returnObject, nil
	}

	var calculate ResourceUsage
	switch req.Kind {
	case PodKind:
		queryPath := fmt.Sprintf(podMetricAPI, req.Namespace, req.ResourceNames[0])
		podMetricBytes, err := ms.client.RESTClient().Get().AbsPath(queryPath).DoRaw(ctx)
		if err != nil {
			return returnObject, err
		}
		podMetric := PodMetrics{}
		if err := json.Unmarshal(podMetricBytes, &podMetric); err != nil {
			klog.Errorf("[Error] unmarshal pod-metric error: %s\n", err)
			return returnObject, err
		}
		calculate = &podMetric
	case NodeKind:
		queryPath := fmt.Sprintf(nodeMetricAPI, req.ResourceNames[0])
		nodeMetricBytes, err := ms.client.RESTClient().Get().AbsPath(queryPath).DoRaw(ctx)
		if err != nil {
			klog.Errorf("[Error] failed to get node metric from metric-server: %s\n", err)
			return returnObject, err
		}
		nodeMetric := NodeMetrics{}
		if err := json.Unmarshal(nodeMetricBytes, &nodeMetric); err != nil {
			klog.Errorf("[Error] unmarshal node-metric error: %s\n", err)
			return returnObject, err
		}
		calculate = &nodeMetric

	default:
		klog.Errorf("[Error] don't support kind %s\n", req.Kind)
		return &obi.GetMetricsResponse{}, fmt.Errorf("do not support kind %s", req.Kind)
	}

	klog.V(4).Infof("Query: %s\n", req.Query)

	value, unit := calculate.SumResources(req.MetricName)
	if unit != "" {
		returnObject.Unit = unit
	}
	returnObject.Records = []*obi.GetMetricsResponseRecord{
		{
			Timestamp: time.Now().UnixMilli(),
			Value:     fmt.Sprintf("%.3f", float64(value)),
		},
	}
	return returnObject, nil
}

func init() {
	cfg, err := clientcmd.BuildConfigFromFlags("", *flags.Kubeconfig)
	if err == nil {
		instance := NewMetricServer(cfg)
		resource.Register(instance)
		klog.Infof("Observer [%s] registration is successful", PluginName)
	} else {
		klog.Warningf("Observer [%s] registration failed", PluginName)
	}
}
