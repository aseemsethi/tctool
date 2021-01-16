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
	Name           string
	PasswordPolicy string
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
	//i.PasswordPolicy = string(resp)
	return true
}

func PwdPolicyOneUpperCaseLetter(i *Iam) {

}

func (i *Iam) Run() {
	fmt.Println("Iam run..")
	PwdPolicyOneUpperCaseLetter(i)
}
