package service

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// Mnemosyne ...
var Mnemosyne *rpc.Client

// MnemosyneConfig ...
type MnemosyneConfig struct {
	Host string `xml:"host"`
	Port string `xml:"port"`
}

func InitMnemosyne(config MnemosyneConfig) {
	conn, err := net.Dial("tcp", config.Host+":"+config.Port)
	if err != nil {
		Logger.Fatal(err)
	}

	Mnemosyne = jsonrpc.NewClient(conn)

	Logger.Info("Connection do Mnemosyne established successfully.")
}
