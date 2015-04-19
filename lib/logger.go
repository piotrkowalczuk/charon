package lib

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	nocolor = 0
	red     = 91
	green   = 92
	yellow  = 93
	cyan    = 96
)

var (
	baseTimestamp time.Time
	isTerminal    bool
)

func init() {
	baseTimestamp = time.Now()
	isTerminal = logrus.IsTerminal()
}

// ConsoleFormatter ...
type ConsoleFormatter struct {
	ForceColors   bool
	DisableColors bool
}

// Format ...
func (cf *ConsoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	isColored := (cf.ForceColors || isTerminal) && !cf.DisableColors
	b := &bytes.Buffer{}

	var keys []string
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if isColored {
		printColored(b, entry, keys)
	} else {
		b.WriteString("[ ")
		b.WriteString(entry.Time.Format(time.RFC3339))
		b.WriteString(" ] ")
		b.WriteString(entry.Message)
		b.WriteString(" ")
		for _, key := range keys {
			fmt.Fprintf(b, "%v=%v ", key, entry.Data[key])
		}
		b.WriteString("\n")
	}

	return b.Bytes(), nil
}

func printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string) {
	var levelColor int
	switch entry.Level {
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	case logrus.DebugLevel:
		levelColor = 95
	default:
		levelColor = cyan
	}

	levelText := strings.ToUpper(entry.Level.String())[0:4]

	fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %-44s", entry.Time.Format(time.RFC3339), levelColor, levelText, entry.Message)
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=%v", levelColor, k, v)
	}
	b.WriteString("\n")
}
