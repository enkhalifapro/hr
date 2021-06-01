package awshelper

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"strings"
)

// S3ListBucketsAPI defines the interface for the ListBuckets function.
// We use this interface to test the function using a mocked service.
type S3API interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type S3 struct {
	client *s3.Client
}

func NewS3(client *s3.Client) *S3 {
	return &S3{
		client: client,
	}
}

// GetAllBuckets retrieves a list of your Amazon Simple Storage Service (Amazon S3) buckets.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a ListBucketsOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ListBuckets.
func GetAllBuckets(c context.Context, api S3API, input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return api.ListBuckets(c, input)
}

// MakeBucket creates an Amazon Simple Storage Service (Amazon S3) bucket.
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a CreateBucketOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to CreateBucket.
func MakeBucket(c context.Context, api S3API, input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return api.CreateBucket(c, input)
}

// PutFile uploads a file to an Amazon Simple Storage Service (Amazon S3) bucket
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a PutObjectOutput object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to PutObject
func PutFile(c context.Context, api S3API, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}
func (s *S3) ListBuckets() error {

	// Create an Amazon S3 service client
	input := &s3.ListBucketsInput{}

	result, err := GetAllBuckets(context.TODO(), s.client, input)
	if err != nil {
		fmt.Println("Got an error retrieving buckets:")
		fmt.Println(err)
		return err
	}

	fmt.Println("Buckets:")

	for _, bucket := range result.Buckets {
		fmt.Println(*bucket.Name + ": " + bucket.CreationDate.Format("2006-01-02 15:04:05 Monday"))
	}
	return nil
}
func (s *S3) Create(bucket string) error {

	if strings.TrimSpace(bucket) == "" {
		return errors.New("You must supply a bucket name ")
	}

	input := &s3.CreateBucketInput{
		Bucket: &bucket,
	}
	_, err := MakeBucket(context.TODO(), s.client, input)
	if err != nil {
		if strings.Contains(err.Error(),"BucketAlreadyExists")  {
			return errors.New("Could not create bucket (" + bucket + ") , Bucket Already Exists")
		} else {
			fmt.Println("Could not create bucket (" + bucket + ") , Bucket Already Exists")
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

func (s *S3) Upload(bucket string,filename string,data string) {

	if bucket == "" || filename == "" {
		fmt.Println("You must supply a bucket name   and file name ")
		return
	}



	// create a reader from data data in memory
	reader := strings.NewReader(data)

	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &filename,
		Body:   reader,
	}

	_, err := PutFile(context.TODO(), s.client, input)
	if err != nil {
		fmt.Println("Got error uploading file:")
		fmt.Println(err)
		return
	}

}