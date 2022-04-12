package internal

import (
	"net"
	"strings"
)

// reference https://golang.org/src/net/ipsock.go?s=4766:4832#L146
func splitHostPort(hostport string) (host, port string, err error) {
	if len(hostport) == 0 {
		return "", "", nil
	}

	const (
		missingPort   = "missing port in address"
		tooManyColons = "too many colons in address"
	)
	addrErr := func(addr, why string) (host, port string, err error) {
		return "", "", &net.AddrError{Err: why, Addr: addr}
	}
	j, k := 0, 0

	// The port starts after the last colon.
	i := strings.LastIndexByte(hostport, ':')

	if hostport[0] == '[' {
		// Expect the first ']' just before the last ':'.
		end := strings.IndexByte(hostport, ']')
		if end < 0 {
			return addrErr(hostport, "missing ']' in address")
		}
		switch end + 1 {
		case len(hostport):
			// There can't be a ':' behind the ']' now.
			host = hostport
		case i:
			// The expected result.
		default:
			// Either ']' isn't followed by a colon, or it is
			// followed by a colon that is not the last one.
			if hostport[end+1] == ':' {
				return addrErr(hostport, tooManyColons)
			}
		}
		host = hostport[1:end]
		j, k = 1, end+1 // there can't be a '[' resp. ']' before these positions
	} else {
		switch {
		case i >= 0:
			host = hostport[:i]
			if strings.IndexByte(host, ':') >= 0 {
				return addrErr(hostport, tooManyColons)
			}
		default:
			host = hostport
		}
	}
	if strings.IndexByte(hostport[j:], '[') >= 0 {
		return addrErr(hostport, "unexpected '[' in address")
	}
	if strings.IndexByte(hostport[k:], ']') >= 0 {
		return addrErr(hostport, "unexpected ']' in address")
	}

	if i >= 0 {
		port = hostport[i+1:]
	}
	return host, port, nil
}
