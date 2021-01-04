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
	//tcIfs []tcIf
	tcIfs map[string]tcIf
	name  string
}

func (tcNew tc) add(c tcIf) {
	//fmt.Printf("\nAdding a component to tc tool..")
	//tcNew.tcIfs = append(tcNew.tcIfs, c)
}

func main() {
	fmt.Printf("\nTest Compliance Tool Starting..")

	tcTool := tc{name: "Test Compliance Tool"}
	tcTool.tcIfs = make(map[string]tcIf)

	tcg := &tcGlobals.TcGlobals{Name: "TC Globals"}
	fmt.Printf("\nTC Tool - adding TC Globals Module")
	tcTool.tcIfs["Globals"] = tcg

	in := &inspector.Inspector{Name: "Inspector"}
	fmt.Printf("\nTC Tool - adding Inspector Module")
	tcTool.tcIfs["Inspector"] = in

	//utils.TestS3()
	fmt.Printf("\nTC Tool - completed, %+v", tcTool)
}
