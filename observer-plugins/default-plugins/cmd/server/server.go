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
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"

	"github.com/kube-arbiter/arbiter-plugins/observer-plugins/default-plugins/pkg"
	"github.com/kube-arbiter/arbiter-plugins/observer-plugins/default-plugins/pkg/flags"
	obi "github.com/kube-arbiter/arbiter/pkg/proto/lib/observer"

	_ "github.com/kube-arbiter/arbiter-plugins/observer-plugins/default-plugins/pkg/install"
)

var (
	shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

func main() {
	_, err := os.Stat(*flags.Endpoint)
	if err != nil && !os.IsNotExist(err) {
		klog.Fatalln(err)
	}
	os.Remove(*flags.Endpoint)

	SetupSignalHandler(*flags.Endpoint)
	server := grpc.NewServer()
	obi.RegisterServerServer(server, pkg.NewServer())

	listen, err := net.Listen("unix", *flags.Endpoint)
	if err != nil {
		klog.Fatal(err)
	}

	klog.Infof("Observer plugin started ...\n")
	klog.Fatalln(server.Serve(listen))
}

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler(socketFile string) {
	c := make(chan os.Signal, 2)
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
