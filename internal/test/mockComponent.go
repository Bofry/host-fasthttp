package test

import "github.com/Bofry/host"

type MockComponent struct {
}

func (c *MockComponent) Runner() host.Runner {
	return &MockComponentRunner{
		prefix: "MockComponent",
	}
}
