package securityHub

import (
	"fmt"
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/sirupsen/logrus"
)

type SecurityHub struct {
	Name string
	svc  *securityhub.SecurityHub
}

var sLog *logrus.Logger

func (i *SecurityHub) Initialize() bool {
	sLog = tcGlobals.Tcg.Log
	i.svc = securityhub.New(tcGlobals.Tcg.Sess, &tcGlobals.Tcg.GConf)

	input := &securityhub.EnableSecurityHubInput{}
	_, err := i.svc.EnableSecurityHub(input)
	if err != nil {
		//fmt.Println("failed EnableSecurityHub: %s", err)
		sLog.WithFields(logrus.Fields{"Test": "SecurityHub"}).Info("Not Enabled: ", err)
	}
	//fmt.Println("EnableSecurityHub...")
	sLog.WithFields(logrus.Fields{"Test": "SecurityHub"}).Info("Enabled")

	return true
}

func listFindings(i *SecurityHub) {
	var nextToken *string
	for {
		input := &securityhub.GetFindingsInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		}
		list, err := i.svc.GetFindings(input)
		if err != nil {
			sLog.WithFields(logrus.Fields{"Test": "securityhub"}).Info("ListFindings failed: ", err)
			return
		}
		sLog.WithFields(logrus.Fields{"Test": "securityhub"}).Info("ListFindings passed")
		for _, v := range list.Findings {
			sLog.WithFields(logrus.Fields{"Test": "securityhub"}).Info(v)
			fmt.Println("Findings: ", v)
		}
		if list.NextToken != nil {
			nextToken = list.NextToken
		} else {
			break
		}
	}
}

func (i *SecurityHub) Run() {
	sLog.WithFields(logrus.Fields{
		"Test": "SecurityHub"}).Info("SecurityHub Run...")
	input := &securityhub.GetEnabledStandardsInput{}
	output, err := i.svc.GetEnabledStandards(input)
	if err != nil {
		sLog.WithFields(logrus.Fields{"Test": "SecurityHub"}).Info("GetEnabledStandards failed: ", err)
		return
	}
	sLog.WithFields(logrus.Fields{"Test": "SecurityHub"}).Info("GetEnabledStandards: ", output.StandardsSubscriptions)
	fmt.Println("Enabled: ", output.StandardsSubscriptions)
	listFindings(i)
}
