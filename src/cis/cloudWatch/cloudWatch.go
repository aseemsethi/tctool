package cloudWatch

import (
	"fmt"
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/sirupsen/logrus"
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
		FilterNamePrefix: aws.String(name),
		LogGroupName:     aws.String(logGroupName),
		NextToken:        nextToken,
	}
	fmt.Printf("Reading CloudWatch Log Metric Filter: %s", input)
	resp, err := i.svc.DescribeMetricFilters(&input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "ResourceNotFoundException" {
			fmt.Println("CloudWatch Log Metric Filter Not Found: ResourceNotFoundException")
		}
		fmt.Println("CloudWatch Log Metric Filter Not Found: %s", err)
	}
	for _, mf := range resp.MetricFilters {
		if *mf.FilterName == name {
			fmt.Println("CloudWatch Log Metric Filter Found:", mf)
		}
	}

	if resp.NextToken != nil {
		lookupCloudWatchLogMetricFilter(i, name, logGroupName, resp.NextToken)
		return
	}
	fmt.Println("CloudWatch Log Metric Filter Not Found: %s", err)
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
	cLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("CloudWatch Run...")
	fmt.Println("CloudWatch Run...")
	result, err := GetLogGroups(i.svc)
	i.LogGroups = result
	fmt.Println(result, err)
	//fmt.Println("LogGroups: ", i.LogGroups)
	//lookupCloudWatchLogMetricFilter(i, )
}
