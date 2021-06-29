package outputgenerator_test

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/interpreter"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/modulecontroller"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/outputgenerator"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"os"
	s "strings"
	"testing"
	"time"
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

func TestGenerateReport(t *testing.T) {
	gopath := os.Getenv("gopath") + `\src\github.com\Jungbusch-Softwareschmiede\jungbusch-auditorium\`

	if err := logger.InitializeLogger(&ConfigStruct{}, []LogMsg{}); err != nil {
		t.Errorf(err.Error())
	}
	e := AuditModule{
		ModuleName:       "ExecuteCommand",
		StepID:           "Hostname",
		Passed:           "true",
		Variables:        initVarMap(),
		ModuleParameters: ParameterMap{"command": "hostname"},
	}
	u := AuditModule{
		ModuleName:       "FileContent",
		StepID:           "Film-Inhalt3",
		Passed:           "%result% === 'Porco Rosso'",
		Variables:        initVarMap(),
		ModuleParameters: ParameterMap{"file": "test"},
	}
	n := AuditModule{
		ModuleName:       "FileContent",
		StepID:           "Film-Inhalt2",
		Passed:           "false",
		Variables:        initVarMap(),
		ModuleParameters: ParameterMap{"file": gopath + "test\\testdata\\filme.txt", "grep": "Porco"},
		NestedModules:    []AuditModule{u},
	}
	p := AuditModule{
		ModuleName:       "FileContent",
		StepID:           "Film-Inhalt1",
		Passed:           "%result% === 'Porco Rosso'",
		Variables:        initVarMap(),
		ModuleParameters: ParameterMap{"file": gopath + "test\\testdata\\filme.txt", "grep": "Porco"},
		NestedModules:    []AuditModule{n},
	}
	a := AuditModule{
		ModuleName: "Script",
		StepID:     "TestScript",
		Passed:     "true",
		Variables:  initVarMap(),
		ModuleParameters: ParameterMap{
			"script": `function runModule() {
params.file = "` + s.ReplaceAll(gopath, "\\", "\\\\") + `test\\testdata\\filme.txt";
params.grep = "Rosso";

r = FileContent(params);
	if (r.result == 'Porco Rosso') {
	params.command = "ipconfig";
	params.grep = ""
	r = ExecuteCommand(params);
}
return r;
}
runModule();
`,
		},
	}

	static.OperatingSystem, _ = modulecontroller.GetOS()
	ms := []AuditModule{p, e, a}

	r, err := interpreter.InterpretAudit(ms, 5, true)

	if err != nil {
		t.Error(err)
	} else {

		if err = outputgenerator.GenerateOutput(r, gopath+"ProgrammOutput/testoutput", 5, false, time.Now(), time.Since(time.Now())); err != nil {
			t.Errorf(err.Error())
		}
	}

	_ = os.RemoveAll(gopath + "ProgrammOutput/testoutput")
	t.Log("Success")
}
