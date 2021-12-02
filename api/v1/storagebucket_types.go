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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StorageBucketSpec defines the desired state of StorageBucket
type StorageBucketSpec struct {
	// The datacenter region of the storage bucket.
	Region string `json:"region,omitempty"`
	// The bucket name.
	Bucket string `json:"bucket,omitempty"`
	// All objects stored in the bucket will have keys those start with this prefix. It is assumed to set a string that
	// identifies the deployment environment of Phalanks. Default value is "phalanks".
	// +optional
	KeyPrefix string `json:"keyPrefix,omitempty"`
	// Secret refers to a secret resource that contains credentials for accessing a storage bucket.
	Secret StorageBucketSecret `json:"secret,omitempty"`
}

// StorageBucketSecret defines a pointer to a secret resource and how to interpret the information in the secret.
type StorageBucketSecret struct {
	// SecretRef refers to a secret resource that contains credentials for accessing a storage bucket.
	SecretName string `json:"secretName,omitempty"`
	// SecretKey
	// +optional
	SecretKey SecretKey `json:"secretKey,omitempty"`
}

// SecretKey describes the mapping between the key in the secret resource and the respective credentials for accessing
// the storage bucket.
type SecretKey struct {
	// AccessKeyId is the key name of the secret that will be used as AWS_ACCESS_KEY_ID.
	// +optional
	AccessKeyId string `json:"accessKeyId,omitempty"`
	// SecretAccessKey is the key name of the secret that will be used as AWS_SECRET_ACCESS_KEY.
	// +optional
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

// StorageBucketStatus defines the observed state of a StorageBucket.
type StorageBucketStatus struct {
	// Available indicates whether the storage bucket is available or not.
	Available bool `json:"available,omitempty"`
	// LastProbeTime is the time when the current state was observed.
	LastProbeTime *metav1.Time `json:"lastProbeTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// StorageBucket is the Schema for the StorageBucket API
// +kubebuilder:printcolumn:name="region",type=string,JSONPath=`.spec.region`
// +kubebuilder:printcolumn:name="bucket",type=string,JSONPath=`.spec.bucket`
// +kubebuilder:printcolumn:name="key_Prefix",type=string,JSONPath=`.spec.keyPrefix`
// +kubebuilder:printcolumn:name="available",type=string,JSONPath=`.status.available`
// +kubebuilder:printcolumn:name="last_probe_time",type=string,format=date,JSONPath=`.status.lastProbeTime`
type StorageBucket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StorageBucketSpec   `json:"spec,omitempty"`
	Status StorageBucketStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StorageBucketList contains a list of StorageBucket
type StorageBucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StorageBucket `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StorageBucket{}, &StorageBucketList{})
}
