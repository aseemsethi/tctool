package main

import (
	"fmt"
	"github.com/aseemsethi/tctool/src/inspector"
	"github.com/aseemsethi/tctool/src/utils"
)

type tcIf interface {
	Initialize()
	Run()
}

type tc struct {
	tcIfs []tcIf
	name  string
}

func (tcNew tc) add(c tcIf) {
	fmt.Printf("\nAdding a component to tc tool..")
	tcNew.tcIfs = append(tcNew.tcIfs, c)
}

func main() {
	fmt.Printf("\nTest Compliance Tool Starting..")
	in := &inspector.Inspector{Name: "Inspector"}
	tcTool := tc{name: "Test Compliance Tool"}
	fmt.Printf("\nTC Tool - adding Inspector Module")
	tcTool.add(in)
	utils.TestS3()
	fmt.Printf("\nTC Tool - completed")
}
