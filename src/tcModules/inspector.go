package main

import (
	"fmt"
)

type inspector struct {
	name string
}

func (i *inspector) init() {
	fmt.Printf("Inspector init..")
}

func (i *inspector) run() {
	fmt.Printf("Inspector run..")
}
