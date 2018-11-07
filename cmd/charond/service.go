package main

import (
	"net"
	"strconv"

	"go.uber.org/zap"
)

func initListener(logger *zap.Logger, host string, port int) net.Listener {
	on := host + ":" + strconv.FormatInt(int64(port), 10)
	listener, err := net.Listen("tcp", on)
	if err != nil {
		logger.Fatal("listener initialization failure", zap.Error(err), zap.String("host", host), zap.Int("port", port))
	}
	return listener
}
