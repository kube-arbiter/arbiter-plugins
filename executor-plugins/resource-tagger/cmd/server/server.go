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

	"google.golang.org/grpc"
	"k8s.io/klog/v2"

	"github.com/kube-arbiter/arbiter-plugins/executor-plugins/resource-tagger/internal/service"
	pb "github.com/kube-arbiter/arbiter/pkg/proto/lib/executor"
)

const (
	protocol = "unix"
	sockAddr = "/plugins/resourcetagger.sock"
)

func main() {
	// Load flags from command line
	klog.InitFlags(nil)
	flag.Parse()

	cleanup := func() {
		if _, err := os.Stat(sockAddr); err == nil {
			if err := os.RemoveAll(sockAddr); err != nil {
				log.Fatal(err)
			}
		}
	}

	cleanup()

	listener, err := net.Listen(protocol, sockAddr)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	execute := service.NewExecuteService()

	pb.RegisterExecuteServer(server, execute)

	klog.Infoln("resource-tagger plugin started...")
	klog.Fatalln(server.Serve(listener))
}
