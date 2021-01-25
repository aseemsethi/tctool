// Test Compliance Tool (TC)
package main

// Learnings and code taken from various sites like
//		https://github.com/jonhadfield/ape
// https://d1.awsstatic.com/whitepapers/compliance/AWS_CIS_Foundations_Benchmark.pdf - v1.2
import (
	"fmt"
	"github.com/aseemsethi/tctool/src/cloudTrail"
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
	tcIfs         map[string]tcIf
	foundSecurity foundSecurity.FoundSecurity
	name          string
}

var tcTool = tc{name: "Test Compliance Tool"}
var mLog *logrus.Logger

func initTool() {
	// Initialize Global Variables
	tcTool.tcIfs = make(map[string]tcIf)
	tcTool.tcIfs["Globals"] = &tcGlobals.Tcg
	tcTool.tcIfs["Globals"].Initialize()
}

func initModules() bool {
	tcTool.tcIfs["Inspector"] = &inspector.Inspector{Name: "Inspector"}
	tcTool.tcIfs["Iam"] = &iam.Iam{Name: "Iam"}
	tcTool.tcIfs["CloudTrail"] = &cloudTrail.CloudTrail{Name: "CloudTrail"}
	for _, tests := range tcTool.tcIfs {
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
	for _, tests := range tcTool.tcIfs {
		tests.Run()
	}
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Test Compliance Completed...CIS AWS Foundations Benchmark controls............")

	/*************************** Test2 *******************/
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Test Starting......AWS Foundational Security Best Practices controls............")
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Ref: https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-cis-controls.html")
	tcTool.foundSecurity = foundSecurity.FoundSecurity{Name: "Foundational Security"}
	tcTool.foundSecurity.Initialize()
	tcTool.foundSecurity.Run()
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Test Completed......AWS Foundational Security Best Practices controls............")
	/*************************** Test3 *******************/

}
