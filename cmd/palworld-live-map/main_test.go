package main

import "testing"

func TestHealthcheckEndpointUsesListenAddressPort(t *testing.T) {
	tests := []struct {
		addr string
		want string
	}{
		{addr: "", want: "http://127.0.0.1:8080/-/health"},
		{addr: ":9090", want: "http://127.0.0.1:9090/-/health"},
		{addr: "0.0.0.0:7777", want: "http://127.0.0.1:7777/-/health"},
		{addr: "[::]:8123", want: "http://127.0.0.1:8123/-/health"},
		{addr: "localhost:9000", want: "http://localhost:9000/-/health"},
		{addr: "[::1]:8081", want: "http://[::1]:8081/-/health"},
	}
	for _, test := range tests {
		got, err := healthcheckEndpoint(test.addr)
		if err != nil {
			t.Errorf("healthcheckEndpoint(%q) error = %v", test.addr, err)
			continue
		}
		if got != test.want {
			t.Errorf("healthcheckEndpoint(%q) = %q, want %q", test.addr, got, test.want)
		}
	}
}

func TestHealthcheckEndpointRejectsInvalidAddress(t *testing.T) {
	if _, err := healthcheckEndpoint("8080"); err == nil {
		t.Fatal("healthcheckEndpoint() error = nil")
	}
}

func TestHealthcheckEndpointAppliesBasePath(t *testing.T) {
	t.Setenv("BASE_PATH", "/palworld-map/")
	got, err := healthcheckEndpoint(":8080")
	if err != nil {
		t.Fatalf("healthcheckEndpoint() error = %v", err)
	}
	if want := "http://127.0.0.1:8080/palworld-map/-/health"; got != want {
		t.Fatalf("healthcheckEndpoint() = %q, want %q", got, want)
	}
}

func TestHealthcheckEndpointRejectsInvalidBasePath(t *testing.T) {
	t.Setenv("BASE_PATH", "/")
	if _, err := healthcheckEndpoint(":8080"); err == nil {
		t.Fatal("healthcheckEndpoint() error = nil")
	}
}
