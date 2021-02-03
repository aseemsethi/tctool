package cloudWatch

import (
	"fmt"
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/sirupsen/logrus"
	"strings"
)

type CloudWatch struct {
	Name      string
	svc       *cloudwatchlogs.CloudWatchLogs
	LogGroups *cloudwatchlogs.DescribeLogGroupsOutput
}

var cLog *logrus.Logger

func (i *CloudWatch) Initialize() bool {
	fmt.Println("CloudWatch Init...")
	cLog = tcGlobals.Tcg.Log
	svc := cloudwatchlogs.New(tcGlobals.Tcg.Sess)
	i.svc = svc

	return true
}

func lookupCloudWatchLogMetricFilter(i *CloudWatch, name, logGroupName string, nextToken *string) {
	input := cloudwatchlogs.DescribeMetricFiltersInput{
		//FilterNamePrefix: aws.String(name),
		LogGroupName: aws.String(logGroupName),
		NextToken:    nextToken,
	}
	//fmt.Printf("Reading CloudWatch Log Metric Filter: %s", input)
	resp, err := i.svc.DescribeMetricFilters(&input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "ResourceNotFoundException" {
			cLog.WithFields(logrus.Fields{
				"Test": "CIS", "Num": 3.3,
			}).Info("CloudWatch Log Metric Filters not retrieved - ResourceNotFoundException: ", err)
			return
		}
		cLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 3.3,
		}).Info("CloudWatch Log Metric Filters not retrieved: ", err)
		return
	}
	for _, mf := range resp.MetricFilters {
		//fmt.Println("\nFilterName: ", mf)
		if strings.Contains(*mf.FilterPattern, "$.userIdentity.type = \"Root\"") {
			//fmt.Println("CloudWatch Log Metric Filter checking for Root found:", mf)
			cLog.WithFields(logrus.Fields{
				"Test": "CIS", "Num": 3.3, "Result": "Passed",
			}).Info("CloudWatch Log Metric Filter checking for Root found: ", mf)
			return
		}
	}

	if resp.NextToken != nil {
		lookupCloudWatchLogMetricFilter(i, name, logGroupName, resp.NextToken)
		return
	}
	//fmt.Println("CloudWatch Log Metric Filter checking for Root Not found:")
	cLog.WithFields(logrus.Fields{
		"Test": "CIS", "Num": 3.3, "Result": "Failed",
	}).Info("CloudWatch Log Metric Filter checking for Root Not found")
}

func GetLogGroups(svc *cloudwatchlogs.CloudWatchLogs) (result *cloudwatchlogs.DescribeLogGroupsOutput, error error) {
	input := &cloudwatchlogs.DescribeLogGroupsInput{}
	data, err := svc.DescribeLogGroups(input)
	if err != nil {
		return nil, err
	}
	token := data.NextToken
	for token != nil {
		input := &cloudwatchlogs.DescribeLogGroupsInput{
			NextToken: token,
		}
		nextResult, err := svc.DescribeLogGroups(input)
		if err != nil {
			return nil, err
		}
		data.LogGroups = append(data.LogGroups, nextResult.LogGroups...)
		token = nextResult.NextToken
	}
	return data, nil
}

func (i *CloudWatch) Run() {
	var err error
	cLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("CloudWatch Run...")
	//fmt.Println("CloudWatch Run...")
	i.LogGroups, err = GetLogGroups(i.svc)
	if err != nil {
		cLog.WithFields(logrus.Fields{
			"Test": "CIS"}).Info("CloudWatch Groups retrieval error: ", err)
		return
	}
	//i.LogGroups = result
	cLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("CloudWatch Groups: ", i.LogGroups)
	//fmt.Println("LogGroups: ", i.LogGroups)
	for _, groups := range i.LogGroups.LogGroups {
		lookupCloudWatchLogMetricFilter(i, "userIdentity.type", *groups.LogGroupName, nil)
	}
}
