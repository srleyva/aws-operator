package controller

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	s3Bucket "github.com/srleyva/aws-operator/pkg/apis/sleyva/v1alpha1"
	"reflect"
	"sync"
	"testing"
)

type MockBucket map[string][]byte

type MockS3 struct {
	s3iface.S3API
	sync.RWMutex
	// bucket: {key: value}
	data map[string]MockBucket
}

func NewMockS3() *MockS3 {
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

func TestS3_CreateS3Bucket(t *testing.T) {
	// TODO Seperate cases
	bucket := s3Bucket.S3Bucket{Spec: s3Bucket.S3BucketSpec{Name: "hello"}}
	s3client := S3{Client: NewMockS3()}
	err := s3client.CreateS3Bucket(bucket)
	if err != nil {
		t.Errorf("error returned when not expected: %s", err)
	}

	err = s3client.CreateS3Bucket(bucket)
	if err != ErrBucketExists {
		t.Errorf("Error not returned as expected: %s", err)
	}

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
