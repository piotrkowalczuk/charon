package main

import (
	"github.com/piotrkowalczuk/charon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type isGrantedHandler struct {
	*handler
}

func (ig *isGrantedHandler) handle(ctx context.Context, r *charon.IsGrantedRequest) (*charon.IsGrantedResponse, error) {
	//	ig.repository.group.FindByUserID()
	return nil, grpc.Errorf(codes.Unimplemented, "is granted is not implemented yet")
}

func (ig *isGrantedHandler) firewall() {

}
