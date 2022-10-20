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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"

	pb "github.com/kube-arbiter/arbiter/pkg/proto/lib/executor"
)

/*
actionData defines the resource passed to the plugins, it can have arbitrary structure
You can use the .Raw data, and marshel a json object
*/
func (e *ExecuteServiceImpl) updateResource(resouceToUpdate *unstructured.Unstructured, message *pb.ExecuteMessage) (err error) {
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
