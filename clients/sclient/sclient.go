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
package sclient

import (
	apiv1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/runtime"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"github.com/jimmy-peng/crd/crds/scrd"
	"github.com/jimmy-peng/crd/types/stype"
)

type Crdclient struct {
	cl     *rest.RESTClient
	ns     string
	plural string
	codec  runtime.ParameterCodec
}
// This file implement all the (CRUD) client methods we need to access our CRD object

func CrdClient(cl *rest.RESTClient, scheme *runtime.Scheme) *Crdclient {
	return &Crdclient{cl: cl, ns: apiv1.NamespaceDefault, plural: scrd.CRDPlural,
		codec: runtime.NewParameterCodec(scheme)}
}

func CreateSettingClient(clientset apiextcs.Interface, cfg *rest.Config) *Crdclient {
	// note: if the CRD exist our CreateCRD function is set to exit without an error
	err := scrd.CreateSettingCRD(clientset)
	if err != nil {
		panic(err)
	}

	// Wait for the CRD to be created before we use it (only needed if its a new one)
	err = scrd.SettingWaitCRDCreateDone(clientset)
	if err != nil {
		panic(err)
	}

	// Create a new clientset which include our CRD schema
	crdcs, scheme, err := scrd.SettingNewClient(cfg)
	if err != nil {
		panic(err)
	}

	// Create a CRD client interface
	return CrdClient(crdcs, scheme)
}



func (f *Crdclient) Create(obj *stype.Crdsetting) (*stype.Crdsetting, error) {
	var result stype.Crdsetting
	err := f.cl.Post().
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) Update(obj *stype.Crdsetting, name string) (*stype.Crdsetting, error) {
	result := stype.Crdsetting{}
	err := f.cl.Put().Name(name).
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *Crdclient) Delete(name string, options *meta_v1.DeleteOptions) error {
	return f.cl.Delete().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Body(options).Do().
		Error()
}

func (f *Crdclient) Get(name string) (*stype.Crdsetting, error) {
	result := stype.Crdsetting{}
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Do().Into(&result)
	return &result, err
}


func (f *Crdclient) GetByVersion(version string) (*stype.Crdsetting, error) {

	slist, err := f.List(meta_v1.ListOptions{})
	for _,item := range slist.Items {
		if item.ResourceVersion == version {
			return &item, err
		}
	}
	return nil, err
}

func (f *Crdclient) List(opts meta_v1.ListOptions) ( *stype.CrdsettingList, error) {
	result := stype.CrdsettingList{}
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		VersionedParams(&opts, f.codec).
		Do().Into(&result)
	return &result, err
}

// Create a new List watch for our TPR
func (f *Crdclient) NewListWatch() *cache.ListWatch {
	return cache.NewListWatchFromClient(f.cl, f.plural, f.ns, fields.Everything())
}
