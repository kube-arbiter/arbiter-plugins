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
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/kube-arbiter/arbiter-plugins/observer-plugins/prometheus/prometheus"
	obi "github.com/kube-arbiter/arbiter/pkg/proto/lib/observer"
)

var (
	endpoint    = flag.String("endpoint", "/var/run/observer.sock", "unix socket domain for current server")
	kubeConfig  = flag.String("kubeconfig", "", "kubernetes auth config file")
	address     = flag.String("address", "", "prometheus server, such as http://localhost:9090")
	stepSeconds = flag.Int64("step", 60, "query steps")
	rangeMinute = flag.Int64("range", 2, "prometheus, the maximum time between two slices within the boundaries.")
)
var (
	shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
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
	// Setup signal watcher to handle cleanup
	SetupSignalHandler(*endpoint)

	server := grpc.NewServer()
	obi.RegisterServerServer(server, prometheus.NewPrometheusServer(*address, conf, *stepSeconds, *rangeMinute))
	listen, err := net.Listen("unix", *endpoint)
	if err != nil {
		log.Fatal(err)
	}

	klog.Infof("%s plugin started ...", prometheus.PluginName)
	klog.Fatalln(server.Serve(listen))
}

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler(socketFile string) {
	c := make(chan os.Signal)
	signal.Notify(c, shutdownSignals...)
	go func() {
		for s := range c {
			switch s {
			case os.Interrupt, syscall.SIGTERM:
				klog.Infoln("Shutting down normally...")
				if err := os.RemoveAll(socketFile); err != nil {
					klog.Fatal(err)
				}
				os.Exit(1)
			default:
				klog.Infoln("Got signal", s)
			}
		}
	}()
}
