package translation

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type transMap map[string]string
type langMap map[string]transMap

// CSVSource ...
type CSVSource struct {
	langs langMap
}

// Get ...
func (c *CSVSource) Get(lang string, key string) (string, error) {
	var translations transMap

	translations, ok := c.langs[lang]

	if !ok {
		return "", ErrNoTranslation
	}

	text, ok := translations[key]

	if !ok {
		return "", ErrNoTranslation
	}

	return text, nil
}

func parseCSVFile(csvFile io.Reader, lang string) *transMap {
	r := csv.NewReader(csvFile)
	records, err := r.ReadAll()

	if err != nil {
		log.Fatalln(err)
	}

	tm := make(transMap, len(records))

	for _, row := range records {
		if len(row) != 2 {
			log.Fatalln("Translation CSV file needs to have exactly 2 entries per line")
		}

		tm[row[0]] = row[1]
	}

	return &tm
}

// NewCSVSource ...
func NewCSVSource(translationsDir string) *CSVSource {
	dir := "data/i18n"
	fileInfos, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Fatalln(err)
	}

	lm := make(langMap, len(fileInfos))

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			lang := strings.TrimSuffix(fileInfo.Name(), filepath.Ext(fileInfo.Name()))
			file, err := os.Open(dir + "/" + fileInfo.Name())

			if err != nil {
				log.Println(err)
			}

			tm := *parseCSVFile(file, lang)
			lm[lang] = tm
		}
	}

	s := &CSVSource{
		langs: lm,
	}

	return s
}
