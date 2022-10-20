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

package service

import (
	"context"
	"flag"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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
	klog.V(4).Infof("ResourceName: %s, namespace: %s, exprval: %f, condval: %v, actionData: %v\n",
		message.ResourceName, message.Namespace, message.ExprVal, message.CondVal, message.ActionData)

	resourceBaseFormat := fmt.Sprintf("%s/%s/%s:%s", message.Group, message.Version, message.Resources, message.ResourceName)
	var (
		config *rest.Config
		err    error
	)

	if *kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		klog.Fatalf("error when building kubeconfig: %s", err.Error())
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	var resouceToUpdate *unstructured.Unstructured
	response := &pb.ExecuteResponse{
		Data: "",
	}
	namespaceableInterface := dynamicClient.Resource(
		schema.GroupVersionResource{Group: message.Group, Version: message.Version, Resource: message.Resources})
	if message.Namespace != "" {
		resouceToUpdate, err = namespaceableInterface.Namespace(message.Namespace).Get(context.Background(), message.ResourceName, metav1.GetOptions{})
	} else {
		resouceToUpdate, err = namespaceableInterface.Get(context.Background(), message.ResourceName, metav1.GetOptions{})
	}
	if err != nil {
		klog.Errorf("get resource %s (int namespace %s) error: %s\n", resourceBaseFormat, message.Namespace, err)
		if errors.IsNotFound(err) {
			response.Data = fmt.Sprintf("Resource %s not found in namespace %s", resourceBaseFormat, message.Namespace)
			return response, nil
		}
		response.Data = fmt.Sprintf("get resource %s error: %s", resourceBaseFormat, err)
		return response, err
	}
	// Let the custom code to handle how to update the resource
	err = e.updateResource(resouceToUpdate, message)
	if err != nil {
		response.Data = err.Error()
		return response, err
	}
	if message.Namespace != "" {
		_, err = namespaceableInterface.Namespace(message.Namespace).Update(context.Background(), resouceToUpdate, metav1.UpdateOptions{})
	} else {
		_, err = namespaceableInterface.Update(context.Background(), resouceToUpdate, metav1.UpdateOptions{})
	}
	if err != nil {
		response.Data = fmt.Sprintf("update resource %s error: %s", resourceBaseFormat, err)
		return response, err
	}
	return response, nil
}
