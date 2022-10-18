/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ReplicatedSecretSpec defines the desired state of ReplicatedSecret
type ReplicatedSecretSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// SecretData is a field of ReplicatedSecret contains all data for replicated secrets
	SecretData map[string][]byte `json:"secret_data,omitempty"`
	// Namespaces are the k8s namespaces in which the secrets should be created.
	Namespaces []string `json:"namespaces,omitempty"`
	// Type is kubernetes secret type https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
	Type string `json:"type,omitempty"`
	// NamespacesPrefix is prefix for namespaces in which the secrets should be created.
	NamespacesPrefix string `json:"namespaces_prefix,omitempty"`
	// NamespacesRegex is regex for namespaces in which the secrets should be created.
	NamespacesRegex string `json:"namespaces_regex,omitempty"`
}

// ReplicatedSecretStatus defines the observed state of ReplicatedSecret
type ReplicatedSecretStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// ReplicatedSecret is the Schema for the replicatedsecrets API
type ReplicatedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReplicatedSecretSpec   `json:"spec,omitempty"`
	Status ReplicatedSecretStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ReplicatedSecretList contains a list of ReplicatedSecret
type ReplicatedSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReplicatedSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ReplicatedSecret{}, &ReplicatedSecretList{})
}
