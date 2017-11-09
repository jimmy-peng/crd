
package ntype

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Definition of our CRD LongHorn Volume class
type Crdnode struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               CrdNodeSpec   `json:"spec"`
	Status             CrdNodeStatus `json:"status,omitempty"`
}

type CrdNodeSpec struct {
	Name string `json:"name"`
	NodeId string   `json:"nodeid"`
	NumReplicas int    `json:"numreplicas,omitempty"`
}

type CrdNodeStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

type CrdnodeList struct {
	meta_v1.TypeMeta             `json:",inline"`
	meta_v1.ListMeta             `json:"metadata"`
	Items            []Crdnode `json:"items"`
}

func LhNode2CRDNode(vinfo *NodeInfo, crdnode *Crdnode, pathname string){
	crdnode.ObjectMeta.Name = vinfo.Name
	crdnode.Spec.Name = pathname
	crdnode.Spec.NodeId = vinfo.NodeID
	crdnode.Spec.NumReplicas = vinfo.NumberOfReplicas
}

func CRDNode2LhNode(crdvolume *Crdnode, vinfo *NodeInfo)  {

	vinfo.Name = crdvolume.ObjectMeta.Name
	vinfo.NodeID = crdvolume.Spec.NodeId
	vinfo.NumberOfReplicas = crdvolume.Spec.NumReplicas
}

type NodeInfo struct {
	// Attributes
	Name                string
	NumberOfReplicas    int
	NodeID   string

}

