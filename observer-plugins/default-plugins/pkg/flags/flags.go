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

package flags

import (
	"flag"

	"k8s.io/klog/v2"
)

// NOTE: if your metric resource need some paramer, please define it here.
var (
	Kubeconfig  = flag.String("kubeconfig", "", "kubernetes auth config file")
	Address     = flag.String("address", "", "prometheus server, such as http://localhost:9090")
	StepSeconds = flag.Int64("step", 60, "query steps")
	Endpoint    = flag.String("endpoint", "/var/run/observer.sock", "unix socket domain for current server")
)

func init() {
	klog.InitFlags(flag.CommandLine)
	flag.Parse()
}
