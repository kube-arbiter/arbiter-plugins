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

package wrapper

import (
	"context"
	"flag"
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	pb "github.com/kube-arbiter/arbiter/pkg/proto/lib/executor"
)

type ExecuteServiceImpl struct {
	pb.UnimplementedExecuteServer
}

var (
	_          pb.ExecuteServer = (*ExecuteServiceImpl)(nil)
	kubeconfig                  = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
)

func NewExecuteService() pb.ExecuteServer {
	return new(ExecuteServiceImpl)
}

func (e *ExecuteServiceImpl) Execute(ctx context.Context, message *pb.ExecuteMessage) (*pb.ExecuteResponse, error) {
	klog.V(10).Infof("kubeconfig path: %s\n", *kubeconfig)
	klog.V(4).Infof("ResourceName: %s, namespace: %s, exprval: %f, condval: %v, actionData: %v, executors: %v\n",
		message.ResourceName, message.Namespace, message.ExprVal, message.CondVal, message.ActionData, message.Executors)
	resourceBaseFormat := fmt.Sprintf("%s/%s/%s:%s", message.Group, message.Version, message.Resources, message.ResourceName)

	if len(message.Executors) == 0 {
		klog.Warningf("%s executor is empty, return..", resourceBaseFormat)
		return &pb.ExecuteResponse{}, nil
	}
	var (
		config   *rest.Config
		err      error
		response *pb.ExecuteResponse
	)
	if *kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		klog.Fatalf("error when building kubeconfig: %s", err.Error())
	}
	for _, executor := range message.Executors {
		instance, ok := GetExecutor(executor)
		if !ok {
			klog.Warningf("%s executor is %s, don't match target 'resourceUpdater'", resourceBaseFormat, executor)
			continue
		}

		if response, err = instance.Execute(ctx, config, message); err != nil {
			klog.Errorf("%s run %s error: %s\n", resourceBaseFormat, executor, err)
			break
		}
	}
	return response, err
}
