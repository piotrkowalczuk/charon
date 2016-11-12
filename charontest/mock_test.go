package charontest_test

import (
	"testing"

	"github.com/piotrkowalczuk/charon/charonc"
	"github.com/piotrkowalczuk/charon/charontest"
)

func TestCharon(t *testing.T) {
	var mock interface{} = &charontest.Client{}

	if _, ok := mock.(charonc.Client); !ok {
		t.Error("mock should implement original interface, but does not")
	}
}
