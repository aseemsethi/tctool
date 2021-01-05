package inspector

import (
	"fmt"
)

type Inspector struct {
	Name string
}

func (i *Inspector) Initialize() {
	fmt.Printf("\nInspector init..")
}

func (i *Inspector) Run() {
	fmt.Printf("\nInspector run..")
}
