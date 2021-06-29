package modules

import (
	"fmt"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"os"
)

func (mh *MethodHandler) PermissionsInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "Permissions",
		ModuleDescription:   "Permissions gibt die Berechtigungen der/des angegebenen Datei/Ordners zurück.",
		ModuleAlias:         []string{"perms", "permissions"},
		ModuleCompatibility: []string{"all"},
		InputParams: ParameterSyntaxMap{
			"path": ParameterSyntax{
				ParamName:        "path",
				ParamAlias:       []string{"pfad"},
				ParamDescription: "Pfad der/des zu überprüfenden Datei/Ordners",
			},
		},
	}
}

// Permissions gibt Berechtigungen der/des übergebenden Datei/Ordners in numerischer Form als String zurück.
func (mh *MethodHandler) Permissions(params ParameterMap) (r ModuleResult) {
	info, err := os.Stat(params["path"])
	if err != nil {
		r.Err = err
		return
	}
	r.Result = fmt.Sprintf("%#o", info.Mode().Perm())
	return
}

func (mh *MethodHandler) PermissionsValidate(params ParameterMap) error {
	return nil
}
