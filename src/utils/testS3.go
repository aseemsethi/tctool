package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"fmt"
)

func TestS3() {
	fmt.Println("Test S3")

	// Init session in us-east-2
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	// Create S3 service client
	s3Svc := s3.New(sess)
	result, err := s3Svc.ListBuckets(nil)
	if err != nil {
		fmt.Println("got error listing buckets: ", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}
