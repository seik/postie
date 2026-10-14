package config

import (
	"context"
	"net"
	"testing"
)

func TestIPFamilyDialsOnlyThatFamily(t *testing.T) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	port := listener.Addr().(*net.TCPAddr).Port

	v4 := ServerConfigToProvider(ServerConfig{Host: "localhost", Port: port, MaxConnections: 1, IPFamily: 4})
	conn, err := v4.Factory(context.Background())
	if err != nil {
		t.Fatalf("IPv4 dial: %v", err)
	}
	_ = conn.Close()

	v6 := ServerConfigToProvider(ServerConfig{Host: "127.0.0.1", Port: port, MaxConnections: 1, IPFamily: 6})
	if conn, err := v6.Factory(context.Background()); err == nil {
		_ = conn.Close()
		t.Fatal("IPv6 dial reached an IPv4 address")
	}

	if any := ServerConfigToProvider(ServerConfig{Host: "localhost", Port: port, MaxConnections: 1}); any.Factory != nil {
		t.Fatal("no family set should keep nntppool's own dialer")
	}
}
