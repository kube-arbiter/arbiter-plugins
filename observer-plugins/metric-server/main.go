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

	//"github.com/smoky8/pkg/lib/go/obi"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/kube-arbiter/arbiter-plugins/observer-plugins/metric-server/server"
	obi "github.com/kube-arbiter/arbiter/pkg/proto/lib/observer"
)

var (
	//metricServer = flag.String("metric-server", "localhost", "metric server address")
	endpoint   = flag.String("endpoint", "/var/run/observer.sock", "unix socket domain for current server")
	kubeconfig = flag.String("kubeconfig", "", "kubernetes auth config file")
)

func main() {
	klog.InitFlags(flag.CommandLine)
	flag.Parse()

	klog.Infoln("Start metric-server...")
	var (
		conf *rest.Config
		err  error
	)
	if *kubeconfig != "" {
		conf, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		conf, err = rest.InClusterConfig()
	}

	if err != nil {
		log.Fatal(err)
	}

	clientSet, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Fatalf("%s create metric client error: %s", server.PluginName, err)
	}
	metricServer := grpc.NewServer()

	obi.RegisterServerServer(metricServer, server.NewServer(clientSet))
	listen, err := net.Listen("unix", *endpoint)
	if err != nil {
		log.Fatal(err)
	}
	klog.Infof("%s starting work...", server.PluginName)

	klog.Fatalln(metricServer.Serve(listen))
}
