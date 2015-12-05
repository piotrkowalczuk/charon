package charontest_test

import (
	"testing"

	"github.com/piotrkowalczuk/charon"
	"github.com/piotrkowalczuk/charon/charontest"
)

func TestCharon(t *testing.T) {
	var mock interface{} = &charontest.Charon{}

	if _, ok := mock.(charon.Charon); !ok {
		t.Errorf("mock should implement original interface, but does not")
	}
}
