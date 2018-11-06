package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/piotrkowalczuk/pqt/pqtgo/pqtgogen"

	"github.com/piotrkowalczuk/ntypespqt"
	"github.com/piotrkowalczuk/pqt/pqtsql"
	"github.com/piotrkowalczuk/qtypespqt"
)

var (
	schema, output string
)

func init() {
	flag.StringVar(&schema, "schema", "charon", "")
	flag.StringVar(&output, "output", "schema.pqt", "")
}

func main() {
	file, err := openFile(output)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sch := databaseSchema()
	genGo := &pqtgogen.Generator{
		Pkg:        "model",
		Version:    9.5,
		Components: pqtgogen.ComponentAll,
		Plugins: []pqtgogen.Plugin{
			&qtypespqt.Plugin{},
			&ntypespqt.Plugin{},
		},
	}
	genSQL := &pqtsql.Generator{}
	if err := genGo.GenerateTo(sch, file); err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(file, "const SQL = `\n")
	if err := genSQL.GenerateTo(sch, file); err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(file, "`")
}

func openFile(output string) (io.WriteCloser, error) {
	return os.OpenFile(output+".go", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0660)
}
