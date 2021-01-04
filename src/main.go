package main

import (
	"fmt"
	"github.com/aseemsethi/tctool/src/inspector"
)

type tcIf interface {
	Initialize()
	Run()
}

type tc struct {
	tcIfs []tcIf
	name  string
}

//type tcTool tc

func (tcNew tc) add(c tcIf) {
	fmt.Printf("Adding a component to tc tool..")
	tcNew.tcIfs = append(tcNew.tcIfs, c)
}

func main() {
	fmt.Printf("Test Compliance Tool Starting..")
	in := &inspector.Inspector{Name: "Inspector"}
	tcTool := tc{name: "Test Compliance Tool"}
	tcTool.add(in)
}
