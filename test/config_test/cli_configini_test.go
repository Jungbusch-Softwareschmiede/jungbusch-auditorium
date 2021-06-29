package config_test

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/config/config-interpreter"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/config/config-parser"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"os"
	s "strings"
	"testing"
)

func TestCLIBasicConfig(t *testing.T) {
	os.Args = []string{gopath, `-config="./test/testdata/cli_testdata/basic_config.ini"`}
	cs, log := config_parser.LoadConfig()
	EvaluateLog(t, log, cs)
	InterpretConfig(t, &cs)
}

func TestCLIMixedOrder(t *testing.T) {
	os.Args = []string{gopath, `-config="./test/testdata/cli_testdata/mixed_order.ini"`}
	cs, log := config_parser.LoadConfig()
	EvaluateLog(t, log, cs)
	InterpretConfig(t, &cs)
}

func TestCLISpacesAndQuotation(t *testing.T) {
	os.Args = []string{gopath, `-config="./test/testdata/cli_testdata/spaces_and_quotation.ini"`}
	cs, log := config_parser.LoadConfig()
	EvaluateLog(t, log, cs)
	InterpretConfig(t, &cs)
}

func TestCLIMissingKey(t *testing.T) {
	os.Args = []string{gopath, `-config="./test/testdata/cli_testdata/missing_key.ini"`}
	cs, log := config_parser.LoadConfig()
	err := GetLogErr(log)

	if !s.HasPrefix(err, "Ungültiger Wert (fehlender Key)") {
		t.Errorf("Fehlgeschlagen. Fehlender Key wurde nicht erkannt.")
	}
	InterpretConfig(t, &cs)
}

func TestCLIMissingValues(t *testing.T) {
	os.Args = []string{gopath, `-config="./test/testdata/cli_testdata/missing_values.ini"`}
	_, log := config_parser.LoadConfig()
	err := GetLogErr(log)
	if !s.HasPrefix(err, "Der Wert eines Keys darf nicht leer sein") {
		t.Errorf("Fehlgeschlagen. Fehlender Wert wurde nicht erkannt.")
	}
}

func TestCLIMissingValuesExceptLoglevel(t *testing.T) {
	os.Args = []string{gopath, `-config="./test/testdata/cli_testdata/missing_values_except_loglevel.ini"`}
	_, log := config_parser.LoadConfig()
	err := GetLogErr(log)
	if !s.HasPrefix(err, "Der Wert eines Keys darf nicht leer sein") {
		t.Errorf("Fehlgeschlagen. Fehlender Wert wurde nicht erkannt.")
	}
}

func TestCLINegativeLoglevel(t *testing.T) {
	os.Args = []string{gopath, `-config="./test/testdata/cli_testdata/negative_loglevel.ini"`}
	cs, log := config_parser.LoadConfig()
	EvaluateLog(t, log, cs)

	_, err := config_interpreter.InterpretConfig(&cs)
	if err != nil {
		t.Errorf("Unbekannter Fehler beim Interpretieren: " + err.Error())
	}

	err = logger.InitializeLogger(&cs, log)
	if err == nil || err.Error() != "Bitte ein Loglevel zwischen 0 und 4 angeben." {
		t.Errorf("Fehler wurde nicht erkannt")
	}
}

func TestCLIStaticPath(t *testing.T) {
	os.Args = []string{gopath, `-config=` + gopath + `/test/testdata/cli_testdata/basic_config.ini`}
	config, log := config_parser.LoadConfig()
	EvaluateLog(t, log, config)
}

func TestEmptyLines(t *testing.T) {
	os.Args = []string{gopath, `-config=./test/testdata/cli_testdata/empty_lines.ini`}
	cs, err := config_parser.LoadConfig()
	strerr := GetLogErr(err)
	if strerr != "" {
		t.Errorf("Fehlgeschlagen: %v, %v", err, cs)
	}
	InterpretConfig(t, &cs)
}

func TestCLIStaticConfig(t *testing.T) {
	os.Args = []string{gopath, `-config=` + getStaticConf("./test/testdata/cli_testdata/static.ini")}
	cs, log := config_parser.LoadConfig()
	EvaluateLog(t, log, cs)
	InterpretConfig(t, &cs)
}
