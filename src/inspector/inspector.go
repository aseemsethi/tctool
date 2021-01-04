package inspector

import (
	"fmt"
)

type Inspector struct {
	Name string
}

func (i *Inspector) Initialize() {
	fmt.Printf("Inspector init..")
}

func (i *Inspector) Run() {
	fmt.Printf("Inspector run..")
}
