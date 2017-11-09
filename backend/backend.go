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
	"strings"
	"github.com/jimmy-peng/crd/clients/nclient"
	"github.com/jimmy-peng/crd/types/ntype"
)

type CRDBackend struct {
	VolumeClient *vclient.Crdclient
	NodeClient *nclient.Crdclient
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
		&vtype.Crdvolume{},
		time.Minute*10,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				/*
				r, ok := obj.(*vtype.Crdvolume)
				if ok {
					result, err := client.Create(r)
					if err != nil {
						if apierrors.IsAlreadyExists(err) {
							fmt.Printf("ALREADY EXISTS: %#v\n", result)
							return
						}
						fmt.Printf("create voleme error: %#v\n", err)
					}
				}
				*/
				fmt.Printf("add: %s \n", obj)
			},
			DeleteFunc: func(obj interface{}) {
				/*
				r, ok := obj.(*vtype.Crdvolume)
				if ok {
					client.Delete(r.ObjectMeta.Name, &meta_v1.DeleteOptions{})
				}*/
				fmt.Printf("del: %s \n", obj)
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
		NodeClient: nclient.CreateNodeClient(clientset, config),
	}

	CreateVolumeController(lhbackend.VolumeClient)
	return lhbackend, nil
}

func (s *CRDBackend) Create(key string, obj interface{}) (uint64, error) {
	r, ok := obj.(vtype.VolumeInfo)
	if ok {
		CRDobj := vtype.Crdvolume{}
		vtype.LhVoulme2CRDVolume(&r, &CRDobj, key)
		result, err := s.VolumeClient.Create(&CRDobj)
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				fmt.Printf("ALREADY EXISTS: %#v\n", result)
			}
			return 0, err
		}
		fmt.Printf("CREATED: %#v\n", result)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}

	n, ok := obj.(ntype.NodeInfo)
	if ok {
		CRDobj := ntype.Crdnode{}
		ntype.LhNode2CRDNode(&n, &CRDobj, key)
		result, err := s.NodeClient.Create(&CRDobj)
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
	r, ok := obj.(vtype.VolumeInfo)
	if ok {
		CRDobj := vtype.Crdvolume{}
		CRDobj.ResourceVersion = strconv.FormatUint(index, 10)
		vtype.LhVoulme2CRDVolume(&r, &CRDobj, key)
		validkey := strings.Split(key, "/")[3]
		result, err := s.VolumeClient.Update(&CRDobj, validkey)
		if err != nil {
			fmt.Printf("UPDATE: %#v\n", err)
			return 0, err
		}
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}
	n, ok := obj.(ntype.NodeInfo)
	if ok {
		CRDobj := ntype.Crdnode{}
		CRDobj.ResourceVersion = strconv.FormatUint(index, 10)
		ntype.LhNode2CRDNode(&n, &CRDobj, key)
		validkey := strings.Split(key, "/")[3]
		result, err := s.NodeClient.Update(&CRDobj, validkey)
		if err != nil {
			fmt.Printf("UPDATE: %#v\n", err)
			return 0, err
		}
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}
	return 0, nil
}

func (s *CRDBackend) Delete(key string) error {

	if strings.HasPrefix(key, "/longhorn_manager_test/volumes/") {
		validkey := strings.Split(key, "/")[3]
		err := s.VolumeClient.Delete(validkey, &meta_v1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	if strings.HasPrefix(key, "/longhorn_manager_test/nodes/") {
		validkey := strings.Split(key, "/")[3]
		err := s.NodeClient.Delete(validkey, &meta_v1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *CRDBackend) Get(key string, obj interface{}) (uint64, error) {
	r, ok := obj.(* vtype.VolumeInfo)
	if ok {
		validkey := strings.Split(key, "/")[3]
		result, err := s.VolumeClient.Get(validkey)
		vtype.CRDVolume2LhVoulme(result, r)
		if err != nil {
			return 0, err
		}
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}

	n, ok := obj.(* ntype.NodeInfo)
	if ok {
		validkey := strings.Split(key, "/")[3]
		result, err := s.NodeClient.Get(validkey)
		ntype.CRDNode2LhNode(result, n)
		if err != nil {
			return 0, err
		}
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}
	return 0, nil
}

func (s *CRDBackend) IsNotFoundError(err error) bool {
	return apierrors.IsNotFound(err)
}


func (s *CRDBackend) Keys(key string) ([]string, error) {
	if strings.HasPrefix(key, "/longhorn_manager_test/volumes/") {

		ret := []string{}
		r, err := s.VolumeClient.List(meta_v1.ListOptions{})
		if err != nil {
			return nil, err
		}
		if len(r.Items) <= 0 {
			return nil, nil
		}
		fmt.Printf("List: %#v\n", r)
		for _, item := range r.Items {
			if err != nil {
				return nil, err
			}
			ret = append(ret, item.ResourceVersion)
		}

		return ret, nil
	}
	if strings.HasPrefix(key, "/longhorn_manager_test/nodes/") {

		ret := []string{}
		r, err := s.NodeClient.List(meta_v1.ListOptions{})
		if err != nil {
			return nil, err
		}
		if len(r.Items) <= 0 {
			return nil, nil
		}
		fmt.Printf("List: %#v\n", r)
		for _, item := range r.Items {
			if err != nil {
				return nil, err
			}
			ret = append(ret, item.ResourceVersion)
		}

		return ret, nil
	}
	return nil, nil
}
