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
	"encoding/json"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"golang.org/x/net/context"
	"k8s.io/client-go/transport"
	"k8s.io/klog/v2"
)

var actionFuncs = map[string]func([]CalculateAux, *DataSeries){
	MaxAction: MaxOp,
	MinAction: MinOp,
	AvgAction: AvgOp,
}

func (p *prometheusServer) NewPrometheusAPI() (v1.API, error) {
	transConf, err := p.restConf.TransportConfig()
	if err != nil {
		return nil, err
	}
	rt, err := transport.New(transConf)
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(api.Config{
		Address:      p.address,
		RoundTripper: rt,
	})
	if err != nil {
		return nil, err
	}

	return v1.NewAPI(client), nil
}

type DataSeries struct {
	Timestamp int64
	Value     string
}

type CalculateAux struct {
	Timestamp int64
	Value     float64
}

func (p *prometheusServer) Query(startTime, endTime time.Time, kind, query, op string) (DataSeries, error) {
	method := "prometheusServer.Query"
	ans := DataSeries{Timestamp: endTime.UnixMilli()}
	prometheusAPI, err := p.NewPrometheusAPI()
	if err != nil {
		klog.Errorf("%s try to get prometheus API erorr: %s\n", method, err)
		return ans, err
	}
	result, warnings, err := prometheusAPI.QueryRange(context.TODO(), query, v1.Range{
		Start: startTime,
		End:   endTime,
		Step:  time.Duration(p.stepSeconds * int64(time.Second)),
	})
	if err != nil {
		klog.Errorf("%s try to query '%s' error: %s\n", method, query, err)
		return ans, err
	}
	if len(warnings) > 0 {
		klog.V(4).Infof("%s quer '%s' result with warnings %v\n", method, warnings)
	}

	// TODO: Use kind as the raw data query, may add a 'rawData: true' property for this?
	if kind == "Pod" || kind == "Node" {
		data, err := formatRawValues(result)
		if err != nil {
			return ans, err
		}
		if f, ok := actionFuncs[op]; ok {
			f(data, &ans)
		}
	} else {
		// Handle raw data if no aggregation defined, just return the json data
		jsonValue, err := json.Marshal(result)
		if err != nil {
			klog.Errorf("failed to marshal result to json: %s", err)
			ans.Value = fmt.Sprintf("failed to get json value: %s " + result.String())
		} else {
			ans.Value = string(jsonValue)
		}
	}
	return ans, nil
}

func formatRawValues(rawValue model.Value) ([]CalculateAux, error) {
	ans := make([]CalculateAux, 0)
	switch rawValue.Type() {
	case model.ValScalar:
		klog.V(4).Info("value type is Scalar\n")
		scalarObj, ok := rawValue.(*model.Scalar)
		if !ok {
			return ans, fmt.Errorf("can't conver to scaler")
		}
		ans = append(ans, CalculateAux{
			Timestamp: scalarObj.Timestamp.Time().UnixMilli(),
			Value:     float64(scalarObj.Value),
		})
	case model.ValMatrix:
		klog.V(4).Info("value type is matrix\n")
		matrixObj, ok := rawValue.(model.Matrix)
		if !ok {
			klog.Errorf("invalid matrix data")
			return ans, fmt.Errorf("can't convert to matrix")
		}

		klog.V(4).Infof("total rows: %d\n", len(matrixObj))
		for idx, v := range matrixObj {
			if idx > 0 {
				klog.V(4).Infof("Currently only supports getting a single data. values: %+v\n", v)
				continue
			}

			klog.V(5).Infof("values length: %d\n", len(v.Values))
			for _, sample := range v.Values {
				ans = append(ans, CalculateAux{
					Timestamp: sample.Timestamp.Time().UnixMilli(),
					Value:     float64(sample.Value),
				})
			}
		}
	case model.ValVector:
		klog.V(4).Info("value type is Vector\n")
		vector, ok := rawValue.(model.Vector)
		if !ok {
			klog.Infof("can't convert to vector")
			return ans, fmt.Errorf("can't conver to vector")
		}

		for _, sample := range vector {
			ans = append(ans, CalculateAux{
				Timestamp: sample.Timestamp.Time().UnixMicro(),
				Value:     float64(sample.Value),
			})
		}

	case model.ValString:
		klog.Errorf("string type response is not supported\n")
		return ans, fmt.Errorf("string type response isn't supported")
	}

	return ans, nil
}

func MaxOp(data []CalculateAux, result *DataSeries) {
	if len(data) == 0 {
		return
	}

	max := data[0]
	for idx := 1; idx < len(data); idx++ {
		if data[idx].Value > max.Value {
			max = data[idx]
		}
	}

	result.Timestamp = max.Timestamp
	result.Value = fmt.Sprintf("%f", max.Value)
}

func MinOp(data []CalculateAux, result *DataSeries) {
	if len(data) == 0 {
		return
	}

	min := data[0]
	for idx := 1; idx < len(data); idx++ {
		if data[idx].Value < min.Value {
			min = data[idx]
		}
	}

	result.Timestamp = min.Timestamp
	result.Value = fmt.Sprintf("%f", min.Value)
}

func AvgOp(data []CalculateAux, result *DataSeries) {
	if len(data) == 0 {
		return
	}

	avg := float64(0)
	for _, v := range data {
		avg += v.Value
	}

	result.Value = fmt.Sprintf("%f", avg/float64(len(data)))
}
