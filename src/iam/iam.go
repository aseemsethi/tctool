package iam

import (
	"github.com/aseemsethi/tctool/src/tcGlobals"
	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	//"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"

	"fmt"
	"github.com/aws/aws-sdk-go/service/iam"
)

type Iam struct {
	Name      string
	PwdPolicy *iam.GetAccountPasswordPolicyOutput
}

func (i *Iam) Initialize() bool {
	fmt.Println("Iam init..")
	// Create a IAM service client.
	svc := iam.New(tcGlobals.Tcg.Sess)

	var params *iam.GetAccountPasswordPolicyInput
	resp, err := svc.GetAccountPasswordPolicy(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		//fmt.Println(err.Error())
		tcGlobals.Tcg.Log.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.5 - 1.11, "Result": "Failed",
		}).Info(err.Error())
		tcGlobals.Tcg.Log.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.5 - 1.11, "Result": "Failed",
		}).Info("Password Policy does not exist")
		return false
	}

	// Pretty-print the response data.
	fmt.Println(resp)
	i.PwdPolicy = resp
	return true
}

func PwdPolicyCheck(i *Iam) {
	if *i.PwdPolicy.PasswordPolicy.RequireUppercaseCharacters ||
		*i.PwdPolicy.PasswordPolicy.RequireLowercaseCharacters ||
		*i.PwdPolicy.PasswordPolicy.RequireNumbers ||
		*i.PwdPolicy.PasswordPolicy.RequireSymbols {
		tcGlobals.Tcg.Log.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.5 - 1.8, "Result": "Failed",
		}).Info("Password Policy doesn't require Uppercase/Lowercase Letters, Numbers and Symbols")
	} else {
		tcGlobals.Tcg.Log.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.5 - 1.8, "Result": "Passed",
		}).Info("Password Policy doesn't require Uppercase/Lowercase Letters, Numbers and Symbols")
	}

	if *i.PwdPolicy.PasswordPolicy.MinimumPasswordLength < 14 {
		tcGlobals.Tcg.Log.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.9, "Result": "Failed",
		}).Info("Minimum Password length less than 14 chars")
	} else {
		tcGlobals.Tcg.Log.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.9, "Result": "Passed",
		}).Info("Minimum Password length is more than 14 chars")
	}

	if i.PwdPolicy.PasswordPolicy.PasswordReusePrevention == nil || *i.PwdPolicy.PasswordPolicy.PasswordReusePrevention < 3 {
		tcGlobals.Tcg.Log.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.10, "Result": "Failed",
		}).Info("Password reuse policy < 3 days or not set - CIS 1.10 failed")
	} else {
		tcGlobals.Tcg.Log.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.10, "Result": "Passed",
		}).Info("Password reuse policy - CIS 1.10 passed")
	}

	if i.PwdPolicy.PasswordPolicy.MaxPasswordAge == nil || *i.PwdPolicy.PasswordPolicy.MaxPasswordAge < 90 {
		tcGlobals.Tcg.Log.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.11, "Result": "Failed",
		}).Info("Passwords don't expire after at least 90 days")
	} else {
		tcGlobals.Tcg.Log.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.11, "Result": "Passed",
		}).Info("Passwords expires after at least 90 days")
	}
}

func (i *Iam) Run() {
	tcGlobals.Tcg.Log.WithFields(logrus.Fields{
		"Test": "CIS",
	}).Info("IAM Run")
	PwdPolicyCheck(i)
}
