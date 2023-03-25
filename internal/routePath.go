package internal

import "fmt"

var (
	_ fmt.Stringer = new(RoutePath)
)

type RoutePath struct {
	Method string
	Path   string
}

func (p *RoutePath) String() string {
	return p.Method + " " + p.Path
}
