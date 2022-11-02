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

package prometheus

import (
	"time"

	"golang.org/x/net/context"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/kube-arbiter/arbiter-plugins/observer-plugins/default-plugins/pkg/flags"
	"github.com/kube-arbiter/arbiter-plugins/observer-plugins/default-plugins/pkg/plugins/resource"
	obi "github.com/kube-arbiter/arbiter/pkg/proto/lib/observer"
)

const (
	PluginName = "prometheus"
	MaxAction  = "max"
	MinAction  = "min"
	AvgAction  = "avg"
	NoneAction = "none"
)

// impl obi interface
type prometheusServer struct {
	obi.UnimplementedServerServer
	address     string
	restConf    *rest.Config
	stepSeconds int64
}

func NewPrometheusServer(address string, restConf *rest.Config, stepSeconds int64) *prometheusServer {
	method := "NewPrometheusServer"
	klog.V(4).Infof("%s stepSecond: %d\n", method, stepSeconds)
	return &prometheusServer{
		address:     address,
		restConf:    restConf,
		stepSeconds: stepSeconds,
	}
}

func (p *prometheusServer) Name() string {
	return PluginName
}

func (p *prometheusServer) Capabilities() map[string]*obi.CapabilityInfo {
	return map[string]*obi.CapabilityInfo{
		"cpu": {
			MetricUnit:  "m",
			Description: "request pod or node cpu information from prometheus",
			Aggregation: []string{MaxAction, MinAction, AvgAction},
		},
		"memory": {
			MetricUnit:  "byte",
			Description: "request pod or node memory information from prometheus",
			Aggregation: []string{MaxAction, MinAction, AvgAction},
		},
	}
}

func (p *prometheusServer) FetchData(ctx context.Context, req *obi.GetMetricsRequest) (*obi.GetMetricsResponse, error) {
	method := "prometheusServer/FetchData"
	klog.V(4).Infof("%s req %s\n", method, req.String())

	startTime := time.Unix(0, req.StartTime*int64(time.Millisecond))
	endTime := time.Unix(0, req.EndTime*int64(time.Millisecond))

	var err error
	klog.V(4).Infof("prometheus query: %s\n", req.Query)
	var resourceName string
	if len(req.ResourceNames) > 0 {
		resourceName = req.ResourceNames[0]
	}
	result := &obi.GetMetricsResponse{
		ResourceName: resourceName,
		Namespace:    req.Namespace,
		Unit:         req.Unit,
		Records:      []*obi.GetMetricsResponseRecord{},
	}

	// use avgerage as the default aggregation action
	op := AvgAction
	if len(req.Aggregation) > 0 {
		op = req.Aggregation[0]
	}
	metricData, err := p.Query(startTime, endTime, req.Kind, req.Query, op)
	if err != nil {
		klog.Errorf("%s query error: %s\n", method, err)
		return result, err
	}
	// only return the latest record
	result.Records = append(result.Records, &obi.GetMetricsResponseRecord{Timestamp: metricData.Timestamp, Value: metricData.Value})

	klog.Infof("query by metric '%s', query '%s' successfully", req.MetricName, req.Query)
	klog.V(5).Infof("%s query by %s, %s result: %v\n", method, req.MetricName, req.Query, metricData)

	return result, nil
}

func init() {
	cfg, err := clientcmd.BuildConfigFromFlags("", *flags.Kubeconfig)
	if err == nil {
		instance := NewPrometheusServer(*flags.Address, cfg, *flags.StepSeconds)
		resource.Register(instance)
		klog.Infof("Observer [%s] registration is successful", PluginName)
	} else {
		klog.Warningf("Observer [%s] registration failed", PluginName)
	}
}
