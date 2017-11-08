package backend

import (
	"github.com/jimmy-peng/crd/clients/vclient"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"github.com/jimmy-peng/crd/types/vtype"
	"strconv"
	"fmt"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/cache"
	"time"
)

type CRDBackend struct {
	VolumeClient *vclient.Crdclient
}

// return rest config, if path not specified assume in cluster config
func getClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func CreateVolumeController(client *vclient.Crdclient) {
	// Example Controller
	// Watch for changes in Example objects and fire Add, Delete, Update callbacks
	_, controller := cache.NewInformer(
		client.NewListWatch(),
		&vtype.Lhvolume{},
		time.Minute*10,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				fmt.Printf("add: %s \n", obj)
			},
			DeleteFunc: func(obj interface{}) {
				fmt.Printf("delete: %s \n", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				fmt.Printf("Update old: %s \n      New: %s\n", oldObj, newObj)
			},
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)
}

func NewCRDBackend(kubeconf string)(*CRDBackend, error)  {
	config, err := getClientConfig(kubeconf)
	if err != nil {
		panic(err.Error())
	}

	// create clientset and create our CRD, this only need to run once
	clientset, err := apiextcs.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	lhbackend := &CRDBackend{
		VolumeClient: vclient.CreateVolumeClient(clientset, config),
	}
	lhbackend.VolumeClient = vclient.CreateVolumeClient(clientset, config)

	CreateVolumeController(lhbackend.VolumeClient)
	return lhbackend, nil
}

func (s *CRDBackend) Create(key string, obj interface{}) (uint64, error) {
	r, ok := obj.(vtype.Lhvolume)
	if ok {
		result, err := s.VolumeClient.Create(&r)
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				fmt.Printf("ALREADY EXISTS: %#v\n", result)
			}
			return 0, err
		}
		fmt.Printf("CREATED: %#v\n", result)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}

	return 0, nil
}

func (s *CRDBackend) Update(key string, obj interface{}, index uint64) (uint64, error) {
	if index == 0 {
		return 0, fmt.Errorf("kvstore index cannot be 0")
	}
	r, ok := obj.(vtype.Lhvolume)
	if ok {
		result, err := s.VolumeClient.Update(&r)
		if err != nil {
			return 0, err
		}
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}

	return 0, nil
}

func (s *CRDBackend) Delete(key string) error {
	err := s.VolumeClient.Delete(key, &meta_v1.DeleteOptions{})

	if err != nil {
		return err
	}
	return nil
}

func (s *CRDBackend) Get(key string, obj interface{}) (uint64, error) {
	_, ok := obj.(vtype.Lhvolume)
	if ok {
		result, err := s.VolumeClient.Get(key)
		obj = result
		fmt.Printf("GET: %#v\n", result)
		if err != nil {
			return 0, err
		}
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}
	return 0, nil
}