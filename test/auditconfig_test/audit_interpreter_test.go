package auditconfig_test

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/interpreter"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"testing"
)

func initVarMap() VariableMap {
	varSlice := make(VariableMap)
	varSlice["%result%"] = Variable{Name: "%result%", IsEnv: true}
	varSlice["%passed%"] = Variable{Name: "%passed%", IsEnv: true}
	varSlice["%unsuccessful%"] = Variable{Name: "%unsuccessful%", IsEnv: true}
	varSlice["%os%"] = Variable{Name: "%os%", Value: static.OperatingSystem, IsEnv: true}
	varSlice["%currentmodule%"] = Variable{Name: "%currentmodule%", IsEnv: true}
	return varSlice
}

func TestInterpreterBasicPassed(t *testing.T) {
	_ = initMe()

	myModule := AuditModule{
		ModuleName:       "ExecuteCommand",
		StepID:           "1",
		Passed:           "%result%.startsWith(\"Te\")",
		Variables:        initVarMap(),
		ModuleParameters: ParameterMap{"command": "echo Test"},
	}

	report, err := interpreter.InterpretAudit([]AuditModule{myModule}, 1, false)

	if err != nil {
		t.Errorf("Fehler aufgetreten: " + err.Error())
	}

	if report[0].Result != "PASSED" {
		t.Errorf("Result: " + report[0].Result + "\nSoll: PASSED")
	}

	t.Log("Success")
}

func TestInterpreterBasicNotPassed(t *testing.T) {
	_ = initMe()

	myModule := AuditModule{
		ModuleName:       "ExecuteCommand",
		StepID:           "1",
		Passed:           "%result%.startsWith(\"Hallo\")",
		Variables:        initVarMap(),
		ModuleParameters: ParameterMap{"command": "echo Test"},
	}

	report, err := interpreter.InterpretAudit([]AuditModule{myModule}, 1, false)

	if err != nil {
		t.Errorf("Fehler aufgetreten: " + err.Error())
	}

	if report[0].Result != "NOTPASSED" {
		t.Errorf("Result: " + report[0].Result + "\nSoll: NOTPASSED")
	}

	t.Log("Success")
}

func TestInterpreterBasicUnsuccessful(t *testing.T) {
	_ = initMe()

	myModule := AuditModule{
		ModuleName:       "FileContent",
		StepID:           "1",
		Variables:        initVarMap(),
		ModuleParameters: ParameterMap{"file": "..."},
	}

	report, err := interpreter.InterpretAudit([]AuditModule{myModule}, 1, false)

	if err != nil {
		t.Errorf("Fehler aufgetreten: " + err.Error())
	}

	if report[0].Result != "UNSUCCESSFUL" {
		t.Errorf("Result: " + report[0].Result + "\nSoll: UNSUCCESSFUL")
	}

	t.Log("Success")
}

func TestInterpreterConditionNotPassed(t *testing.T) {
	_ = initMe()

	myModule := AuditModule{
		Condition: "false",
	}

	report, err := interpreter.InterpretAudit([]AuditModule{myModule}, 1, false)

	if err != nil {
		t.Errorf("Fehler aufgetreten: " + err.Error())
	}

	if report[0].Result != "NOTEXECUTED" {
		t.Errorf("Result: " + report[0].Result + "\nSoll: NOTPASSED")
	}

	t.Log("Success")
}

func TestInterpreterUndeclaredVariable(t *testing.T) {
	_ = initMe()

	myModule := AuditModule{
		ModuleName:       "FileContent",
		StepID:           "1",
		Variables:        VariableMap{},
		ModuleParameters: ParameterMap{"file": "%result%"},
	}

	_, err := interpreter.InterpretAudit([]AuditModule{myModule}, 1, false)

	if err == nil {
		t.Errorf("Nicht deklarierte Variablen wurden nicht abgefangen.")
	}

	t.Log("Success")
}
