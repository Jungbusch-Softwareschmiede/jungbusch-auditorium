// +build windows

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
	s "strings"
)

func (mh *MethodHandler) GetUserSIDInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "GetUserSID",
		ModuleDescription:   "GetUserSID gibt die SID des angegebenen Nutzer aus.",
		ModuleAlias:         []string{"getusersid"},
		ModuleCompatibility: []string{"windows"},
		InputParams: ParameterSyntaxMap{
			"userName": ParameterSyntax{
				ParamName:        "userName",
				ParamAlias:       []string{},
				ParamDescription: "Windows Nutzername dessen SID bestimmt werden soll.",
			},
		},
	}
}

func (mh *MethodHandler) GetUserSID(params ParameterMap) (r ModuleResult) {
	// Ausgabe der SID
	sid, err := util.ExecCommand("cmd.exe /c \"wmic useraccount where name='" + params["userName"] + "' get sid\"")

	if err != nil {
		r.Err = err
		return
	}

	r.Artifacts = append(r.Artifacts, Artifact{
		Name:  "cmd.exe /c \"wmic useraccount where name='" + params["userName"] + "' get sid\"",
		Value: sid,
	})

	r.ResultRaw = sid

	// Prüfen ob ein Fehler aufgetreten ist
	if !s.Contains(sid, "SID") {
		r.Err = errors.New("Der Nutzername ist nicht gültig.")
		return
	}

	r.Result = s.TrimSpace(sid[s.LastIndex(sid, "S-"):])
	return
}
