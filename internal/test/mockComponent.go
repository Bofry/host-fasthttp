package test

import "github.com/Bofry/host"

var _ host.Runable = new(MockComponent)

type MockComponent struct {
}

func (c *MockComponent) Runner() host.Runner {
	return &MockComponentRunner{
		prefix: "MockComponent",
	}
}
