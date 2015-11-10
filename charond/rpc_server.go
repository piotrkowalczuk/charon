package main

import "github.com/go-kit/kit/log"

type rpcServer struct {
	logger  log.Logger
	monitor *monitoring
}
