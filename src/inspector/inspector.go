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
	"time"
)

type Inspector struct {
	Name       string
	Cred       string
	CredReport credentialReport
}

var Access_Key_1_Last_Used_Date = 10
var Access_Key_2_Last_Used_Date = 15

func (i *Inspector) Initialize() bool {
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
		return true
	} else {
		fmt.Println("Report not created as yet")
		return false
	}
}

func RootAccessKeysDisabled(i *Inspector) {
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

func ParseCredentialFile(i *Inspector) {
	var err error
	var credReportItem credentialReportItem

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
		//fmt.Println(userName)
		var (
			passwordEnabled, mfaActive, accessKey1Active, accessKey2Active, cert1Active, cert2Active bool
			userCreationTime, passwordLastUsed, passwordLastChanged, passwordNextRotation,
			accessKey1LastRotated, accessKey1LastUsedDate, accessKey2LastRotated, accessKey2LastUsedDate,
			cert1LastRotated, cert2LastRotated time.Time
		)
		userCreationTime, err = time.Parse(time.RFC3339, record[crUserCreationTime])
		if err != nil {
			// Invoking an empty time.Time struct literal will return Go's zero date.
			userCreationTime = time.Time{}
		}

		passwordEnabled = stringToBool(record[crPasswordEnabled])

		passwordLastUsed, err = time.Parse(time.RFC3339, record[crPasswordLastUsed])
		if err != nil {
			passwordLastUsed = time.Time{}
		}
		passwordLastChanged, err = time.Parse(time.RFC3339, record[crPasswordLastChanged])
		if err != nil {
			passwordLastChanged = time.Time{}
		}

		passwordNextRotation, err = time.Parse(time.RFC3339, record[crPasswordNextRotation])
		if err != nil {
			passwordNextRotation = time.Time{}
		}
		mfaActive = stringToBool(record[crMfaActive])
		accessKey1Active = stringToBool(record[crAccessKey1Active])

		accessKey1LastRotated, err = time.Parse(time.RFC3339, record[crAccessKey1LastRotated])
		if err != nil {
			accessKey1LastRotated = time.Time{}
		}
		accessKey1LastUsedDate, err = time.Parse(time.RFC3339, record[crAccessKey1LastUsedDate])
		if err != nil {
			accessKey1LastUsedDate = time.Time{}
		}
		accessKey2Active = stringToBool(record[crAccessKey2Active])

		accessKey2LastRotated, err = time.Parse(time.RFC3339, record[crAccessKey2LastRotated])
		if err != nil {
			accessKey2LastRotated = time.Time{}
		}
		accessKey2LastUsedDate, err = time.Parse(time.RFC3339, record[crAccessKey2LastUsedDate])
		if err != nil {
			accessKey2LastUsedDate = time.Time{}
		}
		cert1Active = stringToBool(record[crCert1Active])

		cert1LastRotated, err = time.Parse(time.RFC3339, record[crCert1LastRotated])
		if err != nil {
			cert1LastRotated = time.Time{}
		}
		cert2Active = stringToBool(record[crCert2Active])

		cert2LastRotated, err = time.Parse(time.RFC3339, record[crCert2LastRotated])
		if err != nil {
			cert2LastRotated = time.Time{}
			err = nil
		}

		credReportItem = credentialReportItem{
			Arn:                       record[crArn],
			User:                      userName,
			UserCreationTime:          userCreationTime,
			PasswordEnabled:           passwordEnabled,
			PasswordLastUsed:          passwordLastUsed,
			PasswordLastChanged:       passwordLastChanged,
			PasswordNextRotation:      passwordNextRotation,
			MfaActive:                 mfaActive,
			AccessKey1Active:          accessKey1Active,
			AccessKey1LastRotated:     accessKey1LastRotated,
			AccessKey1LastUsedDate:    accessKey1LastUsedDate,
			AccessKey1LastUsedRegion:  record[crAccessKey1LastUsedRegion],
			AccessKey1LastUsedService: record[crAccessKey1LastUsedService],
			AccessKey2Active:          accessKey2Active,
			AccessKey2LastRotated:     accessKey2LastRotated,
			AccessKey2LastUsedDate:    accessKey2LastUsedDate,
			AccessKey2LastUsedRegion:  record[crAccessKey2LastUsedRegion],
			AccessKey2LastUsedService: record[crAccessKey2LastUsedService],
			Cert1Active:               cert1Active,
			Cert1LastRotated:          cert1LastRotated,
			Cert2Active:               cert2Active,
			Cert2LastRotated:          cert2LastRotated,
		}
		i.CredReport = append(i.CredReport, credReportItem)
		//fmt.Printf("%+v", credReportItem)
	}
}

func MFAEnabled(i *Inspector) {
	failed := false
	for _, elem := range i.CredReport {
		//fmt.Println("Check User: ", elem.Arn)
		if elem.MfaActive == false {
			fmt.Println("MFAEnabled - CIS 1.2, 1.13 - failed for User", elem.Arn)
			failed = true
		}
	}
	if failed == false {
		fmt.Println("MFAEnabled - CIS 1.2, 1.13 - Passed")
	}
}

func TimeLastUsedAccessKeys(i *Inspector) {
	failed := false
	for _, elem := range i.CredReport {
		// If the AccessKey is never used, it will show as N/A, and a time coversion on this will yield an error
		// At that tiem, we save null vaule in this time field
		if elem.AccessKey1LastUsedDate.IsZero() == true {
			fmt.Println("AccessKey usage - CIS 1.3 - credentials never used - failed for User", elem.Arn)
			failed = true
		} else {
			diff := time.Now().Sub(elem.AccessKey1LastUsedDate).Hours()
			diff1 := fmt.Sprintf("%.1f", diff)
			//fmt.Println("Time elapsed for User: ", elem.Arn, " is ", diff1, " Hours")
			if diff > 90*24 {
				fmt.Println("TimeLastUsedAccessKeys - CIS 1.3 - failed. Last used Hrs:", diff1, " user: ", elem.Arn)
				failed = true
			} else {
				fmt.Println("TimeLastUsedAccessKeys - CIS 1.3 - passed. Last rotated Hrs:", diff1, " user: ", elem.Arn)
			}
		}
	}
	if failed == false {
		fmt.Println("MFAEnabled - CIS 1.2 - Passed")
	}
}

func TimeLastRotatedAccessKeys(i *Inspector) {
	failed := false
	for _, elem := range i.CredReport {
		// If the AccessKey is never used, it will show as N/A, and a time coversion on this will yield an error
		// At that tiem, we save null vaule in this time field
		if elem.AccessKey1LastRotated.IsZero() == true {
			fmt.Println("AccessKey rotated - CIS 1.4 -  not rotated 90 days - failed for User", elem.Arn)
			failed = true
		} else {
			diff := time.Now().Sub(elem.AccessKey1LastRotated).Hours()
			diff1 := fmt.Sprintf("%.1f", diff)
			//fmt.Println("Time elapsed for User: ", elem.Arn, " is ", diff1, " Hours")
			if diff > 90*24 {
				fmt.Println("TimeLastRotatedAccessKeys - CIS 1.4 - failed. Last rotated Hrs:", diff1, " user: ", elem.Arn)
				failed = true
			} else {
				fmt.Println("TimeLastRotatedAccessKeys - CIS 1.4 - passed. Last rotated Hrs:", diff1, " user: ", elem.Arn)
			}
		}
	}
	if failed == false {
		fmt.Println("MFAEnabled - CIS 1.2 - Passed")
	}
}

func (i *Inspector) Run() {
	fmt.Println("\nInspector run..")
	RootAccessKeysDisabled(i)
	ParseCredentialFile(i)
	MFAEnabled(i)
	TimeLastUsedAccessKeys(i)
	TimeLastRotatedAccessKeys(i)
}
