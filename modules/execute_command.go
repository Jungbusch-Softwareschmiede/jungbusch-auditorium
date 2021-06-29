package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"strings"
)

func (mh *MethodHandler) ExecuteCommandInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "ExecuteCommand",
		ModuleDescription:   "ExecuteCommand führt den übergebenen Befehl aus und überprüft optional, ob der angegebene Suchbegriff im Ergebnis vorhanden ist.",
		ModuleAlias:         []string{"execute_command", "executeCommand"},
		ModuleCompatibility: []string{"all"},
		InputParams: ParameterSyntaxMap{
			"command": ParameterSyntax{
				ParamName:        "command",
				ParamAlias:       []string{"cmd"},
				ParamDescription: "Der auszuführende Befehl",
			},
			"grep": ParameterSyntax{
				ParamName:        "grep",
				IsOptional:       true,
				ParamDescription: "Optionaler Suchbegriff, entspricht dem Pipen des Outputs in grep",
			},
		},
	}
}

// ExecuteCommand führt den übergebenen Befehl aus und speichert das Ergebnis in einem String.
func (mh *MethodHandler) ExecuteCommand(params ParameterMap) (r ModuleResult) {
	out, err := util.ExecCommand(params["command"])
	if err != nil {
		r.Err = err
		return
	}

	out = strings.ReplaceAll(out, "\r", "")
	r.ResultRaw = out
	r.Result = r.ResultRaw

	if params["grep"] != "" {
		r.Result = mh.Grep(ParameterMap{
			"input": r.ResultRaw,
			"grep":  params["grep"],
		}).Result
	}

	r.Artifacts = append(r.Artifacts, Artifact{Name: params["command"], Value: r.ResultRaw})

	return
}

func (mh *MethodHandler) ExecuteCommandValidate(params ParameterMap) error {

	// if params["command"] == "" {
	//	return errors.New("der Command-Parameter darf nicht leer sein")
	// }
	//
	// splitCommand := strings.Fields(params["command"])
	//
	// _, err := exec.Command(splitCommand[0], splitCommand[1:]...).Output()
	// if err != nil {
	//	return errors.New("Modul: ExecuteCommand - " + err.Error())
	// }
	return nil
}
