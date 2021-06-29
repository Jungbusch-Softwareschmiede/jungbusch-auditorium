// +build linux

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
)

func (mh *MethodHandler) StatInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:                 "Stat",
		ModuleDescription:          "Mit dem Befehl stat lassen sich Zugriffs- und Änderungs-Zeitstempel von Dateien und Ordnern anzeigen. Weiterhin werden Informationen zu Rechten, zu Besitzer und Gruppe und zum Dateityp ausgegeben.",
		ModuleAlias:                []string{"stat"},
		ModuleCompatibility:        []string{"linux"},
		RequiresElevatedPrivileges: true,
		InputParams: ParameterSyntaxMap{
			"file": ParameterSyntax{
				ParamName:        "file",
				ParamAlias:       []string{"datei"},
				ParamDescription: "Pfad zur Datei",
			},
		},
	}
}

// Mit dem Befehl Stat lassen sich Zugriffs- und Änderungs-Zeitstempel von Dateien und Ordnern anzeigen.
// Weiterhin werden Informationen zu Rechten, zu Besitzer und Gruppe und zum Dateityp ausgegeben.
func (mh *MethodHandler) Stat(params ParameterMap) (r ModuleResult) {
	res, err := util.ExecCommand("stat " + params["file"])
	if err != nil {
		r.Err = err
		return
	}

	r.ResultRaw = res
	r.Result = mh.Grep(ParameterMap{
		"input": r.ResultRaw,
		"grep":  "Uid:",
	}).Result

	r.Artifacts = append(r.Artifacts, Artifact{Name: "stat " + params["file"], Value: r.Result})

	return
}

func (mh *MethodHandler) StatValidate(params ParameterMap) error {
	if !util.IsFile(params["file"]) {
		return errors.New("Datei ist nicht vorhanden")
	} else {
		return nil
	}
}
