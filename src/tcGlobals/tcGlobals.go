package tcGlobals

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"

	"fmt"
	"os"
)

type TcGlobals struct {
	Name string
	Log  *logrus.Logger
	Sess *session.Session
}

var Tcg = TcGlobals{Name: "TC Globals"}

func (tcg *TcGlobals) Initialize() bool {
	fmt.Printf("\nTcGlobals init..")

	// Setup common session to be used by all Services
	// Init session in us-east-2
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)
	if err != nil {
		fmt.Println("Error creating new session")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	tcg.Sess = sess
	fmt.Println("\nTcGlobals: Session created..")

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
		"Test": "CIS", "Num": 1.0,
	}).Info("Starting CIS v3 Tests.......................")
	return true
}

func (tcg *TcGlobals) Run() {
	fmt.Printf("\nTcGlobals run..")
}
