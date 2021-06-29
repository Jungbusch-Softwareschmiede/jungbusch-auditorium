// +build windows

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
	s "strings"
)

func (mh *MethodHandler) GetWinEnvInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "GetWinEnv",
		ModuleDescription:   "GetWinEnv gibt den Wert der angegebenen Windows Umgebungsvariable aus.",
		ModuleAlias:         []string{},
		ModuleCompatibility: []string{"windows"},
		InputParams: ParameterSyntaxMap{
			"envVar": ParameterSyntax{
				ParamName:        "envVar",
				ParamAlias:       []string{},
				ParamDescription: "Name der Windows-Umgebungsvariable.",
			},
		},
	}
}

func (mh *MethodHandler) GetWinEnv(params ParameterMap) (r ModuleResult) {
	envVar, err := util.ExecCommand("cmd.exe /c \"echo %" + params["envVar"] + "%\"")
	
	if err != nil {
		r.Err = err
		return
	}

	r.Artifacts = append(r.Artifacts, Artifact{
		Name:  "cmd.exe /c \"echo %" + params["envVar"] + "%\"",
		Value: envVar,
	})

	r.ResultRaw = envVar

	// Prüfen ob ein Fehler aufgetreten ist
	if s.TrimSpace(envVar) == params["envVar"] {
		r.Err = errors.New("Bei abfragen der Windows Umgebungsvariable '" + params["envVar"] + "'ist ein Fehler aufgetreten.")
		return
	}

	r.Result = s.TrimSpace(envVar)
	return
}
