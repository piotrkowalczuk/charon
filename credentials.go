package charon

import (
	"fmt"

	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
)

type Credentials struct {
	transportSecurity bool
}

// NewCredentials allocates new Credentials object.
func NewCredentials(ts bool) credentials.Credentials {
	return Credentials{
		transportSecurity: ts,
	}
}

// GetRequestMetadata implements credentials.Credentials interface.
func (c Credentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	at, ok := mnemosyne.AccessTokenFromContext(ctx)
	fmt.Println("GetRequestMetadata", at, ok)
	if ok {
		return map[string]string{
			mnemosyne.AccessTokenMetadataKey: "Bearer " + at.Encode(),
		}, nil
	}
	// maybe someone already set metadata previously
	return nil, nil
}

// RequireTransportSecurity implements credentials.Credentials interface.
func (c Credentials) RequireTransportSecurity() bool {
	return c.transportSecurity
}
