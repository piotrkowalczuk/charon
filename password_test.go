// +build unit,!postgres,!e2e

package charon

import (
	"bytes"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewBCryptPasswordHasher_bellowMin(t *testing.T) {
	ph, err := NewBCryptPasswordHasher(bcrypt.MinCost - 1)

	testNewBCryptPasswordHasher_failure(t, ph, err)
}

func TestNewBCryptPasswordHasher_min(t *testing.T) {
	ph, err := NewBCryptPasswordHasher(bcrypt.MinCost)

	testNewBCryptPasswordHasher_success(t, ph, err)
}

func TestNewBCryptPasswordHasher_between(t *testing.T) {
	ph, err := NewBCryptPasswordHasher(25)

	testNewBCryptPasswordHasher_success(t, ph, err)
}

func TestNewBCryptPasswordHasher_max(t *testing.T) {
	ph, err := NewBCryptPasswordHasher(bcrypt.MaxCost)

	testNewBCryptPasswordHasher_success(t, ph, err)
}

func TestNewBCryptPasswordHasher_aboveMax(t *testing.T) {
	ph, err := NewBCryptPasswordHasher(bcrypt.MaxCost + 1)

	testNewBCryptPasswordHasher_failure(t, ph, err)
}

func testNewBCryptPasswordHasher_success(t *testing.T, ph PasswordHasher, err error) {
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	if ph == nil {
		t.Errorf("password hasher should not be nil")
	}
}

func testNewBCryptPasswordHasher_failure(t *testing.T, ph PasswordHasher, err error) {
	if err == nil {
		t.Error("error expected")
	}
	if ph != nil {
		t.Errorf("password hasher should be nil, but got %v", ph)
	}
}

func TestBCryptPasswordHasher_Hash(t *testing.T) {
	ph, err := NewBCryptPasswordHasher(10)
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
		t.Errorf("length of hash do not match", len(expected), len(got))
	}
	if !bytes.HasPrefix(got, expected[:7]) {
		t.Errorf("hash should have prefix %s but does not: %s", expected[:7], got)
	}
}

func TestBCryptPasswordHasher_Compare(t *testing.T) {
	given := []byte("123")

	ph, err := NewBCryptPasswordHasher(10)
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
