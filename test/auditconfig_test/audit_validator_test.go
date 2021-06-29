package auditconfig_test

import (
	interpreter2 "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/interpreter"
	"os"
	"testing"
)

var (
	gopath = os.Getenv("gopath") + `\src\github.com\Jungbusch-Softwareschmiede\jungbusch-auditorium\`
)

func TestValidateParametersExistingVariable(t *testing.T) {
	loadedModules := initMe()

	m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/validator_testdata/existing_variable.jba")
	if err != nil {
		t.Errorf("Fehler im Parser: %v", err)
	}

	if err := interpreter2.ValidateParameters(m); err != nil {
		t.Errorf("Ungültige Datei wurde angenommen %v", err)
	}
}

func TestValidateParametersMultipleVariables(t *testing.T) {
	loadedModules := initMe()

	m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/validator_testdata/multiple_variables.jba")
	if err != nil {
		t.Errorf("Fehler im Parser: %v", err)
	}

	if err := interpreter2.ValidateParameters(m); err != nil {
		t.Errorf("Ungültige Datei wurde angenommen %v", err)
	}
}

func TestValidateParametersNonExistingVariable(t *testing.T) {
	loadedModules := initMe()

	m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/validator_testdata/non_existing_variable.jba")
	if err != nil {
		t.Errorf("Unbekannter Fehler: " + err.Error())
	}

	if err := interpreter2.ValidateParameters(m); err == nil {
		t.Errorf("Ungültige Datei wurde angenommen %v", err)
	}
}

func TestValidateParametersMultipleNonExistingVariables(t *testing.T) {
	loadedModules := initMe()

	m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/validator_testdata/multiple_non_existing_variables.jba")
	if err != nil {
		t.Errorf("Fehler im Parser: %v", err)
	}

	if err := interpreter2.ValidateParameters(m); err == nil {
		t.Errorf("Ungültige Datei wurde angenommen %v", err)
	}
}
