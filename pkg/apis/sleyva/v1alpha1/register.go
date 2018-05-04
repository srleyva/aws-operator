// Package v1alpha1 for a S3Bucket crd
package v1alpha1

import (
	"reflect"

	opkit "github.com/rook/operator-kit"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = SchemeBuilder.AddToScheme
)

// schemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: "sleyva.io", Version: "v1alpha1"}

var S3BucketResource = opkit.CustomResource{
	Name:       "s3bucket",
	Plural:     "s3buckets",
	Group:      "sleyva.io",
	Version:    "v1alpha1",
	Scope:      apiextensionsv1beta1.NamespaceScoped,
	Kind:       reflect.TypeOf(S3Bucket{}).Name(),
	ShortNames: []string{"s3", "buckets"},
}

var CloudformationResource = opkit.CustomResource{
	Name:       "cloudformation",
	Plural:     "cloudformation",
	Group:      "sleyva.io",
	Version:    "v1alpha1",
	Scope:      apiextensionsv1beta1.NamespaceScoped,
	Kind:       reflect.TypeOf(Cloudformation{}).Name(),
	ShortNames: []string{"cfn", "stacks"},
}

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(addKnownTypes)
}

// Resource takes an unqualified resource and returns back a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&S3Bucket{},
		&S3BucketList{},
	)
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Cloudformation{},
		&CloudformationList{})

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
