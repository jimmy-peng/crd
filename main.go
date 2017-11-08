/*
Copyright 2016 Iguazio Systems Ltd.

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
	"fmt"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"flag"
	"github.com/jimmy-peng/crd/types/vtype"
	"github.com/jimmy-peng/crd/backend"
)


func main() {

	kubeconf := flag.String("kubeconf", "admin.conf", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	backend, err := backend.NewCRDBackend(*kubeconf)

	// Create a new Example object and write to k8s
	example := vtype.Lhvolume{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   "lhvolumestest",
			Labels: map[string]string{"mylabel": "test"},
		},
		Spec: vtype.LhVolumeSpec{
			Name: "valumes",
			NodeId: "12345678",
			NumReplicas: 3,
		},
		Status: vtype.LhVolumeStatus{
			State:   "created",
			Message: "Created, not processed yet",
		},
	}

	result, err := backend.Create(example.Name, example)
	fmt.Printf("CREATED: %d\n", result)
	var ss vtype.Lhvolume
	resul, err := backend.Get("lhvolumestest", ss)
	if err == nil {
		fmt.Printf("GET: %d\n", resul)
	}

	// Wait forever
	select {}
}
