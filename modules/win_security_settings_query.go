// +build windows

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"path/filepath"
	s "strings"
)

func (mh *MethodHandler) SecuritySettingsQueryInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:                 "SecuritySettingsQuery",
		ModuleDescription:          "SecuritySettingsQuery ermöglicht verschiedene Bereiche der Security Settings auszulesen, die mit RegistryQuery nicht auslesbar sind.",
		ModuleCompatibility:        []string{"windows"},
		ModuleAlias:                []string{"securitysettingsquery"},
		RequiresElevatedPrivileges: true,
		InputParams: ParameterSyntaxMap{
			"valueName": ParameterSyntax{
				ParamName:        "valueName",
				ParamDescription: "ValueName der Security Setting. Bitte im Benutzerhandbuch nachschlagen.",
			},
			"path": ParameterSyntax{
				ParamName:        "path",
				ParamDescription: "Pfad zur Dump-Datei. Muss nur angegeben werden, wenn im Dump-Modul ein Pfad spezifiziert wurde.",
				IsOptional:       true,
			},
		},
	}
}

func (mh *MethodHandler) SecuritySettingsQuery(params ParameterMap) (r ModuleResult) {
	if params["path"] == "" {
		params["path"] = static.TempPath + static.PATH_SEPERATOR + "secedit.cfg"
	} else {
		// Wenn nur der Pfad zu einem Ordner angegeben ist wird der Dateiname ergänzt
		if util.IsDir(params["path"]) {
			params["path"] = filepath.Clean(params["path"] + static.PATH_SEPERATOR + "secedit.cfg")
		}
	}

	lines, err := util.ReadUTF16File(params["path"])

	if err != nil {
		r.Err = err
		return
	}

	// Im Stringslice nach Value suchen z.B.: MinimumPasswordAge
	for _, line := range lines {
		if s.Contains(line, params["valueName"]) {
			line = s.ReplaceAll(line, " ", "")
			line = s.ReplaceAll(line, "\"", "")

			r.ResultRaw = s.Split(line, "=")[1]
			r.Result = r.ResultRaw
		}
	}

	if r.Result == "" {
		r.Result = "keyNonExistent"
	}
	return
}
