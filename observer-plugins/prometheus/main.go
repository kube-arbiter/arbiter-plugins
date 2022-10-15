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

package main

import (
	"flag"
	"log"
	"net"

	"github.com/kube-arbiter/arbiter-plugins/observer-plugins/prometheus/prometheus"
	obi "github.com/kube-arbiter/arbiter/pkg/proto/lib/observer"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var (
	endpoint    = flag.String("endpoint", "/var/run/observer.sock", "unix socket domain for current server")
	kubeConfig  = flag.String("kubeconfig", "", "kubernetes auth config file")
	address     = flag.String("address", "", "prometheus server, such as http://localhost:9090")
	stepSeconds = flag.Int64("step", 60, "query steps")
	rangeMinute = flag.Int64("range", 2, "prometheus, the maximum time between two slices within the boundaries.")
)

func main() {
	klog.InitFlags(flag.CommandLine)
	flag.Parse()

	if *address == "" {
		klog.Fatalf("prometheus serve address can not be empty")
	}
	var (
		conf *rest.Config
		err  error
	)
	if *kubeConfig != "" {
		conf, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
	} else {
		conf, err = rest.InClusterConfig()
	}

	if err != nil {
		klog.Fatal(err)
	}

	_, err = kubernetes.NewForConfig(conf)
	if err != nil {
		klog.Fatal(err)
	}

	server := grpc.NewServer()
	obi.RegisterServerServer(server, prometheus.NewPrometheusServer(*address, conf, *stepSeconds, *rangeMinute))
	listen, err := net.Listen("unix", *endpoint)
	if err != nil {
		log.Fatal(err)
	}

	klog.Infof("%s starting work...", prometheus.PluginName)
	server.Serve(listen)
}
