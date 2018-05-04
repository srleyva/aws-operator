package controller

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	awsResources "github.com/srleyva/aws-operator/pkg/apis/sleyva/v1alpha1"
	"github.com/srleyva/aws-operator/pkg/logger"
	"reflect"
	"sync"
	"testing"
)

type MockBucket map[string][]byte

var buffer bytes.Buffer

type MockS3 struct {
	s3iface.S3API
	sync.RWMutex
	// bucket: {key: value}
	data map[string]MockBucket
}

type MockCFN struct {
	cloudformationiface.CloudFormationAPI
}

func NewMockS3() *MockS3 {
	logger.NewLogger(&logrus.TextFormatter{}, logrus.DebugLevel, &buffer)
	return &MockS3{
		data: map[string]MockBucket{},
	}
}

func (self *MockS3) CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	self.Lock()
	defer self.Unlock()
	if _, exists := self.data[*input.Bucket]; exists {
		return nil, ErrBucketExists
	}
	self.data[*input.Bucket] = MockBucket{}
	return &s3.CreateBucketOutput{}, nil
}

func (self *MockS3) WaitUntilBucketExists(input *s3.HeadBucketInput) error {
	return nil
}

func (self *MockS3) PutBucketPolicy(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
	resp := s3.PutBucketPolicyOutput{}
	return &resp, nil
}

func (self *MockS3) DeleteBucket(input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	resp := &s3.DeleteBucketOutput{}
	return resp, nil
}

func TestNewS3Client(t *testing.T) {
	config, err := NewS3Client()
	if err != nil {
		t.Errorf("error returned when not expected: %s", err)
	}
	if reflect.TypeOf(config) != reflect.TypeOf(&S3{}) {
		t.Errorf("Expected: %s \nActual: %s", reflect.TypeOf(config), reflect.TypeOf(S3{}))
	}
}

func TestS3_CreateS3Bucket(t *testing.T) {
	// TODO Seperate cases
	bucket := awsResources.S3Bucket{}
	bucket.Name = "my-test-bucket"
	s3client := S3{Client: NewMockS3()}
	err := s3client.CreateS3Bucket(bucket)
	if err != nil {
		t.Errorf("error returned when not expected: %s", err)
	}
}

func TestSetS3BucketPolicy(t *testing.T) {
	policy := `{
       "Version":"2012-10-17",
       "Statement":[
          {
             "Effect":"Allow",
             "Action":[
                "s3:ListAllMyBuckets"
             ],
             "Resource":"arn:aws:s3:::*"
          },
          {
             "Effect":"Allow",
             "Action":[
                "s3:ListBucket",
                "s3:GetBucketLocation"
             ],
             "Resource":"arn:aws:s3:::examplebucket"
          },
          {
             "Effect":"Allow",
             "Action":[
                "s3:PutObject",
                "s3:PutObjectAcl",
                "s3:GetObject",
                "s3:GetObjectAcl",
                "s3:DeleteObject"
             ],
             "Resource":"arn:aws:s3:::examplebucket/*"
          }
	}`

	s3client := S3{Client: NewMockS3()}
	err := s3client.SetBucketPolicy("test-bucket", policy)
	if err != nil {
		t.Errorf("error returned when not expected: %s", err)
	}

}

func TestDeleteS3Bucket(t *testing.T) {
	s3client := S3{Client: NewMockS3()}
	err := s3client.DeleteS3Bucket("test-bucket")
	if err != nil {
		t.Errorf("error returned when not expected: %s", err)
	}
}

func NewMockCFN() *MockCFN {
	logger.NewLogger(&logrus.TextFormatter{}, logrus.DebugLevel, &buffer)
	return &MockCFN{}
}

func (self *MockCFN) CreateStack(input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	resp := cloudformation.CreateStackOutput{
		StackId: aws.String("1234"),
	}
	return &resp, nil
}

func (self *MockCFN) DeleteStack(input *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error) {
	resp := cloudformation.DeleteStackOutput{}
	return &resp, nil
}

func TestCreateCfnStack(t *testing.T) {
	stack := `{"Resources" : { "dontkickthebucket" : { "Type" : "AWS::S3::Bucket" }}}`

	fakeCfn := awsResources.Cloudformation{}
	fakeCfn.Name = "test-stack"
	fakeCfn.Spec.Template = stack

	cfnClient := CFN{Client: NewMockCFN()}
	err := cfnClient.CreateCfnStack(&fakeCfn)
	if err != nil {
		t.Errorf("error returned when not expected: %s", err)
	}

}

func TestDeleteCfnStack(t *testing.T) {
	fakeCfn := awsResources.Cloudformation{}
	fakeCfn.Name = "test-stack"

	cfnClient := CFN{Client: NewMockCFN()}
	err := cfnClient.DeleteCfnStack(&fakeCfn)
	if err != nil {
		t.Errorf("error returned when not expected: %s", err)
	}

}
