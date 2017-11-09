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
//	"fmt"
	"flag"
//	"github.com/jimmy-peng/crd/types/vtype"
	"github.com/jimmy-peng/crd/backend"
	"fmt"
	"github.com/jimmy-peng/crd/types/vtype"
	"github.com/jimmy-peng/crd/types/ntype"
)


func main() {

	kubeconf := flag.String("kubeconf", "admin.conf", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	backend, err := backend.NewCRDBackend(*kubeconf)
	/*
	// Create a new Example object and write to k8s
	example := vtype.Crdvolume{
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
	}*/

	fmt.Printf("out CREATED: %#v\n", err)
	ee := vtype.VolumeInfo{
			Name:   "lhvolumestest",
			NodeID: "12345678",
			NumberOfReplicas: 3,
	}

	result, err := backend.Create( "/longhorn_manager_test/volumes/" + ee.Name, ee)
	fmt.Printf("out CREATED: %d\n", result)

	e := ntype.NodeInfo{
		Name:   "lhvolumestest",
		NodeID: "12345678",
		NumberOfReplicas: 3,
	}
	resul, err := backend.Create( "/longhorn_manager_test/nodes/" + e.Name, e)
	fmt.Printf("out CREATED: %d\n", resul)
/*
	resu, err := backend.Update("/longhorn_manager_test/volumes/lhvolumestest", e, result)
	if err == nil {
		fmt.Printf("out Update: %d\n", resu)
	}

	var ss vtype.VolumeInfo
	r, err := backend.Get("/longhorn_manager_test/volumes/lhvolumestest", &ss)

	if err == nil {
		fmt.Printf("out GET: %#v %d\n", ss, r)
	}

	rl, err := backend.Keys("/longhorn_manager_test/volumes/lhvolumestest")
	if err == nil {
		fmt.Printf("out List: %#v\n", rl)
	}
*/
	// Wait forever
	select {}
}
