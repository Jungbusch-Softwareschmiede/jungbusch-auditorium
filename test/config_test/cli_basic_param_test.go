package config_test

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/config/config-interpreter"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/config/config-parser"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestAuditConfigFlag(t *testing.T) {
	config_parser.ResetFlags()
	pfad := "./test/testdata/cli_testdata/auditDummy2.jba"
	os.Args = []string{gopath, "-auditConfig=" + pfad}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.AuditConfig != pfad {
		t.Errorf("Audit-Config-Pfad wurde nicht richtig gesetzt:\nIst:  " + cs.AuditConfig + "\nSoll: " + pfad)
	}
}

func TestAuditConfigFlag2(t *testing.T) {
	config_parser.ResetFlags()
	pfad := "./test/testdata/cli_testdata/auditDummy2.jba"
	os.Args = []string{gopath, "-a=" + pfad}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.AuditConfig != pfad {
		t.Errorf("Audit-Config-Pfad wurde nicht richtig gesetzt:\nIst:  " + cs.AuditConfig + "\nSoll: " + pfad)
	}
}

func TestConfigFlag(t *testing.T) {
	os.Args = []string{gopath, "-config=./test/testdata/cli_testdata/basic_config.ini"}
	cs, log := config_parser.LoadConfig()
	pfad := filepath.Clean(gopath + "/test/testdata/cli_testdata/basic_config.ini")

	_, err := config_interpreter.InterpretConfig(&cs)
	if err != nil {
		t.Errorf("Fehler beim Interpretieren.")
	}

	EvaluateLog(t, log, cs)
	if cs.Config != pfad {
		t.Errorf("Config-Pfad wurde nicht richtig gesetzt:\nIst:  " + cs.Config + "\nSoll: " + pfad)
	}
}

func TestConfigFlag2(t *testing.T) {
	os.Args = []string{gopath, "-c=./test/testdata/cli_testdata/basic_config.ini"}
	cs, log := config_parser.LoadConfig()
	pfad := filepath.Clean(gopath + "/test/testdata/cli_testdata/basic_config.ini")

	EvaluateLog(t, log, cs)
	if cs.Config != pfad {
		t.Errorf("Config-Pfad wurde nicht richtig gesetzt:\nIst:  " + cs.Config + "\nSoll: " + pfad)
	}
}

func TestOutputFlag(t *testing.T) {
	pfad := "./test/testdata/cli_testdata/"
	os.Args = []string{gopath, "-outputPath=" + pfad}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.OutputPath != pfad {
		t.Errorf("Output-Pfad wurde nicht richtig gesetzt:\nIst:  " + cs.OutputPath + "\nSoll: " + pfad)
	}
}

func TestOutputFlag2(t *testing.T) {
	pfad := "./test/testdata/cli_testdata/"
	os.Args = []string{gopath, "-o=" + pfad}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.OutputPath != pfad {
		t.Errorf("Output-Pfad wurde nicht richtig gesetzt:\nIst:  " + cs.OutputPath + "\nSoll: " + pfad)
	}
}

func TestVerbosityLogFlag(t *testing.T) {
	os.Args = []string{gopath, "-verbosityLog=2"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.VerbosityLog != 2 {
		t.Errorf("Log-Verbosity wurde nicht richtig gesetzt:\nIst:  " + strconv.Itoa(cs.VerbosityLog) + "\nSoll: " + "2")
	}
}

func TestVerbosityLogFlag2(t *testing.T) {
	os.Args = []string{gopath, "-vl=2"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.VerbosityLog != 2 {
		t.Errorf("Log-Verbosity wurde nicht richtig gesetzt:\nIst:  " + strconv.Itoa(cs.VerbosityLog) + "\nSoll: " + "2")
	}
}

func TestVerbosityConsoleFlag(t *testing.T) {
	os.Args = []string{gopath, "-verbosityConsole=3"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.VerbosityConsole != 3 {
		t.Errorf("Log-Verbosity wurde nicht richtig gesetzt:\nIst:  " + strconv.Itoa(cs.VerbosityConsole) + "\nSoll: " + "3")
	}
}

func TestVerbosityConsoleFlag2(t *testing.T) {
	os.Args = []string{gopath, "-vc=3"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.VerbosityConsole != 3 {
		t.Errorf("Log-Verbosity wurde nicht richtig gesetzt:\nIst:  " + strconv.Itoa(cs.VerbosityConsole) + "\nSoll: " + "3")
	}
}

func TestSkipModuleCompatibilityFlag(t *testing.T) {
	os.Args = []string{gopath, "-skipModuleCompatibilityCheck"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.SkipModuleCompatibilityCheck {
		t.Errorf("SkipModuleCompatibilityCheck wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestSkipModuleCompatibilityFlag2(t *testing.T) {
	os.Args = []string{gopath, "-skipModuleCompatibilityCheck=false"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.SkipModuleCompatibilityCheck {
		t.Errorf("SkipModuleCompatibilityCheck wurde nicht richtig gesetzt:\nIst:  true\nSoll: false")
	}
}

func TestKeepConsoleOpenFlag(t *testing.T) {
	os.Args = []string{gopath, "-keepConsoleOpen"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.KeepConsoleOpen {
		t.Errorf("KeepConsoleOpen wurde nicht richtig gesetzt:\nIst:  true\nSoll: false")
	}
}

func TestForceOSFlag(t *testing.T) {
	os.Args = []string{gopath, "-forceOS=\"testOS\""}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.ForceOS != "testOS" {
		t.Errorf("forceOS wurde nicht richtig gesetzt:\nIst:  " + cs.ForceOS + "\nSoll: testOS")
	}
}

func TestForceOSFlag2(t *testing.T) {
	os.Args = []string{gopath, "-forceOS=testOS"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.ForceOS != "testOS" {
		t.Errorf("forceOS wurde nicht richtig gesetzt:\nIst:  " + cs.ForceOS + "\nSoll: testOS")
	}
}

func TestForceOSFlag3(t *testing.T) {
	os.Args = []string{gopath, "-forceOS= testOS"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.ForceOS != "testOS" {
		t.Errorf("forceOS wurde nicht richtig gesetzt:\nIst:  " + cs.ForceOS + "\nSoll: testOS")
	}
}

func TestForceOSFlag4(t *testing.T) {
	os.Args = []string{gopath, "--forceOS=testOS"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.ForceOS != "testOS" {
		t.Errorf("forceOS wurde nicht richtig gesetzt:\nIst:  " + cs.ForceOS + "\nSoll: testOS")
	}
}

func TestIgnoreMissingPrivilegesFlag(t *testing.T) {
	os.Args = []string{gopath, "-ignoreMissingPrivileges"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.IgnoreMissingPrivileges {
		t.Errorf("IgnoreMissingPrivileges wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestIgnoreMissingPrivilegesFlag2(t *testing.T) {
	os.Args = []string{gopath, "-ignoreMissingPrivileges=true"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.IgnoreMissingPrivileges {
		t.Errorf("IgnoreMissingPrivileges wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestAlwaysPrintProgressFlag(t *testing.T) {
	os.Args = []string{gopath, "-alwaysPrintProgress"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.AlwaysPrintProgress {
		t.Errorf("AlwaysPrintProgress wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestZipFlag(t *testing.T) {
	os.Args = []string{gopath, "-zip"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.Zip {
		t.Errorf("Zip wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestZipFlag2(t *testing.T) {
	os.Args = []string{gopath, "-zip=true"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.Zip {
		t.Errorf("Zip wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestZipOnlyFlag(t *testing.T) {
	os.Args = []string{gopath, "-zipOnly"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.ZipOnly {
		t.Errorf("ZipOnly wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestZipOnlyFlag2(t *testing.T) {
	os.Args = []string{gopath, "-zipOnly=true"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.ZipOnly {
		t.Errorf("ZipOnly wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestZipAndZipOnlyFlag(t *testing.T) {
	os.Args = []string{gopath, "-zipOnly", "-zip"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if _, err := config_interpreter.InterpretConfig(&cs); err == nil {
		t.Errorf("Zip und ZipOnly wurden gleichzeitig gesetzt, der Fehler wurde nicht erkannt.")
	}
}

func TestVersionFlag(t *testing.T) {
	os.Args = []string{gopath, "-version"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.Version {
		t.Errorf("Version wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestShowModulesFlag(t *testing.T) {
	os.Args = []string{gopath, "-showModules"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.ShowModule != "all" {
		t.Errorf("ShowModule wurde nicht richtig gesetzt:\nIst:  " + cs.ShowModule + "\nSoll: all")
	}
}

func TestShowModulesFlag2(t *testing.T) {
	os.Args = []string{gopath, "-showModuleInfo=test"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if cs.ShowModule != "test" {
		t.Errorf("ShowModule wurde nicht richtig gesetzt:\nIst:  " + cs.ShowModule + "\nSoll: test")
	}
}

func TestCheckConfigurationFlag(t *testing.T) {
	os.Args = []string{gopath, "-checkConfiguration"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.CheckConfiguration {
		t.Errorf("CheckConfiguration wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestCheckSyntaxFlag(t *testing.T) {
	os.Args = []string{gopath, "-checkSyntax"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.CheckSyntax {
		t.Errorf("CheckSyntax wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestCheckSyntaxFlag2(t *testing.T) {
	os.Args = []string{gopath, "-syntax"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.CheckSyntax {
		t.Errorf("CheckSyntax wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestSaveConfigurationFlag(t *testing.T) {
	os.Args = []string{gopath, "-saveConfiguration"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.SaveConfiguration {
		t.Errorf("Save wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestSaveConfigurationFlag2(t *testing.T) {
	os.Args = []string{gopath, "-s"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.SaveConfiguration {
		t.Errorf("Save wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}

func TestCreateDefaultConfigFlag(t *testing.T) {
	os.Args = []string{gopath, "-createDefault"}
	cs, log := config_parser.LoadConfig()

	EvaluateLog(t, log, cs)
	if !cs.CreateDefaultConfig {
		t.Errorf("CreateDefaultConfig wurde nicht richtig gesetzt:\nIst:  false\nSoll: true")
	}
}
