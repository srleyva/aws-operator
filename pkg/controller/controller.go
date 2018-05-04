// Package main for a s3 Bucket operator
package controller

import (
	"github.com/Sirupsen/logrus"
	opkit "github.com/rook/operator-kit"
	awsResources "github.com/srleyva/aws-operator/pkg/apis/sleyva/v1alpha1"
	leyvaclient "github.com/srleyva/aws-operator/pkg/client/clientset/versioned/typed/sleyva/v1alpha1"
	"github.com/srleyva/aws-operator/pkg/logger"
	"k8s.io/client-go/tools/cache"
	"os"
)

var s3Client *S3
var cfnClient *CFN
var err error

// LeyvaController represents a controller object for sample custom resources
type LeyvaController struct {
	context        *opkit.Context
	leyvaClientset leyvaclient.SleyvaV1alpha1Interface
}

// NewLeyvaController create controller for watching awsResources custom resources created
func NewLeyvaController(context *opkit.Context, leyvaClientset leyvaclient.SleyvaV1alpha1Interface) *LeyvaController {
	logrus.Info("Initializing Client for AWS")

	if s3Client, err = NewS3Client(); err != nil {
		logger.LogAWSErrorf("Client", "Error initializing client: %+v", err)
		os.Exit(1)
	}

	if cfnClient, err = NewCFNClient(); err != nil {
		logger.LogAWSErrorf("Client", "Error initializing client: %+v", err)
		os.Exit(1)
	}

	return &LeyvaController{
		context:        context,
		leyvaClientset: leyvaClientset,
	}
}

// Watch watches for instances of AWS custom resources and acts on them
func (c *LeyvaController) StartWatch(namespace string, stopCh chan struct{}) error {

	s3resourceHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAddS3,
		UpdateFunc: c.onUpdateS3,
		DeleteFunc: c.onDeleteS3,
	}
	cfnresourceHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAddCfn,
		UpdateFunc: c.onUpdateCfn,
		DeleteFunc: c.onDeleteCfn,
	}

	restClient := c.leyvaClientset.RESTClient()
	s3Watcher := opkit.NewWatcher(awsResources.S3BucketResource, namespace, s3resourceHandlers, restClient)
	cfnWatcher := opkit.NewWatcher(awsResources.CloudformationResource, namespace, cfnresourceHandlers, restClient)
	go s3Watcher.Watch(&awsResources.S3Bucket{}, stopCh)
	go cfnWatcher.Watch(&awsResources.Cloudformation{}, stopCh)
	return nil
}

func (c *LeyvaController) onAddS3(obj interface{}) {
	s := obj.(*awsResources.S3Bucket).DeepCopy()
	s3Client.CreateS3Bucket(*s)
	s3Client.SetBucketPolicy(s.Name, s.Spec.Policy)
}

func (c *LeyvaController) onUpdateS3(oldObj, newObj interface{}) {
	newSample := newObj.(*awsResources.S3Bucket).DeepCopy()
	logger.LogAWSInfof("S3", "Updating Bucket: %s", newSample.Name)
	s3Client.SetBucketPolicy(newSample.Name, newSample.Spec.Policy)
}

func (c *LeyvaController) onDeleteS3(obj interface{}) {
	s := obj.(*awsResources.S3Bucket).DeepCopy()
	s3Client.DeleteS3Bucket(s.Name)
}

func (c *LeyvaController) onAddCfn(obj interface{}) {
	s := obj.(*awsResources.Cloudformation).DeepCopy()
	logger.LogAWSDebugf("Template: %s", s.Spec.Template)
	cfnClient.CreateCfnStack(s)

}

func (c *LeyvaController) onUpdateCfn(oldObj, newObj interface{}) {
	s := newObj.(*awsResources.Cloudformation).DeepCopy()
	logger.LogAWSInfof("S3", "Updating CFN Stack: %s", s.Name)
}

func (c *LeyvaController) onDeleteCfn(obj interface{}) {
	s := obj.(*awsResources.Cloudformation).DeepCopy()
	cfnClient.DeleteCfnStack(s)
}
