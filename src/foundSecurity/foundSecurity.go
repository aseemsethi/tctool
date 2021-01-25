package foundSecurity

import (
	"github.com/aseemsethi/tctool/src/tcGlobals"
	"github.com/sirupsen/logrus"
)

type FoundSecurity struct {
	Name string
}

var fLog *logrus.Logger

func (i *FoundSecurity) Initialize() bool {
	fLog = tcGlobals.Tcg.Log

	return true
}

func (i *FoundSecurity) Run() {
	fLog.WithFields(logrus.Fields{
		"Test": "CIS"}).Info("Foundational Security Run...")
}
