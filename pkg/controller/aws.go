package controller

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	awsResources "github.com/srleyva/aws-operator/pkg/apis/sleyva/v1alpha1"
	"github.com/srleyva/aws-operator/pkg/logger"
	"os"
)

type S3 struct {
	Client s3iface.S3API
}

type CFN struct {
	Client cloudformationiface.CloudFormationAPI
}

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

func (s *S3) CreateS3Bucket(bucket awsResources.S3Bucket) (err error) {
	logger.LogAWSInfof("S3", "Creating S3 Bucket: %s", bucket.Name)
	_, err = s.Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket.Name),
	})

	err = s.Client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket.Name),
	})
	if err != nil {
		logger.LogAWSErrorf("S3", "Error Creating Bucket: %+v", err)
		return err
	}

	logger.LogAWSInfof("S3", "Bucket %s created successfully", bucket.Name)
	return nil

}

func (s *S3) SetBucketPolicy(bucketName, policy string) (err error) {
	logger.LogAWSInfof("S3", "Updating Bucket Policy: %s", bucketName)
	_, err = s.Client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(policy),
	})

	if err != nil {
		logger.LogAWSErrorf("S3", "Error updating bucket policy: %+v", err)
		return err
	}

	logger.LogAWSInfof("S3", "Policy for %s updated successfully", bucketName)
	return err
}

func (s *S3) DeleteS3Bucket(bucketName string) (err error) {
	logger.LogAWSInfof("S3", "Deleting bucket: %s", bucketName)
	_, err = s.Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		logger.LogAWSErrorf("S3", "Error Deleting Bucket %+v", err)
		return err
	}

	logger.LogAWSInfof("S3", "Bucket %s deleted successfully", bucketName)
	return nil
}

func NewCFNClient() (*CFN, error) {
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

	client := CFN{cloudformation.New(sess)}

	return &client, nil
}

func (s *CFN) CreateCfnStack(cfn *awsResources.Cloudformation) (err error) {
	logger.LogAWSInfof("S3", "Cloudformation", "Create New Stack: %s", cfn.Name)
	_, err = s.Client.CreateStack(&cloudformation.CreateStackInput{
		StackName:    aws.String(cfn.Name),
		TemplateBody: aws.String(cfn.Spec.Template),
	})

	if err != nil {
		logger.LogAWSErrorf("S3", "Error Creating Stack: %+v", err)
		return err
	}
	// TODO Return status
	logger.LogAWSInfof("S3", "CFN Stack %s create initiated successfully. Please check status.", cfn.Name)
	return nil
}

func (s *CFN) DeleteCfnStack(cfn *awsResources.Cloudformation) (err error) {
	logger.LogAWSInfof("S3", "Cloudformation", "Deleting Stack: %s", cfn.Name)
	_, err = s.Client.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: aws.String(cfn.Name),
	})

	if err != nil {
		logger.LogAWSErrorf("S3", "Error Deleting Stack: %+v", err)
		return err
	}
	// TODO Watch status and return to spec status
	logger.LogAWSInfof("S3", "CFN Stack %s deleted successfully. Please check status.", cfn.Name)
	return nil
}
