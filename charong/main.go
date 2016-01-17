package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/piotrkowalczuk/pqt/pqtgo"
	"github.com/piotrkowalczuk/pqt/pqtsql"
)

var (
	schema, output string
	acronyms       = map[string]string{
		"id":   "ID",
		"http": "HTTP",
		"ip":   "IP",
		"net":  "NET",
		"irc":  "IRC",
		"io":   "IO",
		"kpi":  "KPI",
		"ssp":  "SSP",
		"uuid": "UUID",
	}
)

func init() {
	flag.StringVar(&schema, "schema", "charon", "")
	flag.StringVar(&output, "output", "schema", "")
}

func main() {
	sqlFile, goFile, err := openFiles(output)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlFile.Close()
	defer goFile.Close()

	sch := databaseSchema()
	if err := pqtsql.NewGenerator().GenerateTo(sch, sqlFile); err != nil {
		log.Fatal(err)
	}
	if err := pqtgo.NewGenerator().
		AddImport("github.com/piotrkowalczuk/nilt").
		SetAcronyms(acronyms).
		GenerateTo(sch, goFile); err != nil {
		log.Fatal(err)
	}

	log.Println("success")
}

func openFiles(output string) (sqlFile io.WriteCloser, goFile io.WriteCloser, err error) {
	sqlFile, err = os.OpenFile(output+".sql", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return
	}
	goFile, err = os.OpenFile(output+".go", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0660)
	return
}
