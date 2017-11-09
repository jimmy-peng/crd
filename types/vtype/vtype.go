
package vtype

import (
meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/jimmy-peng/crd/tools/crdcopy"
	"github.com/rancher/longhorn-manager/types"
)

// Definition of our CRD LongHorn Volume class
type Crdvolume struct {
meta_v1.TypeMeta   `json:",inline"`
meta_v1.ObjectMeta `json:"metadata"`
Spec               CrdVolumeSpec   `json:"spec"`
Status             CrdVolumeStatus `json:"status,omitempty"`
}



type CrdVolumeSpec struct {
	// Attributes
	Name                string 	`json:"name"`
	Size                int64  	`json:",string"`
	BaseImage           string 	`json:"baseimage,omitempty"`
	FromBackup          string 	`json:"frombackup,omitempty"`
	NumberOfReplicas    int    	`json:"numreplicas,omitempty"`
	StaleReplicaTimeout int	   	`json:"stalereplicatimeout,omitempty"`

	// Running spec
	TargetNodeID  string		`json:"targetnodeid,omitempty"`
	DesireState   types.VolumeState	`json:"desirestate,omitempty"`
	RecurringJobs []types.RecurringJob`json:"recurringjobs,omitempty"`

	// Running state
	Created  string				`json:"create,omitempty"`
	NodeID   string 			`json:"nodeid,omitempty"`
	State    types.VolumeState	`json:"state,omitempty"`
	Endpoint string				`json:"endpoint,omitempty"`

	types.KVMetadata
}

type CrdVolumeStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

type CrdvolumeList struct {
	meta_v1.TypeMeta             `json:",inline"`
	meta_v1.ListMeta             `json:"metadata"`
	Items            []Crdvolume `json:"items"`
}

func LhVoulme2CRDVolume(vinfo *types.VolumeInfo, crdvolume *Crdvolume, pathname string){
	crdvolume.ObjectMeta.Name = vinfo.Name
	crdcopy.CRDDeepCopy(&crdvolume.Spec, vinfo)
}

func CRDVolume2LhVoulme(crdvolume *Crdvolume, vinfo *types.VolumeInfo)  {
	crdcopy.CRDDeepCopy(vinfo, &crdvolume.Spec)
}
/*
type KVMetadata struct {
	KVIndex uint64 `json:"-"`
}

type RecurringJob struct {
	Name   string           `json:"name"`
	Type   RecurringJobType `json:"task"`
	Cron   string           `json:"cron"`
	Retain int              `json:"retain"`
}

type VolumeInfo struct {
	// Attributes
	Name                string
	Size                int64 `json:",string"`
	BaseImage           string
	FromBackup          string
	NumberOfReplicas    int
	StaleReplicaTimeout int

	// Running spec
	TargetNodeID  string
	DesireState   VolumeState
	RecurringJobs []RecurringJob

	// Running state
	Created  string
	NodeID   string
	State    VolumeState
	Endpoint string

	KVMetadata
}
type RecurringJobType string

type VolumeState string

*/