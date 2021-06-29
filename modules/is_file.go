package modules

import (
	"fmt"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
)

func (mh *MethodHandler) IsFileInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "IsFile",
		ModuleDescription:   "IsFile prüft ob eine Datei existiert.",
		ModuleCompatibility: []string{"all"},
		InputParams: ParameterSyntaxMap{
			"path": ParameterSyntax{
				ParamName:        "path",
				ParamDescription: "Pfad zur Datei",
			},
		},
	}
}

func (mh *MethodHandler) IsFile(params ParameterMap) (r ModuleResult) {
	path := params["path"]
	isFile := util.IsFile(path)
	if isFile {
		r.ResultRaw = fmt.Sprintf("Dateipfad: %v, Ist eine Datei: true", path)
		r.Result = "true"
	} else {
		r.ResultRaw = fmt.Sprintf("Dateipfad: %v, Ist eine Datei: false", path)
		r.Result = "false"
	}
	return
}

func (mh *MethodHandler) IsFileValidate(params ParameterMap) error {
	if _, err := util.GetAbsolutePath(params["path"]); err != nil {
		return errors.New("Der angegebene Pfad ist ungültig.")
	}
	return nil
}
