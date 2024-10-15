package test

import (
	"fmt"

	"github.com/Bofry/host"
)

var _ host.Runner = new(MockComponentRunner)

type MockComponentRunner struct {
	prefix string
}

func (c *MockComponentRunner) Start() {
	fmt.Println(c.prefix + ".Start()")
}

func (c *MockComponentRunner) Stop() {
	fmt.Println(c.prefix + ".Stop()")
}
