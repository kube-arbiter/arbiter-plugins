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

package resourceupdater

import (
	"context"
	"sync"

	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	pb "github.com/kube-arbiter/arbiter/pkg/proto/lib/executor"
)

type Executor interface {
	Name() string
	Execute(context.Context, *rest.Config, *pb.ExecuteMessage) (*pb.ExecuteResponse, error)
}

var (
	once         sync.Once
	mustRegister map[string]Executor
)

func Register(name string, instance Executor) {
	once.Do(func() {
		mustRegister = make(map[string]Executor)
	})
	if _, ok := mustRegister[name]; ok {
		klog.Warningf("Executor %s already exists", name)
	}

	mustRegister[name] = instance
}

func GetExecutor(name string) (Executor, bool) {
	v, ok := mustRegister[name]
	return v, ok
}
