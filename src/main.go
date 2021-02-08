// Test Compliance Tool (TC)
package main

// Learnings and code taken from various sites like
//		https://github.com/jonhadfield/ape
// https://d1.awsstatic.com/whitepapers/compliance/AWS_CIS_Foundations_Benchmark.pdf - v1.2
import (
	"fmt"
	"github.com/aseemsethi/tctool/src/cis/cloudTrail"
	"github.com/aseemsethi/tctool/src/cis/cloudWatch"
	"github.com/aseemsethi/tctool/src/cis/credReport"
	"github.com/aseemsethi/tctool/src/cis/iam"
	"github.com/aseemsethi/tctool/src/inspector"
	"github.com/aseemsethi/tctool/src/securityHub"
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/sirupsen/logrus"
	"os"
)

// All modules implement this interface
// Using the Structural Component Pattern
type tcIf interface {
	Initialize() bool
	Run()
}

// Main Control Struct for TC Tool
type tc struct {
	cisModules  map[string]tcIf
	securityHub securityHub.SecurityHub
	inspector   inspector.InspectorStruct
	name        string
}

var tcTool = tc{name: "Test Compliance Tool"}
var mLog *logrus.Logger

func initTool(region string, account string) {
	// Initialize Global Variables
	tcTool.cisModules = make(map[string]tcIf)
	tcGlobals.Tcg.Initialize(region, account)
}

func initModules() bool {
	tcTool.cisModules["credReport"] = &credReport.CredentialReport{Name: "credReport"}
	tcTool.cisModules["Iam"] = &iam.Iam{Name: "Iam"}
	tcTool.cisModules["CloudTrail"] = &cloudTrail.CloudTrail{Name: "CloudTrail"}
	tcTool.cisModules["CloudWatch"] = &cloudWatch.CloudWatch{Name: "CloudWatch"}
	for _, tests := range tcTool.cisModules {
		status := tests.Initialize()
		if status == false {
			return status
		}
	}
	return true
}

// Call with tctool <region> <accountid>
func main() {
	fmt.Printf("\nTest Compliance Tool Starting..")

	if len(os.Args) < 3 {
		fmt.Println("Usage: tctool <region> <accountid>")
		return
	}
	initTool(os.Args[1], os.Args[2])
	mLog = tcGlobals.Tcg.Log
	status := initModules()
	if status == false {
		mLog.WithFields(logrus.Fields{
			"Test": "CIS"}).Info("Modules init error: exit")
		return
	}
	/*************************** Test1 *******************/
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Test Compliance Starting...CIS AWS Foundations Benchmark controls............")
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Ref: https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-cis-controls.html")
	// Run the 1st Test Plan - 49 Security Config Controls
	// CIS v3 - https://d1.awsstatic.com/whitepapers/compliance/AWS_CIS_Foundations_Benchmark.pdf
	for _, tests := range tcTool.cisModules {
		tests.Run()
	}
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Test Compliance Completed...CIS AWS Foundations Benchmark controls............")

	/*************************** Test2 *******************/
	mLog.WithFields(logrus.Fields{
		"Test": "Inspector"}).Info("**************************** AWS Inspector ***********************************")
	tcTool.inspector = inspector.InspectorStruct{Name: "Inspector"}
	tcTool.inspector.Initialize()
	//tcTool.inspector.Run()
	mLog.WithFields(logrus.Fields{
		"Test": "Inspector"}).Info("Test Completed......AWS Inspector...........")

	/*************************** Test3 *******************/
	mLog.WithFields(logrus.Fields{
		"Test": "SecurityHub"}).Info("**************************** AWS Security Hub *******************************")
	tcTool.securityHub = securityHub.SecurityHub{Name: "Security Hub"}
	tcTool.securityHub.Initialize()
	tcTool.securityHub.Run()
	mLog.WithFields(logrus.Fields{
		"Test": "SecurityHub"}).Info("Test Completed......AWS Security Hub............")
}
