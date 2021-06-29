// +build windows

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
)

func (mh *MethodHandler) DumpSecuritySettingsInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:                 "DumpSecuritySettings",
		ModuleDescription:          "DumpSecuritySettings ermöglicht das Dumpen der aktuellen Security-Settings in eine Datei, welche mit dem Modul SecuritySettingsQuery ausgelesen werden kann.",
		ModuleAlias:                []string{"dumpsecuritysettings"},
		ModuleCompatibility:        []string{"windows"},
		RequiresElevatedPrivileges: true,
		InputParams: ParameterSyntaxMap{
			"path": ParameterSyntax{
				ParamName:        "path",
				IsOptional:       true,
				ParamDescription: "Hier kann ein Pfad angegeben werden, an welchem die zu dumpende Datei abgelegt werden soll. Wird kein Pfad angegeben, wird sie in einem Temporären Ordner abgelegt. Unabhängig davon, wird sie so oder so in die Modul-Artefakte und somit den Output aufgenommen. Wird hier ein Pfad angegeben, muss dieser auch im Query-Modul explizit angegeben werden.",
			},
		},
	}
}

func (mh *MethodHandler) DumpSecuritySettings(params ParameterMap) (r ModuleResult) {
	var path string
	var err error

	// Pfad bestimmen
	if params["path"] != "" {
		path, err = util.GetAbsolutePath(params["path"] + "secedit.cfg")
	} else {
		if !util.IsDir(params["path"]) {
			r.Err = errors.New("Der angegebene Pfad ist ungültig.")
			return
		}
		path = static.TempPath + static.PATH_SEPERATOR + "secedit.cfg"
	}

	if err != nil {
		r.Err = err
	}

	// secedit.cfg Datei generieren
	exportErr, err := util.ExecCommand("secedit /export /cfg \"" + path + "\" ")
	if err != nil {
		r.Err = err
		return
	}

	r.Artifacts = append(r.Artifacts, Artifact{
		Name:   "file",
		Value:  path,
		IsFile: true,
	})

	if !util.IsFile(path) {
		r.Err = errors.New("Beim Exportieren der Security-Policy-Datei ist möglicherweise ein Fehler aufgetreten:\n" + exportErr)
		return
	}

	r.Result = "Die Security-Policy-Datei wurde exportiert."
	return
}
