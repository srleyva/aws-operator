package controller

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	s3Bucket "github.com/srleyva/aws-operator/pkg/apis/sleyva/v1alpha1"
)

type S3 struct {
	Client s3iface.S3API
}

var (
	ErrNoSuchBucket  = errors.New("NoSuchBucket: The specified bucket does not exist")
	ErrBucketExists  = errors.New("Bucket already exists")
	ErrBucketHasKeys = errors.New("Bucket has keys so cannot be deleted")
)

func NewS3Client() (*S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")}, //TODO: Extract to config
	)

	if err != nil {
		return nil, err

	}

	client := S3{s3.New(sess)}

	return &client, nil
}

func (s *S3) CreateS3Bucket(bucket s3Bucket.S3Bucket) (err error) {

	_, err = s.Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket.Spec.Name),
	})

	if err != nil {
		return err
	}
	err = s.Client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket.Spec.Name),
	})
	if err != nil {
		return err
	}
	return nil

}
