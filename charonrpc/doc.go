package charonrpc

//go:generate protoc -I=. -I=../vendor -I=${GOPATH}/src --go_out=plugins=grpc:. auth.proto user.proto group.proto permission.proto
//go:generate python -m grpc_tools.protoc -I=. -I=/usr/include -I=../vendor --python_out=. --grpc_python_out=. auth.proto user.proto group.proto permission.proto
//go:generate goimports -w auth.pb.go
