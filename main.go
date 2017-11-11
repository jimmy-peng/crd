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
	//"flag"
	"github.com/jimmy-peng/crd/backend"
	"fmt"
	"github.com/rancher/longhorn-manager/types"
)


func main() {

	//kubeconf := flag.String("kubeconf", "admin.conf", "Path to a kube config. Only required if out-of-cluster.")
	//flag.Parse()
	kubeconf := ""
	backend, err := backend.NewCRDBackend(kubeconf)

	fmt.Printf("out CREATED: %#v\n", err)
	ee := types.ReplicaInfo{
		InstanceInfo: types.InstanceInfo{
			Name:   "lh.volumes-test",
			NodeID: "12345678",
			ID:     "adf",
		},
	}
	//longhorn_manager_test/volumes/bb/instances/replicas/bb-replica-62ef5c8b-f0de-4a94

	result, err := backend.Create( "/longhorn_manager_test/volumes/" + ee.Name + "/instances/replicas/" + ee.NodeID, ee)
	fmt.Printf("out CREATED: %d\n", result)

	e := types.ControllerInfo{
		InstanceInfo: types.InstanceInfo{
			Name: "lh.volumes-test",
			ID:   "123456786",
			IP:   "192.168.1.2",
		},
	}
	resul, err := backend.Create( "/longhorn_manager_test/nodes/" + e.Name + "/instances/controller", e)
	fmt.Printf("out CREATED: %d\n", resul)


	u := types.ReplicaInfo{
		InstanceInfo: types.InstanceInfo{
			Name:             "lh.volumes-hello",
			NodeID:           "12345678",
			ID: "8",
		},
	}

	resu, err := backend.Update("/longhorn_manager_test/volumes/" + u.Name + "/instances/replicas/"  + u.NodeID, u, result)
	if err == nil {
		fmt.Printf("out Update: %d\n", resu)
	}

	var ss types.ReplicaInfo
	r, err := backend.Get("/longhorn_manager_test/volumes/lh.volumes-test/instances/replicas/12345678", &ss)

	if err == nil {
		fmt.Printf("out GET: %#v %d\n", ss, r)
	}

	var s types.ControllerInfo
	ns, err := backend.Get("/longhorn_manager_test/volumes/lh.volumes-test/instances/controller", &s)

	if err == nil {
		fmt.Printf("out GET: %#v %d\n", s, ns)
	}

	rl, err := backend.Keys("/longhorn_manager_test/volumes/lh.volumes-test/instances/replicas/12345678")
	if err == nil {
		fmt.Printf("out List: %#v\n", rl)
	}
/*
	er := backend.Delete("/longhorn_manager_test/volumes/lh.volumes-test/base")
	if er != nil {
		fmt.Printf("out Delete %#v\n", er)
	}

	a, err := backend.Keys("/longhorn_manager_test/volumes")
	if err == nil {
		fmt.Printf("out after List: %#v\n", a)
	}*/
	// Wait forever
	select {}
}
