package password_test

import (
	"bytes"
	"testing"

	"github.com/piotrkowalczuk/charon/internal/password"

	"golang.org/x/crypto/bcrypt"
)

func TestNewBCryptHasher_bellowMin(t *testing.T) {
	ph, err := password.NewBCryptHasher(bcrypt.MinCost - 1)

	testNewBCryptHasherFailure(t, ph, err)
}

func TestNewBCryptHasher_min(t *testing.T) {
	ph, err := password.NewBCryptHasher(bcrypt.MinCost)

	testNewBCryptHasherSuccess(t, ph, err)
}

func TestNewBCryptHasher_between(t *testing.T) {
	ph, err := password.NewBCryptHasher(25)

	testNewBCryptHasherSuccess(t, ph, err)
}

func TestNewBCryptHasher_max(t *testing.T) {
	ph, err := password.NewBCryptHasher(bcrypt.MaxCost)

	testNewBCryptHasherSuccess(t, ph, err)
}

func TestNewBCryptHasher_aboveMax(t *testing.T) {
	ph, err := password.NewBCryptHasher(bcrypt.MaxCost + 1)

	testNewBCryptHasherFailure(t, ph, err)
}

func testNewBCryptHasherSuccess(t *testing.T, ph password.Hasher, err error) {
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	if ph == nil {
		t.Error("password hasher should not be nil")
	}
}

func testNewBCryptHasherFailure(t *testing.T, ph password.Hasher, err error) {
	if err == nil {
		t.Error("error expected")
	}
	if ph != nil {
		t.Errorf("password hasher should be nil, but got %v", ph)
	}
}

func TestBCryptHasher_Hash(t *testing.T) {
	ph, err := password.NewBCryptHasher(10)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	given := []byte("123")
	expected := []byte("$2a$10$NF5jon4vHytVzwVz5wKAe.AycwRQ8mmeRXEoxTMu4kh4He7K1YCRe")

	got, err := ph.Hash(given)
	if err != nil {
		t.Fatalf("hash returned unexpected error: %s", err.Error())
	}
	if len(expected) != len(got) {
		t.Errorf("length of hash do not match, expected %d but got %d", len(expected), len(got))
	}
	if !bytes.HasPrefix(got, expected[:7]) {
		t.Errorf("hash should have prefix %s but does not: %s", expected[:7], got)
	}
}

func TestBCryptHasher_Compare(t *testing.T) {
	given := []byte("123")

	ph, err := password.NewBCryptHasher(10)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	got, err := ph.Hash(given)
	if err != nil {
		t.Fatalf("hash returned unexpected error: %s", err.Error())
	}

	if !ph.Compare(got, given) {
		t.Error("password do not match")
	}
}
