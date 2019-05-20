package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IncomingPortalSpec defines the desired state of IncomingPortal
// +k8s:openapi-gen=true
type IncomingPortalSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	LocalPorts []int32 `json:"localPorts,omitempty"`
}

// IncomingPortalStatus defines the observed state of IncomingPortal
// +k8s:openapi-gen=true
type IncomingPortalStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IncomingPortal is the Schema for the incomingportals API
// +k8s:openapi-gen=true
type IncomingPortal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IncomingPortalSpec   `json:"spec,omitempty"`
	Status IncomingPortalStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IncomingPortalList contains a list of IncomingPortal
type IncomingPortalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IncomingPortal `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IncomingPortal{}, &IncomingPortalList{})
}
