package controller

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	s3Bucket "github.com/srleyva/aws-operator/pkg/apis/sleyva/v1alpha1"
	"os"
	"fmt"
)

type S3 struct {
	Client s3iface.S3API
}

//type Policy struct {
//	Version   string `json:"Version"`
//	Statement []struct {
//		Effect   string   `json:"Effect"`
//		Action   []string `json:"Action"`
//		Resource string   `json:"Resource"`
//	} `json:"Statement"`
//}

var (
	ErrNoSuchBucket  = errors.New("NoSuchBucket: The specified bucket does not exist")
	ErrBucketExists  = errors.New("Bucket already exists")
	ErrBucketHasKeys = errors.New("Bucket has keys so cannot be deleted")
)

func NewS3Client() (*S3, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		return nil, err
	}

	client := S3{s3.New(sess)}

	return &client, nil
}

func (s *S3) CreateS3Bucket(bucket s3Bucket.S3Bucket) (err error) {

	_, err = s.Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket.Name),
	})

	if err != nil {
		return err
	}
	err = s.Client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket.Name),
	})
	if err != nil {
		return err
	}

	return nil

}

func (s *S3) SetBucketPolicy(bucketName, policy string) (err error) {
	_, err = s.Client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(policy),
	})
	fmt.Errorf("err: %s \n", err)
	return err
}
