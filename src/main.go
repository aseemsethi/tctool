// Test Compliance Tool (TC)
package main

// Learnings and code taken from various sites like
//		https://github.com/jonhadfield/ape
// https://d1.awsstatic.com/whitepapers/compliance/AWS_CIS_Foundations_Benchmark.pdf - v1.2
import (
	"fmt"
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
	tcIfs map[string]tcIf
	name  string
}

var tcTool = tc{name: "Test Compliance Tool"}
var mLog *logrus.Logger

func initTool() {
	// Initialize Global Variables
	tcTool.tcIfs = make(map[string]tcIf)
	fmt.Printf("\nTC Tool - adding TC Globals Module")
	tcTool.tcIfs["Globals"] = &tcGlobals.Tcg
	tcTool.tcIfs["Globals"].Initialize()
}

func initInspector() bool {
	// Init Inspector
	in := &inspector.Inspector{Name: "Inspector"}
	fmt.Printf("\nTC Tool - adding Inspector Module")
	tcTool.tcIfs["Inspector"] = in
	cont := in.Initialize()
	return cont
}

func initIam() bool {
	// Init Iam
	in := &iam.Iam{Name: "Iam"}
	fmt.Printf("\nTC Tool - adding IAM Module")
	tcTool.tcIfs["Iam"] = in
	cont := in.Initialize()
	return cont
}

func main() {
	fmt.Printf("\nTest Compliance Tool Starting..")

	initTool()
	mLog = tcGlobals.Tcg.Log
	/*************************** Test1 *******************/
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Test Compliance Starting...CIS AWS Foundations Benchmark controls............")
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Ref: https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-cis-controls.html")
	// Run the 1st Test Plan - 49 Security Config Controls
	// CIS v3 - https://d1.awsstatic.com/whitepapers/compliance/AWS_CIS_Foundations_Benchmark.pdf
	statusInspector := initInspector() // Credential Report
	if statusInspector == false {
		return
	} else {
		tcTool.tcIfs["Inspector"].Run()
	}
	statusIam := initIam() // IAM Password Policy Report
	// The pwd policy is created and retrievable only when the admin goes to IAM->AccSettings->PasswordPolicy
	if statusIam == false {
		return
	} else {
		tcTool.tcIfs["Iam"].Run()
	}
	//utils.TestS3()
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Test Compliance Completed...CIS AWS Foundations Benchmark controls............")
	/*************************** Test2 *******************/
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Test Starting......AWS Foundational Security Best Practices controls............")
	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Ref: https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-cis-controls.html")

	mLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Test Completed......AWS Foundational Security Best Practices controls............")

	/*************************** Test3 *******************/

}
