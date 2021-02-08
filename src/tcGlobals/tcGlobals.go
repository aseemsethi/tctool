package tcGlobals

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"
	"os"
)

type TcGlobals struct {
	Name    string
	Log     *logrus.Logger
	Sess    *session.Session
	GRegion string
	GArn    string
	GConf   aws.Config
}

// from https://github.com/aws/aws-sdk-go-v2/issues/225
type Value string
type Policy struct {
	// 2012-10-17 or 2008-10-17 old policies, do NOT use this for new policies
	Version    string       `json:"Version"`
	Id         string       `json:"Id,omitempty"`
	Statements []Statement1 `json:"Statement"`
}

type Statement1 struct {
	Sid          string           `json:"Sid,omitempty"`          // statement ID, service specific
	Effect       string           `json:"Effect"`                 // Allow or Deny
	Principal    map[string]Value `json:"Principal,omitempty"`    // principal that is allowed or denied
	NotPrincipal map[string]Value `json:"NotPrincipal,omitempty"` // exception to a list of principals
	Action       Value            `json:"Action"`                 // allowed or denied action
	NotAction    Value            `json:"NotAction,omitempty"`    // matches everything except
	Resource     Value            `json:"Resource,omitempty"`     // object or objects that the statement covers
	NotResource  Value            `json:"NotResource,omitempty"`  // matches everything except
	Condition    json.RawMessage  `json:"Condition,omitempty"`    // conditions for when a policy is in effect
}

var Tcg = TcGlobals{Name: "TC Globals"}

// Not used
func CheckPolicy(str *string) {
	fmt.Println("CheckPolicy Testing...")
	var b = []byte(*str)
	var raw interface{}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		fmt.Println("CheckPolicy unmarshal error: ", err)
		return
	}

	var p []string
	//  value can be string or []string, convert everything to []string
	switch v := raw.(type) {
	case string:
		p = []string{v}
	case []interface{}:
		var items []string
		for _, item := range v {
			items = append(items, fmt.Sprintf("%v", item))
		}
		p = items
	default:
		fmt.Println("invalid value element: allowed is only string or []string")
	}
	fmt.Printf("%+v", p)
}

// TBD: Does not check for Principal: *, need to check S3 Policies manually
// User JSON decoder in policyDecoder going forwaard
// str is Jaon Policy formatted *string
func CheckPolicyForAllowAll(str *string) bool {
	var p Policy
	var jsonData = []byte(*str)

	//fmt.Println("Called with string: ", *str)
	err := json.Unmarshal(jsonData, &p)
	if err != nil {
		//fmt.Println("CheckPolicyForAllowAll: unexpected error parsing policy", err)
		Tcg.Log.WithFields(logrus.Fields{
			"Test": "Globals"}).Info("CheckPolicyForAllowAll: unexpected error parsing policy: ", err)
		return false
	}
	//fmt.Printf("%+v", p)
	for _, val := range p.Statements {
		//fmt.Println("\nEffect/Allow: ", val.Effect, val.Principal)
		if val.Effect == "Allow" && val.Principal["AWS"] == "*" {
			return true
		}
	}
	return false
}

func (tcg *TcGlobals) Initialize() bool {
	// Setup common session to be used by all Services
	// Init session in us-east-2
	//sess, err := session.NewSession(&aws.Config{
	//	Region: aws.String("us-east-2")},
	//)
	sess, err := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile:           "default",
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Println("Error creating new session")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	tcg.Sess = sess
	tcg.GRegion = "us-east-1"
	tcg.GArn = "arn:aws:iam::329914591859:role/KVAccess"
	tcg.GConf = aws.Config{Region: aws.String(tcg.GRegion)}
	tcg.GConf.Credentials = stscreds.NewCredentials(tcg.Sess, tcg.GArn, func(p *stscreds.AssumeRoleProvider) {})

	tcg.Log = logrus.New()
	file, err := os.OpenFile("tctool.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		tcg.Log.Fatal(err)
	}
	//defer file.Close()
	tcg.Log.SetOutput(file)
	tcg.Log.SetFormatter(&logrus.JSONFormatter{})
	tcg.Log.SetLevel(logrus.InfoLevel)
	tcg.Log.WithFields(logrus.Fields{
		"Test": "Globals"}).Info("**************************Globals Initialized...")
	return true
}

func (tcg *TcGlobals) Run() {
	tcg.Log.WithFields(logrus.Fields{
		"Test": "Globals"}).Info("nTcGlobals Run...")
}
