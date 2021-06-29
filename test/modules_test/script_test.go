package modules_test

import (
	"fmt"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/modules"
	"os"
	"strings"
	"testing"
)

func TestScript(t *testing.T) {
	path := strings.ReplaceAll(os.Getenv("gopath"), "\\", "\\\\")

	mh := modules.MethodHandler{}
	params := models.ParameterMap{
		"script": `function runModule() {
params.file = "` + path + `\\src\\github.com\\Jungbusch-Softwareschmiede\\jungbusch-auditorium\\test\\testdata\\filme.txt";
params.grep = "Rosso";

res = FileContent(params);
if(test == 'TEST_WERT') {
test = 'GEÄNDERT';
return res;
}
if(res.result === "Porco Rosso\n"){
	params.command = "hostname"
	params.grep = ""
	res = ExecuteCommand(params)
}
return res;
}
runModule();
`,
	}
	variables := models.VariableMap{"%test%": models.Variable{
		Name:  "%test%",
		Value: "TEST_WERT",
		IsEnv: false,
	}}
	result := mh.Script(params, &variables)

	if result.Err != nil {
		t.Errorf("Fehlgeschlagen: %v", result.Err)
	}

	fmt.Println(result.Result)
	fmt.Println(variables)
}

func TestCustomResult(t *testing.T) {
	mh := modules.MethodHandler{}
	params := models.ParameterMap{
		"script": `function runModule() {
result = newResult('Hallo', 'Hallo dies ist ein Test', null);

return result;
}
runModule();`,
	}

	result := mh.Script(params, &models.VariableMap{})

	if result.Result != "Hallo" {
		t.Errorf("Result nicht richtig: " + result.Result)
	} else if result.ResultRaw != "Hallo dies ist ein Test" {
		t.Errorf("ResultRaw nicht richtig: " + result.ResultRaw)
	}
}
