package iam

import (
	"github.com/aseemsethi/tctool/src/tcGlobals"

	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	//"github.com/aws/aws-sdk-go/aws/session"
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
		fmt.Println(err.Error())
		fmt.Println("Password Policy does not exist: CIS 1.5 - 1.11 failed")
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
		fmt.Println("Password Policy doesn't require Uppercase/Lowercase Letters, Numbers and Symbols - CIS 1.5 - 1.8 failed")
	} else {
		fmt.Println("Password Policy doesn't require Uppercase/Lowercase Letters, Numbers and Symbols - CIS 1.5 - 1.8 passed")
	}

	if *i.PwdPolicy.PasswordPolicy.MinimumPasswordLength < 14 {
		fmt.Println("Minimum Password length less than 14 chars - CIS 1.9 failed")
	} else {
		fmt.Println("Minimum Password length less than 14 chars - CIS 1.9 passed")
	}

	if i.PwdPolicy.PasswordPolicy.PasswordReusePrevention == nil || *i.PwdPolicy.PasswordPolicy.PasswordReusePrevention < 3 {
		fmt.Println("Last 3 Passwords can be reused - CIS 1.10 failed")
	} else {
		fmt.Println("Minimum Password length less than 14 chars - CIS 1.10 passed")
	}

	if i.PwdPolicy.PasswordPolicy.MaxPasswordAge == nil || *i.PwdPolicy.PasswordPolicy.MaxPasswordAge < 90 {
		fmt.Println("Passwords don't expire after at least 90 days - CIS 1.11 failed")
	} else {
		fmt.Println("Passwords don't expire after at least 90 days - CIS 1.11 failed")
	}
}

func (i *Iam) Run() {
	fmt.Println("Iam run..")
	PwdPolicyCheck(i)
}
