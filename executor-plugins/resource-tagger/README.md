# Plugin(sidecar) Implementation Documentation

## Introduction

The `executor` mainly watches the `ObservabilityActionPolicy` and determines how to add labels to `Pods` and `Nodes`
by calculating whether the results meet expectations using the expressions provided by the `ObservabilityActionPolicy`.

To implement a plugin, you only need to implement a `Grpc` interface. The interface only needs to contain three functions:


- **Execute**
    
    This interface is mainly used to add labesl to Pod or Node.

---

## Grpc Interface definition

```protobuf
syntax = "proto3";

package execute;
option go_package = "lib/executor";

service Execute {
    rpc Execute (ExecuteMessage) returns (ExecuteResponse) {}
}

enum Kind {
    pod = 0;
    node = 1;
}

message ExecuteMessage {
    string resourceName = 1;
    Kind kind = 2;
    string namespace = 3;
    double exprVal = 4;
    bool condVal = 5;
    map<string, string> parameters = 6;
    string group = 7;
    string version = 8;
    string resources = 9;
}

message ExecuteResponse {
    string data = 1;
}
```

