// Package main for a s3 Bucket operator
package controller

import (
	"github.com/Sirupsen/logrus"
	opkit "github.com/rook/operator-kit"
	s3Bucket "github.com/srleyva/aws-operator/pkg/apis/sleyva/v1alpha1"
	leyvaclient "github.com/srleyva/aws-operator/pkg/client/clientset/versioned/typed/sleyva/v1alpha1"
	"github.com/srleyva/aws-operator/pkg/logger"
	"k8s.io/client-go/tools/cache"
	"os"
)

var s3Client *S3
var err error

// LeyvaController represents a controller object for sample custom resources
type LeyvaController struct {
	context        *opkit.Context
	leyvaClientset leyvaclient.SleyvaV1alpha1Interface
}

// NewLeyvaController create controller for watching s3Bucket custom resources created
func NewLeyvaController(context *opkit.Context, leyvaClientset leyvaclient.SleyvaV1alpha1Interface) *LeyvaController {
	logrus.Info("Initializing S3 Client for AWS")

	if s3Client, err = NewS3Client(); err != nil {
		logger.LogS3Errorf("Error initializing client: %+v", err)
		os.Exit(1)
	}
	return &LeyvaController{
		context:        context,
		leyvaClientset: leyvaClientset,
	}
}

// Watch watches for instances of AWS custom resources and acts on them
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
	s3Client.CreateS3Bucket(*s)
	s3Client.SetBucketPolicy(s.Name, s.Spec.Policy)
}

func (c *LeyvaController) onUpdate(oldObj, newObj interface{}) {
	newSample := newObj.(*s3Bucket.S3Bucket).DeepCopy()
	logger.LogS3Infof("Updating Bucket: %s", newSample.Name)
	s3Client.SetBucketPolicy(newSample.Name, newSample.Spec.Policy)
}

func (c *LeyvaController) onDelete(obj interface{}) {
	s := obj.(*s3Bucket.S3Bucket).DeepCopy()
	s3Client.DeleteS3Bucket(s.Name)
}
