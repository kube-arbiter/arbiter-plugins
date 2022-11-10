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

package plugins

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	pb "github.com/kube-arbiter/arbiter/pkg/proto/lib/executor"
)

type ResourceUpdateExecutor struct {
	name string
}

func (l *ResourceUpdateExecutor) Name() string {
	return l.name
}

/*
actionData defines the resource passed to the plugins, it can have arbitrary structure
You can use the .Raw data, and marshel a json object

Here is an example executor plugin to show how to update the associated resource using the data from OBI & OAP
it'll add/remove labels from the value that evaluated by OAP
*/
func UpdateResource(resouceToUpdate *unstructured.Unstructured, message *pb.ExecuteMessage) (err error) {
	resourceName := fmt.Sprintf("%s/%s", resouceToUpdate.GetKind(), resouceToUpdate.GetName())
	klog.Infof("start processing resource %s", resourceName)

	metaObj := metav1.ObjectMeta{}
	if err = json.Unmarshal(message.ActionData.Raw, &metaObj); err != nil {
		klog.Errorf("Failed to unmarshal the raw message %s with error %s", message.ActionData.Raw, err)
		return err
	}
	labels := resouceToUpdate.GetLabels()
	if message.CondVal {
		if labels == nil {
			labels = make(map[string]string)
		}
		for key, value := range metaObj.Labels {
			labels[key] = value
		}
		klog.Infof("Resource %s in namesapce '%s' is labeled", resourceName, message.Namespace)
	} else {
		for key := range metaObj.Labels {
			delete(labels, key)
		}
		klog.Infof("Resource %s in namespace '%s' is un-labeled", resourceName, message.Namespace)
	}
	resouceToUpdate.SetLabels(labels)
	klog.Infof("%s updated successfully.", resourceName)
	return nil
}

func (l *ResourceUpdateExecutor) Execute(ctx context.Context, cfg *rest.Config, message *pb.ExecuteMessage) (*pb.ExecuteResponse, error) {
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

func NewResourceUpdateExecutor(name string) *ResourceUpdateExecutor {
	return &ResourceUpdateExecutor{name: name}
}
