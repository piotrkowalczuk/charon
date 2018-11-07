package charond

import (
	"net"
	"testing"

	"go.uber.org/zap"

	"github.com/piotrkowalczuk/mnemosyne/mnemosyned"
	"golang.org/x/crypto/bcrypt"
)

func TestDaemon_Run(t *testing.T) {
	certPath := "../../data/test-selfsigned.crt"
	keyPath := "../../data/test-selfsigned.key"

	mnemosyneRPCListener, err := net.Listen("tcp", "localhost:0") // any available address
	if err != nil {
		t.Fatalf("mnemosyne daemon tcp listener setup error: %s", err.Error())
	}
	mnemosyneDaemon, err := mnemosyned.NewDaemon(&mnemosyned.DaemonOpts{
		IsTest:            true,
		ClusterListenAddr: mnemosyneRPCListener.Addr().String(),
		Logger:            zap.L(),
		PostgresAddress:   testPostgresAddress,
		PostgresTable:     "session",
		PostgresSchema:    "mnemosyne",
		RPCListener:       mnemosyneRPCListener,
		TLS:               true,
		TLSKeyFile:        keyPath,
		TLSCertFile:       certPath,
	})
	if err != nil {
		t.Fatalf("mnemosyne daemon cannot be instantiated: %s", err.Error())
	}
	if err := mnemosyneDaemon.Run(); err != nil {
		t.Fatalf("mnemosyne daemon start error: %s", err.Error())
	}

	charonRPCListener, err := net.Listen("tcp", "localhost:0") // any available address
	if err != nil {
		t.Fatalf("charon daemon tcp listener setup error: %s", err.Error())
	}
	debugRPCListener, err := net.Listen("tcp", "localhost:0") // any available address
	if err != nil {
		t.Fatalf("charon daemon tcp listener setup error: %s", err.Error())
	}

	logger := zap.L()

	d := NewDaemon(DaemonOpts{
		Test:                 true,
		Monitoring:           false,
		MnemosyneAddress:     mnemosyneDaemon.Addr().String(),
		MnemosyneTLS:         true,
		MnemosyneTLSCertFile: certPath,
		Logger:               logger,
		PostgresAddress:      testPostgresAddress,
		RPCListener:          charonRPCListener,
		DebugListener:        debugRPCListener,
		PasswordBCryptCost:   bcrypt.MinCost,
		TLS:                  true,
		TLSKeyFile:           keyPath,
		TLSCertFile:          certPath,
	})
	if err := d.Run(); err != nil {
		t.Fatalf("charon daemon start error: %s", err.Error())
	}
	if err := d.Close(); err != nil {
		t.Fatalf("charon daemon close error: %s", err.Error())
	}
}
