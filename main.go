// Package main for a aws-operator
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	opkit "github.com/rook/operator-kit"
	s3Bucket "github.com/srleyva/aws-operator/pkg/apis/sleyva/v1alpha1"
	leyvaClient "github.com/srleyva/aws-operator/pkg/client/clientset/versioned/typed/sleyva/v1alpha1"
	"github.com/srleyva/aws-operator/pkg/controller"
	"k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	fmt.Println("Getting kubernetes context")
	context, leyvaClientset, err := createContext()
	if err != nil {
		panic(fmt.Errorf("failed to create context. %+v\n", err))
	}

	// Create and wait for CRD resources
	fmt.Println("Registering the sample resource")
	resources := []opkit.CustomResource{s3Bucket.S3BucketResource}
	err = opkit.CreateCustomResources(*context, resources)
	if err != nil {
		panic(fmt.Errorf("failed to create CRDs. %+v\n", err))
	}

	// create signals to stop watching the resources
	signalChan := make(chan os.Signal, 1)
	stopChan := make(chan struct{})
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// start watching the sample resource
	fmt.Println("Watching the s3Bucket resource")
	controller := controller.NewLeyvaController(context, leyvaClientset)
	controller.StartWatch(v1.NamespaceAll, stopChan)

	for {
		select {
		case <-signalChan:
			fmt.Println("shutdown signal received, exiting...")
			close(stopChan)
			return
		}
	}
}

func createContext() (*opkit.Context, leyvaClient.SleyvaV1alpha1Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get k8s config. %+v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get k8s client. %+v", err)
	}

	apiExtClientset, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create k8s API extension clientset. %+v", err)
	}

	sampleClientset, err := leyvaClient.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create sample clientset. %+v", err)
	}

	context := &opkit.Context{
		Clientset:             clientset,
		APIExtensionClientset: apiExtClientset,
		Interval:              500 * time.Millisecond,
		Timeout:               60 * time.Second,
	}
	return context, sampleClientset, nil

}
