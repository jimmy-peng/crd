
package vtype

import (
meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Definition of our CRD LongHorn Volume class
type Lhvolume struct {
meta_v1.TypeMeta   `json:",inline"`
meta_v1.ObjectMeta `json:"metadata"`
Spec               LhVolumeSpec   `json:"spec"`
Status             LhVolumeStatus `json:"status,omitempty"`
}

type LhVolumeSpec struct {
Name string `json:"name"`
NodeId string   `json:"nodeid"`
NumReplicas int    `json:"numreplicas,omitempty"`
}

type LhVolumeStatus struct {
State   string `json:"state,omitempty"`
Message string `json:"message,omitempty"`
}

type LhvolumeList struct {
meta_v1.TypeMeta            `json:",inline"`
meta_v1.ListMeta            `json:"metadata"`
Items            []Lhvolume `json:"items"`
}