package iam

import (
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	//"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"
)

type Iam struct {
	Name      string
	svc       iamiface.IAMAPI
	PwdPolicy *iam.GetAccountPasswordPolicyOutput
}

var iamLog *logrus.Logger

func (i *Iam) Initialize() bool {
	iamLog = tcGlobals.Tcg.Log

	// Create a IAM service client.
	svc := iam.New(tcGlobals.Tcg.Sess, &tcGlobals.Tcg.GConf)
	i.svc = svc

	var params *iam.GetAccountPasswordPolicyInput
	resp, err := svc.GetAccountPasswordPolicy(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		iamLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": "1.5 - 1.11", "Result": "Failed",
		}).Info(err.Error())
		iamLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": "1.5 - 1.11", "Result": "Failed",
		}).Info("Password Policy does not exist")
		return false
	}

	// Pretty-print the response data.
	iamLog.WithFields(logrus.Fields{
		"Test": "CIS", "Num": "1.5 - 1.11"}).Info("Password Policy dump: ", resp)
	i.PwdPolicy = resp
	return true
}

func PwdPolicyCheck(i *Iam) {
	if *i.PwdPolicy.PasswordPolicy.RequireUppercaseCharacters ||
		*i.PwdPolicy.PasswordPolicy.RequireLowercaseCharacters ||
		*i.PwdPolicy.PasswordPolicy.RequireNumbers ||
		*i.PwdPolicy.PasswordPolicy.RequireSymbols {
		iamLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": "1.5 - 1.8", "Result": "Failed",
		}).Info("Password Policy doesn't require Uppercase/Lowercase Letters, Numbers and Symbols")
	} else {
		iamLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": "1.5 - 1.8", "Result": "Passed",
		}).Info("Password Policy doesn't require Uppercase/Lowercase Letters, Numbers and Symbols")
	}

	if *i.PwdPolicy.PasswordPolicy.MinimumPasswordLength < 14 {
		iamLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.9, "Result": "Failed",
		}).Info("Minimum Password length less than 14 chars")
	} else {
		iamLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.9, "Result": "Passed",
		}).Info("Minimum Password length is more than 14 chars")
	}

	if i.PwdPolicy.PasswordPolicy.PasswordReusePrevention == nil || *i.PwdPolicy.PasswordPolicy.PasswordReusePrevention < 3 {
		iamLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.10, "Result": "Failed",
		}).Info("Password reuse policy < 3 days or not set - CIS 1.10 failed")
	} else {
		iamLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.10, "Result": "Passed",
		}).Info("Password reuse policy - CIS 1.10 passed")
	}

	if i.PwdPolicy.PasswordPolicy.MaxPasswordAge == nil || *i.PwdPolicy.PasswordPolicy.MaxPasswordAge < 90 {
		iamLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.11, "Result": "Failed",
		}).Info("Passwords don't expire after at least 90 days")
	} else {
		iamLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.11, "Result": "Passed",
		}).Info("Passwords expires after at least 90 days")
	}
}

// TBD: We need to really check MFA HArdware for root
func mfsDeviceCheck(i *Iam) {
	found := false
	mfaDevices, err := i.svc.ListMFADevices(&iam.ListMFADevicesInput{UserName: aws.String("admin")}) // TBD: Need to check for root user only
	if err != nil {
		iamLog.WithFields(logrus.Fields{
			"Test": "CIS", "Num": 1.14,
		}).Info("Failed to list mfa devices - %s", err)
	} else {
		for _, device := range mfaDevices.MFADevices {
			iamLog.WithFields(logrus.Fields{
				"Test": "CIS", "Num": 1.14, "Result": "Passed",
			}).Info("MFA enabled for admin with Device: ", device.SerialNumber)
			found = true
		}
		if found == false {
			iamLog.WithFields(logrus.Fields{
				"Test": "CIS", "Num": 1.14, "Result": "Failed",
			}).Info("MFA not enabled for admin")
		}
	}
}

func (i *Iam) Run() {
	iamLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("IAM Run...")
	PwdPolicyCheck(i)
	mfsDeviceCheck(i)
}
