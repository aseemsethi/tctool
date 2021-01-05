package tcGlobals

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"fmt"
	"os"
)

type TcGlobals struct {
	Name string
	Sess *session.Session
}

func (tcg *TcGlobals) Initialize() {
	fmt.Printf("TcGlobals init..")

	// Setup common session to be used by all Services
	// Init session in us-east-2
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)
	if err != nil {
		fmt.Println("Ecreating new session")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	tcg.Sess = sess
	fmt.Println("TcGlobals: Session created..")
}

func (tcg *TcGlobals) Run() {
	fmt.Printf("TcGlobals run..")
}
