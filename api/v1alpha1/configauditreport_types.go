/*
Copyright 2021.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConfigAuditReportSpec defines the desired state of ConfigAuditReport
type ConfigAuditReportSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ConfigAuditReport. Edit configauditreport_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// ConfigAuditReportStatus defines the observed state of ConfigAuditReport
type ConfigAuditReportStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ConfigAuditReport is the Schema for the configauditreports API
type ConfigAuditReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigAuditReportSpec   `json:"spec,omitempty"`
	Status ConfigAuditReportStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConfigAuditReportList contains a list of ConfigAuditReport
type ConfigAuditReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigAuditReport `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigAuditReport{}, &ConfigAuditReportList{})
}
