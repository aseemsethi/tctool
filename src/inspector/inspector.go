package inspector

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/ec2"
	"fmt"
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/sirupsen/logrus"
	"time"
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
		// Using the NewSessionWithOptions with SharedConfigState set to SharedConfigEnable will
		// create the session as if the AWS_SDK_LOAD_CONFIG environment variable was set.
		SharedConfigState: session.SharedConfigEnable,
	})
	//_, err := sess.Config.Credentials.Get()
	//fmt.Println("err: ", err)
	svc := inspector.New(sess)

	/** EC2 reading **/
	ec2Svc := ec2.New(sess)
	ec2Instances, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		iLog.WithFields(logrus.Fields{"Test": "CIS"}).Info("Inspector cannot get ec2s: ", err)
		return
	}
	for idx := range ec2Instances.Reservations {
		for _, inst := range ec2Instances.Reservations[idx].Instances {
			fmt.Println("\t# Type", *inst.InstanceType, " ID: ", *inst.InstanceId, " State: ", *inst.State.Name)
		}
	}
	/**********/

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
	iLog.WithFields(logrus.Fields{"Test": "CIS"}).Info("Inspector Resource Group created: ", *rg.ResourceGroupArn)
	//return *rg.ResourceGroupArn

	// 2. Create assessment target
	ati := &inspector.CreateAssessmentTargetInput{
		AssessmentTargetName: aws.String("InspectorRun" + "_AssessmentTarget_" + time.Now().Format("2006-01-02_15.04.05")),
		ResourceGroupArn:     rg.ResourceGroupArn,
	}
	at, aterr := svc.CreateAssessmentTarget(ati)
	if aterr != nil {
		iLog.WithFields(logrus.Fields{"Test": "CIS"}).Info("Inspector Asessment Target ceration failed: ", aterr)
		return
	}
	iLog.WithFields(logrus.Fields{"Test": "CIS"}).Info("Inspector Asessment Target created")
	fmt.Println("at: ", at)

	//return *at.AssessmentTargetArn

}
