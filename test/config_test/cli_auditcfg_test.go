package config_test

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/config/config-interpreter"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/config/config-parser"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"os"
	"testing"
)

/*
Ordner: Test1
CLI-Parameter: Keine
Audit-Datei: Wird über CLI gesetzt
Beschreibung: Die Audit-Config, welche über die CLI gesetzt wird, sollte ausgewählt werden
*/
func TestAuditJbaCLI(t *testing.T) {
	// Pfad zu dem Ordner, in welchem wir die Dateien testen wollen
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test1\`, `-auditConfig=./folder_with_audit/cli_audit.jba`}
	cs := LoadConfig(t)
	InterpretConfig(t, &cs)
	Evaluate(t, &cs, gopath+`test\testdata\cli_testdata\Test1\folder_with_audit\cli_audit.jba`)
}

/*
Ordner: Test2
CLI-Parameter: Config wird angegeben
Audit-Datei: Wird über Config gesetzt
Beschreibung: Die Audit-Datei, die in der Config angegeben ist, soll angegeben werden
*/
func TestAuditConfig(t *testing.T) {
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test2\`, `-config=./config.ini`}
	cs := LoadConfig(t)
	InterpretConfig(t, &cs)
	Evaluate(t, &cs, gopath+`test\testdata\cli_testdata\Test2\folder_with_audit\config_audit.jba`)
}

/*
Ordner: Test3
CLI-Parameter: Keine
Audit-Datei: Eine einzelne Audit mit zufälligem Name im Pfad
Beschreibung: Die einzelne Audit-Datei sollte ausgewählt werden
*/
func TestAuditJbaSingleFileInFolder(t *testing.T) {
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test3\`}
	cs := LoadConfig(t)
	InterpretConfig(t, &cs)
	Evaluate(t, &cs, gopath+`test\testdata\cli_testdata\Test3\testauditxy.jba`)
}

/*
Ordner: Test4
CLI-Parameter: Keine
Audit-Datei: Audit-Datei mit Default-Name
Beschreibung: Mehrere .jba-Dateien, eine davon mit Default-Name. Diese sollte ausgewählt werden
*/
func TestAuditJbaDefaultName(t *testing.T) {
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test4\`}
	cs := LoadConfig(t)
	InterpretConfig(t, &cs)
	Evaluate(t, &cs, gopath+`test\testdata\cli_testdata\Test4\audit.jba`)
}

/*
Ordner: Test5
CLI-Parameter: Keine
Audit-Datei: Keine
Beschreibung: Zwei .jba-Dateien, keine mit Default-Name in einem Ordner. Es sollte keine ausgewählt werden.
*/
func TestAuditJbaError(t *testing.T) {
	// Pfad zu dem Ordner, in welchem wir die Dateien testen wollen
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test5\`}
	cs := LoadConfig(t)
	_, err := config_interpreter.InterpretConfig(&cs)

	if err == nil {
		t.Errorf("Fehlgeschlagen: Datei wurde fälschlicherweise genommen: %v", cs.AuditConfig)
	}
}

/*
Ordner: Test6
CLI-Parameter: Keine
Audit-Datei: Keine
Beschreibung: Keine Audit-Datei im Pfad vorhanden.
*/
func TestAuditJbaEmptyFolder(t *testing.T) {
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test6\`}
	cs := LoadConfig(t)
	_, err := config_interpreter.InterpretConfig(&cs)

	if err == nil {
		t.Errorf("Fehlgeschlagen: Datei wurde fälschlicherweise genommen: %v", cs.AuditConfig)
	}
}

/*
Ordner: Test7
CLI-Parameter: Config wird angegeben
Audit-Datei: In Config und CLI gesetzt
Beschreibung: Eine Audit-Datei wird über das CLI und eine in der Config angegeben, die von dem CLI sollte ausgewählt werden.
*/
func TestAuditJbaSelectionConfigCLI(t *testing.T) {
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test7\`, `-auditConfig=./audit1.jba`, `-config=./config.ini`}
	cs := LoadConfig(t)
	InterpretConfig(t, &cs)
	Evaluate(t, &cs, gopath+`test\testdata\cli_testdata\Test7\audit1.jba`)
}

/*
Ordner: Test8
CLI-Parameter: Keine
Audit-Datei: In CLI gesetzt
Beschreibung: Die Audit-Datei welche über das CLI angegeben wurde existiert nicht.
*/
func TestAuditJbaNonExistentAuditCLI(t *testing.T) {
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test8\`, `-auditConfig=./audit.jba`}
	cs := LoadConfig(t)
	_, err := config_interpreter.InterpretConfig(&cs)

	if err == nil {
		t.Errorf("Fehlgeschlagen: Datei wurde fälschlicherweise genommen: %v", cs.AuditConfig)
	}
}

/*
Ordner: Test9
CLI-Parameter: Config wird angegeben
Audit-Datei: In Config gesetzt
Beschreibung: Die Audit-Datei welche über die config angegeben wurde existiert nicht.
*/
func TestAuditJbaNonExistentAuditConfig(t *testing.T) {
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test9\`, `-config=./config.ini`}
	cs := LoadConfig(t)
	_, err := config_interpreter.InterpretConfig(&cs)

	if err == nil {
		t.Errorf("Fehlgeschlagen: Datei wurde fälschlicherweise genommen: %v", cs.AuditConfig)
	}
}

/*
Ordner: Test10
CLI-Parameter: Config wird angegeben
Audit-Datei: In Config gesetzt
Beschreibung: Die Audit-Datei welche über die config angegeben wurde enthält mehrere "/".
*/
func TestAuditJbaTooManySlashes(t *testing.T) {
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test10\`, `-config=./config.ini`}
	cs := LoadConfig(t)
	_, err := config_interpreter.InterpretConfig(&cs)

	if err == nil {
		t.Errorf("Fehlgeschlagen: Datei wurde fälschlicherweise genommen: %v", cs.AuditConfig)
	}
}

/*
Ordner: Test11
CLI-Parameter: Config wird angegeben
Audit-Datei: In Config gesetzt
Beschreibung: Die Audit-Datei welche über die config angegeben wurde existiert nicht.
*/
func TestAuditJbaPathToDirectory(t *testing.T) {
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test11\`, `-config=./config.ini`}
	cs := LoadConfig(t)
	_, err := config_interpreter.InterpretConfig(&cs)

	if err == nil {
		t.Errorf("Fehlgeschlagen: Datei wurde fälschlicherweise genommen: %v", cs.AuditConfig)
	}
}

/*
Ordner: Test12
CLI-Parameter: Config wird angegeben
Audit-Datei: In Config gesetzt
Beschreibung: Der angegebene Pfad enthält keine Audit.jba.
*/
func TestAuditJbaPathWithoutAuditJBA(t *testing.T) {
	os.Args = []string{gopath + `test\testdata\cli_testdata\Test12\`, `-config=./config.ini`}
	cs := LoadConfig(t)
	_, err := config_interpreter.InterpretConfig(&cs)

	if err == nil {
		t.Errorf("Fehlgeschlagen: Datei wurde fälschlicherweise genommen: %v", cs.AuditConfig)
	}
}

//
//
//
// *=*=*=*=*=*=*=*=*=*=*=*=
//
// Utility
//
// *=*=*=*=*=*=*=*=*=*=*=*=
//
//
//
func LoadConfig(t *testing.T) ConfigStruct {
	config_parser.ResetFlags()
	cs, log := config_parser.LoadConfig()
	stringerr := GetLogErr(log)
	if stringerr != "" {
		t.Errorf("Error beim Parsen der Konfiguration: " + stringerr)
	}
	return cs
}

func Evaluate(t *testing.T, cs *ConfigStruct, expectedPath string) {
	if cs.AuditConfig != expectedPath {
		t.Errorf("Fehlgeschlagen: \nExpected=%v\nActual=  %v", expectedPath, cs.AuditConfig)
	}
}
