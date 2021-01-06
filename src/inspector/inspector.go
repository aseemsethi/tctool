package inspector

import (
	"github.com/aseemsethi/tctool/src/tcGlobals"

	//	"github.com/aws/aws-sdk-go/aws"
	//	"github.com/aws/aws-sdk-go/aws/awserr"
	//	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"

	"fmt"
)

type Inspector struct {
	Name string
}

func (i *Inspector) Initialize() {
	fmt.Printf("\nInspector init..")
	// Create a IAM service client.
	svc := iam.New(tcGlobals.Tcg.Sess)
	resp, err := svc.GenerateCredentialReport(&iam.GenerateCredentialReportInput{})
	if err != nil {
		fmt.Println(err.Error())
	}
	if *resp.State == "COMPLETE" {
		fmt.Printf("\nInspector GetCredRept..")
		resp, err := svc.GetCredentialReport(&iam.GetCredentialReportInput{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(resp.Content))
		fmt.Println(resp.GeneratedTime)
	}
}

func (i *Inspector) Run() {
	fmt.Printf("\nInspector run..")
}
