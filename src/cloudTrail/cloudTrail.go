package cloudTrail

import (
	"fmt"
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/sirupsen/logrus"
)

type CloudTrail struct {
	Name string
	svc  *cloudtrail.CloudTrail
}

var cLog *logrus.Logger

func (i *CloudTrail) Initialize() bool {
	cLog = tcGlobals.Tcg.Log

	// Create a CloudTrail service client.
	svc := cloudtrail.New(tcGlobals.Tcg.Sess)
	i.svc = svc

	return true
}

/* AWS CloudTrail is now enabled by default for ALL CUSTOMERS and will provide visibility
 * into the past seven days of account activity without the need for you to configure a
 * trail in the service to get started
 * We thus check if any trail is configured.
 */
func checkIfEnabled(i *CloudTrail) {
	resp, err := i.svc.DescribeTrails(&cloudtrail.DescribeTrailsInput{TrailNameList: nil})
	if err != nil {
		fmt.Println("Got error calling CreateTrail:")
		fmt.Println(err.Error())
	}

	fmt.Println("Found", len(resp.TrailList), "trail(s) in", "us-west-2")
	if len(resp.TrailList) == 0 {
		cLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 2.1, "Result": "Failed",
		}).Info("CloudTrail is disabled")
	} else {
		for _, trail := range resp.TrailList {
			fmt.Println("Trail name:  " + *trail.Name)
			fmt.Println("Bucket name: " + *trail.S3BucketName)
			fmt.Println("")
		}
		cLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 2.1, "Result": "Passed",
		}).Info("CloudTrail is enabled")
	}
}

func (i *CloudTrail) Run() {
	checkIfEnabled(i)
}
