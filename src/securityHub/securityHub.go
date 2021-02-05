package securityHub

import (
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/sirupsen/logrus"
)

type SecurityHub struct {
	Name string
}

var fLog *logrus.Logger

func (i *SecurityHub) Initialize() bool {
	fLog = tcGlobals.Tcg.Log

	return true
}

func (i *SecurityHub) Run() {
	fLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("SecurityHub Run...")
}
