package charond

import (
	"flag"
	"net"
	"os"
	"testing"
)

var (
	testPostgresAddress string
)

func init() {
	flag.StringVar(&testPostgresAddress, "p.address", "postgres://postgres:@localhost/test?sslmode=disable", "")
}

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

type suite interface {
	setup(testing.T)
	teardown(testing.T)
}

func listenTCP(t *testing.T) net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0") // any available address
	if err != nil {
		t.Fatalf("net.Listen tcp :0: %s", err.Error())
	}
	return l
}
