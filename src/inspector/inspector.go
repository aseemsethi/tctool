package inspector

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/ec2"
	//"fmt"
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
	iLog.WithFields(logrus.Fields{"Test": "CIS"}).Info("Inspector Run...")
	sess, _ := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: "default",

		// Provide SDK Config options, such as Region.
		//Config: aws.Config{Region: aws.String("us-east-1")},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})
	//_, err := sess.Config.Credentials.Get()
	//fmt.Println("err: ", err)
	svc := inspector.New(sess)

	//svc := inspector.New(sess, aws.NewConfig().WithRegion("us-east-1"))

	rgi := &inspector.CreateResourceGroupInput{
		ResourceGroupTags: []*inspector.ResourceGroupTag{
			{
				Key:   aws.String("inspector"),
				Value: aws.String("true"),
			},
		},
	}
	iLog.WithFields(logrus.Fields{"Test": "CIS"}).Info("Inspector ResGrp created")
	rg, rgerr := svc.CreateResourceGroup(rgi)
	if rgerr != nil {
		iLog.WithFields(logrus.Fields{"Test": "CIS"}).Info("Inspector ResGrp creation failed:", rgerr)
		return
	}
	iLog.WithFields(logrus.Fields{"Test": "CIS"}).Info("Inspector ResGrp created: ", *rg.ResourceGroupArn)

	//return *rg.ResourceGroupArn

}
