// Package charon ...
package charon

//go:generate protoc -I=. -I=./vendor -I=${GOPATH}/src --go_out=plugins=grpc:. charon.proto
//go:generate goimports -w charon.pb.go
