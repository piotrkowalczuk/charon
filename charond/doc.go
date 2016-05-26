// Package charond ...
package charond

//go:generate charong
//go:generate mockery -all -inpkg -output_file=mocks_test.go
//go:generate goimports -w schema.pqt.go mocks_test.go
