// Package v1alpha1 for a S3Bucket crd
package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type S3Bucket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              S3BucketSpec `json:"spec"`
}

type S3BucketSpec struct {
	Name   string `json:"name"`
	Policy string `json:"policy"`
}

//type policy struct {
//	Version   string `json:"Version"`
//	Statement []struct {
//		Effect   string   `json:"Effect"`
//		Action   []string `json:"Action"`
//		Resource string   `json:"Resource"`
//	} `json:"Statement"`
//}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type S3BucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []S3Bucket `json:"items"`
}
