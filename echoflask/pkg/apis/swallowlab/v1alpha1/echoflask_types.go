package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EchoFlaskSpec defines the desired state of EchoFlask
type EchoFlaskSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    // Size is the size of the FlaskEcho deployment
        Size int32 `json:"size"`

}

// EchoFlaskStatus defines the observed state of EchoFlask
type EchoFlaskStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    // Nodes are the names of the FlaskEcho pods
        Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EchoFlask is the Schema for the echoflasks API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=echoflasks,scope=Namespaced
type EchoFlask struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EchoFlaskSpec   `json:"spec,omitempty"`
	Status EchoFlaskStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EchoFlaskList contains a list of EchoFlask
type EchoFlaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EchoFlask `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EchoFlask{}, &EchoFlaskList{})
}
