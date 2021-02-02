// Test Compliance Tool (TC)
package main

// Learnings and code taken from various sites like
//		https://github.com/jonhadfield/ape
// https://d1.awsstatic.com/whitepapers/compliance/AWS_CIS_Foundations_Benchmark.pdf - v1.2
import (
	"fmt"
	"github.com/aseemsethi/tctool/src/cloudTrail"
	"github.com/aseemsethi/tctool/src/credReport"
	"github.com/aseemsethi/tctool/src/foundSecurity"
	"github.com/aseemsethi/tctool/src/iam"
	"github.com/aseemsethi/tctool/src/inspector"
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/sirupsen/logrus"
)

// All modules implement this interface
// Using the Structural Component Pattern
type tcIf interface {
	Initialize() bool
	Run()
}

// Main Control Struct for TC Tool
type tc struct {
	cisModules    map[string]tcIf
	foundSecurity foundSecurity.FoundSecurity
	inspector     inspector.InspectorStruct
	name          string
}

var tcTool = tc{name: "Test Compliance Tool"}
var mLog *logrus.Logger

func initTool() {
	// Initialize Global Variables
	tcTool.cisModules = make(map[string]tcIf)
	tcGlobals.Tcg.Initialize()
}

func initModules() bool {
	tcTool.cisModules["credReport"] = &credReport.CredentialReport{Name: "credReport"}
	tcTool.cisModules["Iam"] = &iam.Iam{Name: "Iam"}
	tcTool.cisModules["CloudTrail"] = &cloudTrail.CloudTrail{Name: "CloudTrail"}
	for _, tests := range tcTool.cisModules {
		status := tests.Initialize()
		if status == false {
			return status
		}
	}
	return true
}

func main() {
	fmt.Printf("\nTest Compliance Tool Starting..")

	initTool()
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
		"Test": "AWS Security"}).Info("Test Starting......AWS Foundational Security Best Practices controls............")
	mLog.WithFields(logrus.Fields{
		"Test": "AWS Security"}).Info("Ref: https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-cis-controls.html")
	tcTool.foundSecurity = foundSecurity.FoundSecurity{Name: "Foundational Security"}
	tcTool.foundSecurity.Initialize()
	tcTool.foundSecurity.Run()
	mLog.WithFields(logrus.Fields{
		"Test": "AWS Security"}).Info("Test Completed......AWS Foundational Security Best Practices controls............")
	/*************************** Test3 *******************/
	mLog.WithFields(logrus.Fields{
		"Test": "Inspector"}).Info("**************************** AWS Inspector ***********************************")
	tcTool.inspector = inspector.InspectorStruct{Name: "Inspector"}
	tcTool.inspector.Initialize()
	//tcTool.inspector.Run()
	mLog.WithFields(logrus.Fields{
		"Test": "Inspector"}).Info("Test Completed......AWS Inspector...........")
}
