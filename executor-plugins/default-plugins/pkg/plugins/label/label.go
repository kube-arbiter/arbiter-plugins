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

package label

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	pb "github.com/kube-arbiter/arbiter/pkg/proto/lib/executor"
)

type LabelExecutor struct {
	name string
}

func (l *LabelExecutor) Name() string {
	return l.name
}

func (l *LabelExecutor) Execute(ctx context.Context, cfg *rest.Config, message *pb.ExecuteMessage) (*pb.ExecuteResponse, error) {
	resourceBaseFormat := fmt.Sprintf("%s/%s/%s:%s", message.Group, message.Version, message.Resources, message.ResourceName)

	dynamicClient, err := dynamic.NewForConfig(cfg)
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
		klog.Errorf("get resource %s (in namespace %s) error: %s\n", resourceBaseFormat, message.Namespace, err)
		if errors.IsNotFound(err) {
			response.Data = fmt.Sprintf("Resource %s not found in namespace %s", resourceBaseFormat, message.Namespace)
			return response, nil
		}
		response.Data = fmt.Sprintf("get resource %s error: %s", resourceBaseFormat, err)
		return response, err
	}

	err = UpdateResource(resouceToUpdate, message)
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

func NewLabelExecutor(name string) *LabelExecutor {
	return &LabelExecutor{name: name}
}
