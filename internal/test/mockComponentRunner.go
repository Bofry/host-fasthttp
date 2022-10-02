package test

import "fmt"

type MockComponentRunner struct {
	prefix string
}

func (c *MockComponentRunner) Start() {
	fmt.Println(c.prefix + ".Start()")
}

func (c *MockComponentRunner) Stop() {
	fmt.Println(c.prefix + ".Stop()")
}
