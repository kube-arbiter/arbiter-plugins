# Plugin(sidecar) Implementation Documentation

## Introduction

The main function of `observer` is to watch `observabilityIndicant`, to obtain metrics data periodically, and the plugin to be implemented is to serve as a data provider.

To implement a plugin, you only need to implement a `Grpc` interface. The interface only needs to contain three functions:


- **GetPluginName**
    
    This interface is mainly used to get the plugin name, so that the observer can ignore some crs that it does not pay attention to.

    
- **GetCapabilities**
    
    Get the capabilities supported by the current plugin, for example, get the maximum value of CPU over a period of time, or the average value of memory, etc.
    
    When the customer writes the aggregated value of the indicator to be obtained in cr, the observer will integrate the support ability obtained from the current plugin with what the user writes. For example:
    
    The user needs the aggregation of max and min with the cpu, but the plugin only provides the aggregation capability of min, so when the observer plugin requests data, it will only request the aggregation value of min.

- **GetMetrics**
    
    Returns aggregated metrics data for the target Pod.
    
After implementing the relevant interface, start the rpc process and deploy it together with the observer program (sidecar mode), the two communicate through the `unix domain socket`.

---

## Grpc Interface definition

```protobuf
syntax = "proto3";
package obi.v1;

option go_package = "lib/observer";

service Server {
    rpc GetPluginName (GetPluginNameRequest)
        returns (GetPluginNameResponse) {}
    
    rpc PluginCapabilities (PluginCapabilitiesRequest)
        returns (PluginCapabilitiesResponse) {}

    rpc GetMetrics (GetMetricsRequest)
        returns (GetMetricsResponse) {}
}


message GetPluginNameRequest {
    // nothing
}

message GetPluginNameResponse {
    string name = 1;
    map<string, string> attr = 2;
}

message PluginCapabilitiesRequest {
    // nothing
}

message MetricInfo {
    // cpu_usage
    string metric_unit = 1;
    string  description = 2;
    repeated string aggregation = 3; 
}

message PluginCapabilitiesResponse {
    // {"cpu_usage": {"metric_unit":"c","aggregation": ["a","b"]}}
    map<string, MetricInfo> metric_info = 2;
}

message GetMetricsRequest {
    // pod name or node name.
    // if there is a query field, ignore the resource_name field
    repeated string resource_names = 1;

    string namespace = 2;

    // cpu_usage, memory_usage eg.
    string metric_name = 3;

    // max, min, avg. And so on, some ways to aggregate data
    repeated string aggregation = 4;

    // such as promql
    string query = 5;

    // 1650867976563 million second
    // 这个请求时间做一个保留
    int64 start_time = 6;
    int64 end_time = 7;

    // resource kind, Pod, Node etc.
    string kind = 8;

    string unit = 9;
}

message GetMetricsResponse {
    // maybe empty,
    // For example, metrics-server can return resource_name,
    // which is the acquisition of the instantaneous value of the data.
    string resource_name = 1;

    string namespace = 2;

    // maybe empty
    string unit = 3;

    // the value field is string that is serialized by json.
    // {"cpu_usage": 0.1}
    message record {
        int64 timestamp = 1;
        string value = 2;
    }
    repeated record records = 4;
    //map<string,double> values = 4;
}
```

---

## Example

[metric-server](./observer-plugins/metric-server/), [prometheus](./observer-plugins/prometheus/)。

