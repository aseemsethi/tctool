package main

import (
	"fmt"
	"github.com/aseemsethi/tctool/tcModules"
)

type tcIf interface {
	init()
	run()
}

type tc struct {
	tcIfs []tcIf
	name  string
}

func (tcNew *tc) add(c tcIf) {
	fmt.Printf("Adding a component to tc tool..")
	tcNew.tcIfs = append(tcNew.tcIfs, c)
}

func main() {
	fmt.Printf("Test Compliance Tool Starting..")
	in := &inspector{name: "Inspector"}
	tcTool := &tc{name: "Test Compliance Tool"}
	tcTool.add(in)
}
