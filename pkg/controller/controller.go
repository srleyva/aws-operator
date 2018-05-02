// Package main for a s3 Bucket operator
package controller

import (
	"fmt"
	opkit "github.com/rook/operator-kit"
	s3Bucket "github.com/srleyva/aws-operator/pkg/apis/sleyva/v1alpha1"
	leyvaclient "github.com/srleyva/aws-operator/pkg/client/clientset/versioned/typed/sleyva/v1alpha1"
	"k8s.io/client-go/tools/cache"
)

// SampleController represents a controller object for sample custom resources
type LeyvaController struct {
	context         *opkit.Context
	leyvaClientset leyvaclient.SleyvaV1alpha1Interface
}

// NewLeyvaController create controller for watching s3Bucket custom resources created
func NewLeyvaController(context *opkit.Context, leyvaClientset leyvaclient.SleyvaV1alpha1Interface) *LeyvaController {
	return &LeyvaController{
		context:         context,
		leyvaClientset: leyvaClientset,
	}
}

// Watch watches for instances of Sample custom resources and acts on them
func (c *LeyvaController) StartWatch(namespace string, stopCh chan struct{}) error {

	resourceHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	}
	restClient := c.leyvaClientset.RESTClient()
	watcher := opkit.NewWatcher(s3Bucket.S3BucketResource, namespace, resourceHandlers, restClient)
	go watcher.Watch(&s3Bucket.S3Bucket{}, stopCh)
	return nil
}

func (c *LeyvaController) onAdd(obj interface{}) {
	s := obj.(*s3Bucket.S3Bucket).DeepCopy()

	fmt.Printf("Create S3Bucket '%s' in region %s with bucketname %s \n", s.Name, s.Spec.Region, s.Spec.Name)
}

func (c *LeyvaController) onUpdate(oldObj, newObj interface{}) {
	oldSample := oldObj.(*s3Bucket.S3Bucket).DeepCopy()
	newSample := newObj.(*s3Bucket.S3Bucket).DeepCopy()

	fmt.Printf("Updated S3Bucket '%s' to '%s\n", newSample.Name, oldSample.Name)
}

func (c *LeyvaController) onDelete(obj interface{}) {
	s := obj.(*s3Bucket.S3Bucket).DeepCopy()

	fmt.Printf("Deleted S3 Bucket '%s'\n", s.Name)
}
