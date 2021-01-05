// Test Compliance Tool (TC)
package main

import (
	"fmt"
	"github.com/aseemsethi/tctool/src/inspector"
	"github.com/aseemsethi/tctool/src/tcGlobals"
)

// All modules implement this interface
// Using the Structural Component Pattern
type tcIf interface {
	Initialize()
	Run()
}

// Main Control Struct for TC Tool
type tc struct {
	tcIfs map[string]tcIf
	name  string
}

var tcg = &tcGlobals.TcGlobals{Name: "TC Globals"}
var tcTool = tc{name: "Test Compliance Tool"}

func initTool() {
	// Initialize Global Variables
	tcTool.tcIfs = make(map[string]tcIf)
	fmt.Printf("\nTC Tool - adding TC Globals Module")
	tcTool.tcIfs["Globals"] = tcg
	tcg.Initialize()
}

func initInspector() {
	// Init Inspector
	in := &inspector.Inspector{Name: "Inspector"}
	fmt.Printf("\nTC Tool - adding Inspector Module")
	tcTool.tcIfs["Inspector"] = in
	in.Initialize()
}

func runInspector() {
	tcTool.tcIfs["Inspector"].Run()
}

func main() {
	fmt.Printf("\nTest Compliance Tool Starting..")

	initTool()
	initInspector()
	runInspector()

	//utils.TestS3()
	fmt.Printf("\nTC Tool - completed, %+v", tcTool)
}
