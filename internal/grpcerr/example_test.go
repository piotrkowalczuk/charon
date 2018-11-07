package grpcerr_test

import (
	"fmt"

	"github.com/piotrkowalczuk/charon/internal/grpcerr"
)

func ExampleError() {
	// Single error.
	e1 := grpcerr.E(grpcerr.Op("Get"), grpcerr.Kind("io"), "network unreachable")
	fmt.Println("\nSimple error:")
	fmt.Println(e1)
	// Nested error.
	fmt.Println("\nNested error:")
	e2 := grpcerr.E(grpcerr.Op("Read"), grpcerr.Kind("other"), e1)
	fmt.Println(e2)
	// Output:
	//
	// Simple error:
	// Get: io: network unreachable
	//
	// Nested error:
	// Read: io:
	//	Get: network unreachable
}
func ExampleMatch() {
	err := grpcerr.E("network unreachable")
	// Construct an error, one we pretend to have received from a test.
	got := grpcerr.E(grpcerr.Op("Get"), grpcerr.Kind("io"), err)
	// Now construct a reference error, which might not have all
	// the fields of the error from the test.
	expect := grpcerr.E(grpcerr.Kind("io"), err)
	fmt.Println("Match:", grpcerr.Match(expect, got))
	// Now one that's incorrect - wrong Kind.
	got = grpcerr.E(grpcerr.Op("Get"), grpcerr.Kind("permission"), err)
	fmt.Println("Mismatch:", grpcerr.Match(expect, got))
	// Output:
	//
	// Match: true
	// Mismatch: false
}
