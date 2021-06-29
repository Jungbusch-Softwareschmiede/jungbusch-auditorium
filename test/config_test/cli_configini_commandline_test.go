package config_test

import (
	"fmt"
	config_parser "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/config/config-parser"
	"os"
	"strconv"
	"testing"
)

func TestCmdOverwriteAuditConfig(t *testing.T) {
	os.Args = []string{gopath, `-config="./test/testdata/cli_testdata/overwrite_audconf.ini"`, `-a="test2"`}
	cs, log := config_parser.LoadConfig()
	EvaluateLog(t, log, cs)

	if cs.AuditConfig != "test2" {
		t.Error("Parameter wurde nicht korrekt überschrieben: " + cs.AuditConfig)
	}
}

func TestCmdOverwriteOutput(t *testing.T) {
	os.Args = []string{gopath, `-config="./test/testdata/cli_testdata/overwrite_output.ini"`, `-o="test2"`}
	cs, log := config_parser.LoadConfig()
	EvaluateLog(t, log, cs)

	if cs.OutputPath != "test2" {
		t.Error("Parameter wurde nicht korrekt überschrieben: " + cs.AuditConfig)
	}
}

func TestCmdOverwriteVerbosity(t *testing.T) {
	os.Args = []string{gopath, `-config="./test/testdata/cli_testdata/overwrite_verbosity.ini"`, `-verbosityConsole=4`, `-verbosityLog=4`}
	cs, log := config_parser.LoadConfig()
	fmt.Println(log)
	EvaluateLog(t, log, cs)

	if cs.VerbosityLog != 4 || cs.VerbosityConsole != 4 {
		t.Error("Parameter wurde nicht korrekt überschrieben: " + strconv.Itoa(cs.VerbosityLog) + " " + strconv.Itoa(cs.VerbosityConsole))
	}
}

func TestCmdOverwriteBooleans(t *testing.T) {
	os.Args = []string{gopath, `-config="./test/testdata/cli_testdata/overwrite_booleans.ini"`, `-skipModuleCompatibilityCheck=false`, `-keepConsoleOpen=false`, "-ignoreMissingPrivileges=false", "-alwaysPrintProgress=false", "-zip=false", "-zipOnly=false"}
	cs, log := config_parser.LoadConfig()
	fmt.Println(log)
	EvaluateLog(t, log, cs)

	if cs.SkipModuleCompatibilityCheck || cs.KeepConsoleOpen || cs.IgnoreMissingPrivileges || cs.AlwaysPrintProgress || cs.Zip || cs.ZipOnly {
		t.Error("Einer der Booleans wurde nicht korrekt überschrieben: " + fmt.Sprint(cs.SkipModuleCompatibilityCheck) + fmt.Sprint(cs.KeepConsoleOpen) + fmt.Sprint(cs.IgnoreMissingPrivileges) + fmt.Sprint(cs.AlwaysPrintProgress) + fmt.Sprint(cs.Zip) + fmt.Sprint(cs.ZipOnly))
	}
}
