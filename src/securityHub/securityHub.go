package securityHub

import (
	"fmt"
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/sirupsen/logrus"
)

type SecurityHub struct {
	Name string
}

var sLog *logrus.Logger

func (i *SecurityHub) Initialize() bool {
	sLog = tcGlobals.Tcg.Log
	input := &securityhub.EnableSecurityHubInput{}
	_, err := securityhub.New(tcGlobals.Tcg.Sess, &tcGlobals.Tcg.GConf).EnableSecurityHub(input)
	if err != nil {
		fmt.Println("failed EnableSecurityHub: %s", err)
		sLog.WithFields(logrus.Fields{"Test": "SecurityHub"}).Info("Not Enabled: ", err)
		return false
	}
	fmt.Println("EnableSecurityHub...")
	sLog.WithFields(logrus.Fields{"Test": "SecurityHub"}).Info("Enabled")

	return true
}

func (i *SecurityHub) Run() {
	sLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("SecurityHub Run...")
}
