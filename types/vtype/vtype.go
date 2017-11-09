
package vtype

import (
meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Definition of our CRD LongHorn Volume class
type Crdvolume struct {
meta_v1.TypeMeta   `json:",inline"`
meta_v1.ObjectMeta `json:"metadata"`
Spec               CrdVolumeSpec   `json:"spec"`
Status             CrdVolumeStatus `json:"status,omitempty"`
}

type CrdVolumeSpec struct {
	Name string `json:"name"`
	NodeId string   `json:"nodeid"`
	NumReplicas int    `json:"numreplicas,omitempty"`
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

func LhVoulme2CRDVolume(vinfo *VolumeInfo, crdvolume *Crdvolume, pathname string){
	crdvolume.ObjectMeta.Name = vinfo.Name
	crdvolume.Spec.Name = pathname
	crdvolume.Spec.NodeId = vinfo.NodeID
	crdvolume.Spec.NumReplicas = vinfo.NumberOfReplicas
}

func CRDVolume2LhVoulme(crdvolume *Crdvolume, vinfo *VolumeInfo)  {

	vinfo.Name = crdvolume.ObjectMeta.Name
	vinfo.NodeID = crdvolume.Spec.NodeId
	vinfo.NumberOfReplicas = crdvolume.Spec.NumReplicas
}

type VolumeInfo struct {
	// Attributes
	Name                string
	NumberOfReplicas    int
	NodeID   string

}

