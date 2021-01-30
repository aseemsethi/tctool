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

func getSpecificTagValue(key string, tags []*ec2.Tag) string {
	for _, tag := range tags {
		if *(tag.Key) == key {
			return *tag.Value
		}
	}
	return "--"
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
			inspectorTag := getSpecificTagValue("inspector", inst.Tags)
			fmt.Println("\nType", *inst.InstanceType, " ID: ", *inst.InstanceId, " State: ", *inst.State.Name, " InspectorTag: ", inspectorTag)
			if inspectorTag == "true" {
				fmt.Println("Included in Inspector run")
			}
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
	iLog.WithFields(logrus.Fields{"Test": "CIS"}).Info("Inspector Asessment Target created: ", at)
	fmt.Println("AssessmentTarget: ", at)
	//return *at.AssessmentTargetArn

	// 3. create rules package input
	rpi := &inspector.ListRulesPackagesInput{
		MaxResults: aws.Int64(100),
	}
	rp, erp := svc.ListRulesPackages(rpi)
	if erp != nil {
		fmt.Println(erp.Error())
		return
	}
	fmt.Println("List Rules Pkg: ", rp) // we selct all rules, i,e, N/W, CVE etc.

	// 4. create assessment template
	atli := &inspector.CreateAssessmentTemplateInput{
		AssessmentTargetArn:    aws.String(*at.AssessmentTargetArn),
		AssessmentTemplateName: aws.String("InspectorRun" + "_AssessmentTemplate_" + time.Now().Format("2006-01-02_15.04.05")),
		DurationInSeconds:      aws.Int64(300),
		RulesPackageArns:       rp.RulesPackageArns,
		UserAttributesForFindings: []*inspector.Attribute{
			{
				Key:   aws.String("inspection-type"),
				Value: aws.String("InspectorRun"),
			},
		},
	}

	atl, atlerr := svc.CreateAssessmentTemplate(atli)
	if atlerr != nil {
		fmt.Println(atlerr)
		return
	}
	fmt.Println("Asessment Template: ", atl)

	// 6. start assessment template run
	ari := &inspector.StartAssessmentRunInput{
		AssessmentRunName:     aws.String("InspectorRun" + "_Run_" + time.Now().Format("2006-01-02_15.04.05")),
		AssessmentTemplateArn: aws.String(*atl.AssessmentTemplateArn),
	}

	ar, arerr := svc.StartAssessmentRun(ari)
	if arerr != nil {
		fmt.Println(arerr.Error())
	}
	fmt.Println("Asessment Run start: ", ar)
	time.Sleep(360 * time.Second)

	fmt.Println("Asessment Run complete: ")
	input := &inspector.ListFindingsInput{
		AssessmentRunArns: []*string{
			aws.String(*ar.AssessmentRunArn),
		},
		MaxResults: aws.Int64(123),
	}
	list, err := svc.ListFindings(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Findings: ", list)

	var resp []float64
	for _, v := range list.FindingArns {
		input := &inspector.DescribeFindingsInput{
			FindingArns: []*string{
				aws.String(*v),
			},
		}
		fmt.Println("String: ", *v)

		result, err := svc.DescribeFindings(input)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(*result.Findings[0].NumericSeverity)
		fmt.Println(*result.Findings[0])
		fmt.Println("\n")
		resp = append(resp, *result.Findings[0].NumericSeverity)
	}
	fmt.Println("Findings: ", resp)
}
