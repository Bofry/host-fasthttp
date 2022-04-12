package internal

import "testing"

func TestSplitHostPort_EmptyHostAndPort80(t *testing.T) {
	var hostport = ":80"
	var (
		expectedHost = ""
		expectedPort = "80"
	)
	host, port, err := splitHostPort(hostport)
	if err != nil {
		t.Error(err)
	}
	if host != expectedHost {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedPort, port)
	}
}

func TestSplitHostPort_EmptyHostAndEmptyPort(t *testing.T) {
	var hostport = ""
	var (
		expectedHost = ""
		expectedPort = ""
	)
	host, port, err := splitHostPort(hostport)
	if err != nil {
		t.Error(err)
	}
	if host != expectedHost {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedPort, port)
	}
}

func TestSplitHostPort_HostLocalAndEmptyPort(t *testing.T) {
	var hostport = "127.0.0.1"
	var (
		expectedHost = "127.0.0.1"
		expectedPort = ""
	)
	host, port, err := splitHostPort(hostport)
	if err != nil {
		t.Error(err)
	}
	if host != expectedHost {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedPort, port)
	}
}

func TestSplitHostPort_HostLocalAndPort80(t *testing.T) {
	var hostport = "127.0.0.1:80"
	var (
		expectedHost = "127.0.0.1"
		expectedPort = "80"
	)
	host, port, err := splitHostPort(hostport)
	if err != nil {
		t.Error(err)
	}
	if host != expectedHost {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedPort, port)
	}
}

func TestSplitHostPort_IP6_EmptyHostAndPort80(t *testing.T) {
	var hostport = "[::]:80"
	var (
		expectedHost = "::"
		expectedPort = "80"
	)
	host, port, err := splitHostPort(hostport)
	if err != nil {
		t.Error(err)
	}
	if host != expectedHost {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedPort, port)
	}
}

func TestSplitHostPort_EmptyHostAndPort10074(t *testing.T) {
	var hostport = ":10074"
	var (
		expectedHost = ""
		expectedPort = "10074"
	)
	host, port, err := splitHostPort(hostport)
	if err != nil {
		t.Error(err)
	}
	if host != expectedHost {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("assert 'http.Response.StatusCode':: expected '%v', got '%v'", expectedPort, port)
	}
}
