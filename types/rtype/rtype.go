package rtype

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/jimmy-peng/crd/tools/crdcopy"
	"github.com/rancher/longhorn-manager/types"
)

// Definition of our CRD LongHorn Volume class
type Crdreplicas struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               CrdReplicasSpec   `json:"spec"`
	Status             CrdReplicasStatus `json:"status,omitempty"`
}



type CrdReplicasSpec struct {
	ID         string
	Type       types.InstanceType
	Name       string
	NodeID     string
	IP         string
	Running    bool
	VolumeName string

	types.KVMetadata
	FailedAt string
}

type CrdReplicasStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

type CrdreplicasList struct {
	meta_v1.TypeMeta             `json:",inline"`
	meta_v1.ListMeta             `json:"metadata"`
	Items            []Crdreplicas `json:"items"`
}

func LhReplicas2CRDReplicas(rinfo *types.ReplicaInfo, crdreplica *Crdreplicas, key string) {
	crdreplica.ObjectMeta.Name = key
	crdcopy.CRDDeepCopy(&crdreplica.Spec, rinfo)
}

func CRDReplicas2LhReplicas(crdreplicas *Crdreplicas, rinfo *types.ReplicaInfo) {
	crdcopy.CRDDeepCopy(rinfo, &crdreplicas.Spec)
}