package inspector

import (
	"github.com/aseemsethi/tctool/src/tcGlobals"

	//	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	//	"github.com/aws/aws-sdk-go/aws/session"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/service/iam"
	"io"
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
		resp, get_err := svc.GetCredentialReport(&iam.GetCredentialReportInput{})
		if get_err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeCredentialReportNotPresentException:
					fmt.Println(iam.ErrCodeCredentialReportNotPresentException, aerr.Error())
				case iam.ErrCodeCredentialReportExpiredException:
					fmt.Println(iam.ErrCodeCredentialReportExpiredException, aerr.Error())
				case iam.ErrCodeCredentialReportNotReadyException:
					fmt.Println(iam.ErrCodeCredentialReportNotReadyException, aerr.Error())
				case iam.ErrCodeServiceFailureException:
					fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(get_err.Error())
			}
		}

		//fmt.Println("\n", string(resp.Content))
		//fmt.Println(resp.GeneratedTime)
		i.Cred = string(resp.Content)
	}
}

func RootAccessKeysDisabled(i *Inspector) {
	fmt.Println("RootAccessKeysDisabled - CIS 1.12")
	s := strings.Split(i.Cred, "\n")

	for _, each := range s {
		//1.1 Avoid the use of the "root" account
		//fmt.Println("\n...", each)
		if strings.Contains(each, "<root_account>") {
			root_account := csv.NewReader(strings.NewReader(each))
			record, err := root_account.Read()
			if err != nil {
				log.Fatal(err)
			}
			if record[Access_Key_1_Last_Used_Date] != "N/A" && record[Access_Key_2_Last_Used_Date] != "N/A" {
				fmt.Println("RootAccessKeysDisabled - CIS 1.12 - failed")
			} else {
				fmt.Println("RootAccessKeysDisabled - CIS 1.12 - passed")
			}
		}
		//fmt.Println(index, each)
	}
}

type credentialReport []credentialReportItem

func ParseCredentialFile(i *Inspector) {
	fmt.Println("ParseCredentialFile")
	reader := csv.NewReader(strings.NewReader(i.Cred))
	var readErr error
	var record []string
	//var credReportItem credentialReportItem
	for {
		record, readErr = reader.Read()
		if len(record) > 0 && record[0] == "user" && record[1] == "arn" {
			continue
		}
		if readErr == io.EOF {
			break
		}
		var userName string
		if record[crUser] == "<root_account>" {
			userName = "root"
		} else {
			userName = record[crUser]
		}
		fmt.Println(userName)
	}
}

func (i *Inspector) Run() {
	fmt.Println("\nInspector run..")
	RootAccessKeysDisabled(i)
	ParseCredentialFile(i)
}
