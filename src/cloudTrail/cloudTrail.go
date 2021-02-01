package cloudTrail

import (
	"bytes"
	"encoding/json"
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
)

type CloudTrail struct {
	Name  string
	svc   *cloudtrail.CloudTrail
	s3Svc *s3.S3
}

var cLog *logrus.Logger

func (i *CloudTrail) Initialize() bool {
	cLog = tcGlobals.Tcg.Log

	// Create a CloudTrail service client.
	svc := cloudtrail.New(tcGlobals.Tcg.Sess)
	i.svc = svc

	// Create S3 service client
	i.s3Svc = s3.New(tcGlobals.Tcg.Sess)
	//i.s3Svc = s3.New(tcGlobals.Tcg.Sess, aws.NewConfig().WithRegion("us-east-1"))

	return true
}

func checkS3(i *CloudTrail, bucketName *string) {
	// Get the bucket name configured for CloudTrail
	//fmt.Println("Search Bucket: ", *bucketName)
	cLog.WithFields(logrus.Fields{
		"Test": "CIS", "Num": 2.3,
	}).Info("Search S3 Bucket for CloudTrail: ", *bucketName)
	in := &s3.HeadBucketInput{
		Bucket: aws.String(*bucketName),
	}
	_, err := i.s3Svc.HeadBucket(in)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			cLog.WithFields(logrus.Fields{
				"Test": "CIS", "Num": 2.3, "Result": "Failed",
			}).Info("S3 Bucket not found: ", aerr.Code())
		} else {
			cLog.WithFields(logrus.Fields{
				"Test": "CIS", "Num": 2.3, "Result": "Failed",
			}).Info("S3 Bucket not found..: ", err.Error())
		}
		return
	}
	cLog.WithFields(logrus.Fields{
		"Test": "CIS", "Num": 2.3, "Result": "Passed",
	}).Info("S3 Bucket found..: ")

	// Ensure the policy does not contain a Statement having an Effect set to
	// Allow and a Principal set to "*" or {"AWS" : "*"}
	// Call S3 to retrieve the JSON formatted policy for the selected bucket.
	result, err := i.s3Svc.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(*bucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			cLog.WithFields(logrus.Fields{
				"Test": "CIS", "Num": 2.3, "Result": "Failed",
			}).Info("S3 Bucket Policy not found: ", aerr.Code())
		} else {
			cLog.WithFields(logrus.Fields{
				"Test": "CIS", "Num": 2.3, "Result": "Failed",
			}).Info("S3 Bucket Policy not found..: ", err.Error())
		}
		return
	}

	out := bytes.Buffer{}
	policyStr := aws.StringValue(result.Policy)
	if err := json.Indent(&out, []byte(policyStr), "", "  "); err != nil {
		cLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 2.3, "Result": "Failed",
		}).Info("Failed to pretty the S3 Policy: ", err)
	}
	//fmt.Printf("Bucket Policy:\n")
	//fmt.Println(out.String())
	cLog.WithFields(logrus.Fields{
		"Test": "CIS", "Num": 2.3, "Result": "Failed",
	}).Info("S3 Bucket Policy: ", out.String())
	allow := tcGlobals.CheckPolicyForAllowAll(result.Policy)
	if allow == true {
		cLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 2.3, "Result": "Failed",
		}).Info("S3 Policy allows Public access: ", err)
	} else {
		cLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 2.3, "Result": "Passed",
		}).Info("S3 Policy does not allows Public access: ")
	}
}

/* AWS CloudTrail is now enabled by default for ALL CUSTOMERS and will provide visibility
 * into the past seven days of account activity without the need for you to configure a
 * trail in the service to get started
 * We thus check if any trail is configured.
 */
func checkIfEnabled(i *CloudTrail) {
	resp, err := i.svc.DescribeTrails(&cloudtrail.DescribeTrailsInput{TrailNameList: nil})
	if err != nil {
		cLog.WithFields(logrus.Fields{
			"Test": "CIS"}).Info("Error getting trail: ", err.Error())
	}

	cLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Found trail len: ", len(resp.TrailList))
	if len(resp.TrailList) == 0 {
		cLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 2.1, "Result": "Failed",
		}).Info("CloudTrail is disabled")
	} else {
		cLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 2.1, "Result": "Passed",
		}).Info("CloudTrail is enabled")
		for _, trail := range resp.TrailList {
			cLog.WithFields(logrus.Fields{
				"Test": "CIS"}).Info("Found Trail: ", *trail.Name, " Bucket: ", *trail.S3BucketName)
			if trail.LogFileValidationEnabled == nil || *trail.LogFileValidationEnabled == false {
				cLog.WithFields(logrus.Fields{
					"Test": "CIS", "Num": 2.2, "Result": "Failed",
				}).Info("CloudTrail LogFileValidationEnabled is disabled for trail: ", *trail.Name)
			} else {
				cLog.WithFields(logrus.Fields{
					"Test": "CIS", "Num": 2.2, "Result": "Passed",
				}).Info("CloudTrail LogFileValidationEnabled is enabled for trail: ", *trail.Name)
			}
			checkS3(i, trail.S3BucketName)
		}
	}
}

func (i *CloudTrail) Run() {
	cLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("CloudTrail Run...")
	checkIfEnabled(i)
}
