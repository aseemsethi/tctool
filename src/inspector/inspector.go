package inspector

import (
	"github.com/aseemsethi/tctool/src/tcGlobals"

	//	"github.com/aws/aws-sdk-go/aws"
	//	"github.com/aws/aws-sdk-go/aws/awserr"
	//	"github.com/aws/aws-sdk-go/aws/session"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/service/iam"
	"log"
	"strings"
)

type Inspector struct {
	Name string
	Cred string
}

var Access_Key_1_Last_Used_Date = 10
var Access_Key_2_Last_Used_Date = 15

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
		fmt.Println("\n", string(resp.Content))
		fmt.Println(resp.GeneratedTime)
		i.Cred = string(resp.Content)
	}
}

func test1(i *Inspector) {
	s := strings.Split(i.Cred, "\n")

	for _, each := range s {
		//1.1 Avoid the use of the "root" account
		fmt.Println("\n...", each)
		if strings.Contains(each, "<root_account>") {
			root_account := csv.NewReader(strings.NewReader(each))
			record, err := root_account.Read()
			if err != nil {
				log.Fatal(err)
			}
			if record[Access_Key_1_Last_Used_Date] != "N/A" && record[Access_Key_2_Last_Used_Date] != "N/A" {
				fmt.Println("Disable root keys")
			} else {
				fmt.Println("Root good")
			}
		}
		//fmt.Println(index, each)
	}
}

func (i *Inspector) Run() {
	fmt.Printf("\nInspector run..")
	test1(i)
}
