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
	"github.com/rancher/longhorn-manager/types"
	"github.com/jimmy-peng/crd/clients/rclient"
	"github.com/jimmy-peng/crd/types/rtype"
	"github.com/jimmy-peng/crd/types/ctype"
	"github.com/jimmy-peng/crd/clients/cclient"
	"github.com/jimmy-peng/crd/clients/sclient"
	"github.com/jimmy-peng/crd/types/stype"
	"github.com/pkg/errors"
	"github.com/rancher/longhorn-manager/manager"
)

type CRDBackend struct {
	VolumeClient *vclient.Crdclient
	NodeClient *nclient.Crdclient
	ReplicasClient *rclient.Crdclient
	ControllerClient *cclient.Crdclient
	SettingClient *sclient.Crdclient
}

// return rest config, if path not specified assume in cluster config
func getClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}


func CreateVolumeController(m *manager.VolumeManager, b *CRDBackend ) {
	// Example Controller
	// Watch for changes in Example objects and fire Add, Delete, Update callbacks
	_, controller := cache.NewInformer(
		b.VolumeClient.NewListWatch(),
		&vtype.Crdvolume{},
		time.Minute*10,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {

				v, ok := obj.(*vtype.Crdvolume)
				if ok && v.Spec.TargetNodeID == "" {
					var result types.VolumeInfo
					vtype.CRDVolume2LhVoulme(v, &result)
					m.CRDVolumeCreate(&result, v.ObjectMeta.Name)
					fmt.Printf("add get: %#v \n", result)
				}

			},

			DeleteFunc: func(obj interface{}) {
				v, ok := obj.(*vtype.Crdvolume)
				if ok {
					var result types.VolumeInfo
					vtype.CRDVolume2LhVoulme(v, &result)
					m.CRDVolumeDelete(&result)
				}
				fmt.Printf("del: %s \n", obj)
			},

			UpdateFunc: func(oldObj, newObj interface{}) {
				nv := newObj.(*vtype.Crdvolume)
				ov, ok := oldObj.(*vtype.Crdvolume)
				if ok && ov.Spec.DesireState != nv.Spec.DesireState {
					var or types.VolumeInfo
					var nr types.VolumeInfo
					vtype.CRDVolume2LhVoulme(ov, &or)
					vtype.CRDVolume2LhVoulme(nv, &nr)
					kindex, err := strconv.ParseUint(nv.ObjectMeta.ResourceVersion, 10, 64)
					if err != nil {
						fmt.Println("Parse index error")
						return
					}
					m.CRDVolumeAttachDetach(&or, &nr, kindex)
				}
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
		ReplicasClient: rclient.CreateReplicaClient(clientset, config),
		ControllerClient: cclient.CreateControllerClient(clientset, config),
		SettingClient: sclient.CreateSettingClient(clientset, config),
	}

	//CreateVolumeController(lhbackend.VolumeClient)
	return lhbackend, nil
}

func (s *CRDBackend) Create(key string, obj interface{}) (uint64, error) {
	v, ok := obj.(*types.VolumeInfo)
	if ok {
		CRDobj := vtype.Crdvolume{}
		vtype.LhVoulme2CRDVolume(v, &CRDobj)
		result, err := s.VolumeClient.Create(&CRDobj)
		if err != nil {
			fmt.Printf("ERROR: %#v\n", err)
			if apierrors.IsAlreadyExists(err) {
				fmt.Printf("ALREADY EXISTS: %#v\n", result)
			}
			return 0, err
		}
		fmt.Printf("CREATED: %#v\n", result)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}

	n, ok := obj.(*types.NodeInfo)
	if ok {
		CRDobj := ntype.Crdnode{}
		validkey := strings.Split(key, "/")[3]
		ntype.LhNode2CRDNode(n, &CRDobj, validkey)
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


	r, ok := obj.(*types.ReplicaInfo)
	if ok {
		CRDobj := rtype.Crdreplica{}
		validkey := strings.Split(key, "/")[6]
		rtype.LhReplicas2CRDReplicas(r, &CRDobj, validkey)
		result, err := s.ReplicasClient.Create(&CRDobj)
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				fmt.Printf("ALREADY EXISTS: %#v\n", result)
			}
			return 0, err
		}
		fmt.Printf("CREATED: %#v\n", result)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}


	c, ok := obj.(*types.ControllerInfo)
	if ok {
		CRDobj := ctype.Crdcontroller{}
		validkey := strings.Split(key, "/")[3]
		ctype.LhController2CRDController(c, &CRDobj, validkey)
		result, err := s.ControllerClient.Create(&CRDobj)
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				fmt.Printf("ALREADY EXISTS: %#v\n", result)
			}
			return 0, err
		}
		fmt.Printf("CREATED: %#v\n", result)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}

	set, ok := obj.(* types.SettingsInfo)
	if ok {
		CRDobj := stype.Crdsetting{}
		validkey := strings.Split(key, "/")[2]
		stype.LhSetting2CRDSetting(set, &CRDobj, validkey)
		result, err := s.SettingClient.Create(&CRDobj)
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				fmt.Printf("ALREADY EXISTS: %#v\n", result)
			}
			return 0, err
		}
		fmt.Printf("CREATED: %#v\n", result)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}
	panic(key)
	return 0, nil
}

func (s *CRDBackend) Update(key string, obj interface{}, index uint64) (uint64, error) {
	fmt.Println("Update key string", key, index)
	if index == 0 {
		return 0, fmt.Errorf("kvstore index cannot be 0")
	}
	v, ok := obj.(*types.VolumeInfo)
	if ok {
		CRDobj := vtype.Crdvolume{}
		CRDobj.ResourceVersion = strconv.FormatUint(index, 10)
		vtype.LhVoulme2CRDVolume(v, &CRDobj)
		validkey := strings.Split(key, "/")[3]
		result, err := s.VolumeClient.Update(&CRDobj, validkey)
		if err != nil {
			return 0, err
		}
		//fmt.Printf("UPDATE: %#v\n", result)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}
	n, ok := obj.(*types.NodeInfo)
	if ok {
		CRDobj := ntype.Crdnode{}
		CRDobj.ResourceVersion = strconv.FormatUint(index, 10)
		validkey := strings.Split(key, "/")[3]
		ntype.LhNode2CRDNode(n, &CRDobj, validkey)
		result, err := s.NodeClient.Update(&CRDobj, validkey)
		if err != nil {
			return 0, err
		}
		//fmt.Printf("UPDATE: %#v\n", CRDobj)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}


	r, ok := obj.(*types.ReplicaInfo)
	if ok {
		CRDobj := rtype.Crdreplica{}
		CRDobj.ResourceVersion = strconv.FormatUint(index, 10)
		validkey := strings.Split(key, "/")[6]
		rtype.LhReplicas2CRDReplicas(r, &CRDobj, validkey)
		result, err := s.ReplicasClient.Update(&CRDobj, validkey)
		if err != nil {

			return 0, err
		}
	//	fmt.Printf("UPDATE: %#v\n", CRDobj)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}


	c, ok := obj.(*types.ControllerInfo)
	if ok {
		CRDobj := ctype.Crdcontroller{}
		CRDobj.ResourceVersion = strconv.FormatUint(index, 10)
		validkey := strings.Split(key, "/")[3]
		ctype.LhController2CRDController(c, &CRDobj, validkey)
		result, err := s.ControllerClient.Update(&CRDobj, validkey)
		if err != nil {

			return 0, err
		}
	//	fmt.Printf("UPDATE: %#v\n", CRDobj)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}


	ss, ok := obj.(*types.SettingsInfo)
	if ok {
		CRDobj := stype.Crdsetting{}
		CRDobj.ResourceVersion = strconv.FormatUint(index, 10)
		validkey := strings.Split(key, "/")[2]
		stype.LhSetting2CRDSetting(ss, &CRDobj, validkey)
		result, err := s.SettingClient.Update(&CRDobj, validkey)
		if err != nil {
			return 0, err
		}
	//	fmt.Printf("UPDATE: %#v\n", CRDobj)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}
	panic(key)
	return 0, nil
}

func (s *CRDBackend) Delete(key string) error {
	fmt.Println("Delete key string %s", key)
	if strings.HasPrefix(key, "/longhorn_manager_test/volumes/") {
		validkey := strings.Split(key, "/")[3]
		err := s.VolumeClient.Delete(validkey, &meta_v1.DeleteOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	if strings.HasPrefix(key, "/longhorn_manager_test/nodes/") {
		validkey := strings.Split(key, "/")[3]
		err := s.NodeClient.Delete(validkey, &meta_v1.DeleteOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	if strings.HasPrefix(key, "/longhorn_manager_test/volumes/") &&
		strings.Contains(key, "/instances/replicas/") {
		validkey := strings.Split(key, "/")[6]
		err := s.ReplicasClient.Delete(validkey, &meta_v1.DeleteOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	if strings.HasPrefix(key, "/longhorn_manager_test/volumes/") &&
		strings.HasSuffix(key, "/instances/controller") {
		validkey := strings.Split(key, "/")[3]
		err := s.ControllerClient.Delete(validkey, &meta_v1.DeleteOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	if key == "/longhorn_manager_test/settings" {
		validkey := strings.Split(key, "/")[2]
		err := s.SettingClient.Delete(validkey, &meta_v1.DeleteOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	panic(key)
	return nil
}

func (s *CRDBackend) Get(key string, obj interface{}) (uint64, error) {
	fmt.Println("GET key string", key)
	v, ok := obj.(* types.VolumeInfo)
	if ok {
		var result *vtype.Crdvolume
		var err error
		if strings.Contains(key, "volumes") {
			validkey := strings.Split(key, "/")[3]
			result, err = s.VolumeClient.Get(validkey)
		} else {
			validkey := strings.Split(key, "/")[0]
			result, err = s.VolumeClient.GetByVersion(validkey)
		}

		if err != nil || result == nil {
			return 0, errors.New("crd not found")
		}

		vtype.CRDVolume2LhVoulme(result, v)
		//fmt.Printf("GET volume string %#v \n\n", v)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}

	n, ok := obj.(* types.NodeInfo)
	if ok {
		var result *ntype.Crdnode
		var err error
		if strings.Contains(key, "nodes") {
			validkey := strings.Split(key, "/")[3]
			result, err = s.NodeClient.Get(validkey)
		} else {

			result, err = s.NodeClient.GetByVersion(key)
		}

		if err != nil || result == nil {
			return 0, errors.New("crd not found")
		}
		ntype.CRDNode2LhNode(result, n)
		//fmt.Printf("GET node string %#v \n\n", n)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)

	}

	r, ok := obj.(* types.ReplicaInfo)
	if ok {
		var result *rtype.Crdreplica
		var err error
		if strings.Contains(key, "volumes") {
			validkey := strings.Split(key, "/")[6]
			result, err = s.ReplicasClient.Get(validkey)
		} else {
			result, err = s.ReplicasClient.GetByVersion(key)
		}
		if err != nil || result == nil{
			return 0, errors.New("crd not found")
		}

		rtype.CRDReplicas2LhReplicas(result, r)
		//fmt.Printf("GET replica string %#v \n\n", r)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}


	c, ok := obj.(* types.ControllerInfo)
	if ok {
		var result *ctype.Crdcontroller
		var err error

		if strings.Contains(key, "volumes") {
			validkey := strings.Split(key, "/")[3]
			result, err = s.ControllerClient.Get(validkey)
		} else {
			result, err = s.ControllerClient.GetByVersion(key)
		}

		if err != nil || result == nil{
			return 0, errors.New("crd not found")
		}
		fmt.Printf("GET controller string %#v \n\n", c)
		ctype.CRDController2LhController(result, c)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}

	ss, ok := obj.(* types.SettingsInfo)
	if ok {
		var result *stype.Crdsetting
		var err error

		if strings.Contains(key, "settings") {
			validkey := strings.Split(key, "/")[2]
			result, err = s.SettingClient.Get(validkey)
		} else {
			result, err = s.SettingClient.GetByVersion(key)
		}

		if err != nil || result == nil{
			return 0, errors.New("crd not found")
		}
		fmt.Printf("GET setting string %#v \n\n", ss)
		stype.CRDSetting2LhSetting(result, ss)
		return strconv.ParseUint(result.ResourceVersion, 10, 64)
	}
	panic(key)
	return 0, nil
}

func (s *CRDBackend) IsNotFoundError(err error) bool {
	if strings.Contains(err.Error(), "crd not found"){
		return true
	} else {
		return false
	}
}


func (s *CRDBackend) Keys(key string) ([]string, error) {
	if key == "/longhorn_manager_test/volumes" {
		ret := []string{}
		r, err := s.VolumeClient.List(meta_v1.ListOptions{})
		if err != nil {
			return nil, err
		}
		if len(r.Items) <= 0 {
			return nil, nil
		}
		//fmt.Printf("List: %#v\n", r)
		for _, item := range r.Items {
			if err != nil {
				return nil, err
			}
			ret = append(ret, item.ResourceVersion)
		}

		return ret, nil
	}

	if key == "/longhorn_manager_test/nodes" {

		ret := []string{}
		r, err := s.NodeClient.List(meta_v1.ListOptions{})
		if err != nil {
			return nil, err
		}
		if len(r.Items) <= 0 {
			return nil, nil
		}
		//fmt.Printf("List: %#v\n", r)
		for _, item := range r.Items {
			if err != nil {
				return nil, err
			}
			ret = append(ret, item.ResourceVersion)
		}

		return ret, nil
	}
	if strings.HasPrefix(key, "/longhorn_manager_test/volumes/") &&
		strings.Contains(key, "/instances/replicas") {

		ret := []string{}
		r, err := s.ReplicasClient.List(meta_v1.ListOptions{})
		if err != nil {
			return nil, err
		}
		if len(r.Items) <= 0 {
			return nil, nil
		}
		//fmt.Printf("List: %#v\n", r)
		for _, item := range r.Items {
			if err != nil {
				return nil, err
			}
			ret = append(ret, item.ResourceVersion)
		}

		return ret, nil
	}

	if strings.HasPrefix(key, "/longhorn_manager_test/volumes/") &&
		strings.HasSuffix(key, "/instances/controller") {
		ret := []string{}
		r, err := s.ControllerClient.List(meta_v1.ListOptions{})
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


	if key == "/longhorn_manager_test/settings" {
		ret := []string{}
		r, err := s.SettingClient.List(meta_v1.ListOptions{})
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
	panic(key)
	return nil, nil
}
