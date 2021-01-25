package inspector

import (
	"fmt"
	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/sirupsen/logrus"
)

type InspectorStruct struct {
	Name string
}

var iLog *logrus.Logger

func (i *InspectorStruct) Initialize() bool {
	iLog = tcGlobals.Tcg.Log

	return true
}

func (i *InspectorStruct) Run() {
	iLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Inspector Run...")
	sess, _ := session.NewSession()
	svc := inspector.New(sess)
	fmt.Println(svc)
}
