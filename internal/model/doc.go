package model

//go:generate charong
//go:generate mockery -all -inpkg -output_file=mocks.go
//go:generate goimports -w schema.pqt.go mocks.go
