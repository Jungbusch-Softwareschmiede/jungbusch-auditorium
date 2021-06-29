package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"os/exec"
)

func (mh *MethodHandler) IsNotInstalledInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "IsNotInstalled",
		ModuleDescription:   "IsNotInstalled überprüft, ob das angegebene Package nicht installiert ist.",
		ModuleAlias:         []string{"isNotInstalled", "is_not_installed", "notInstalled", "not_installed"},
		ModuleCompatibility: []string{"linux", "darwin"},
		InputParams: ParameterSyntaxMap{
			"package": ParameterSyntax{
				ParamName:        "package",
				ParamAlias:       []string{"name", "packagename", "pkg"},
				ParamDescription: "Name des Packages",
			},
		},
	}
}

// IsInstalled liefert Informationen darüber, ob das angegebene Package nicht installiert ist.
func (mh *MethodHandler) IsNotInstalled(params ParameterMap) (r ModuleResult) {
	res, err := exec.LookPath(params["package"])
	if err != nil {
		r.ResultRaw = params["package"] + " scheint nicht installiert zu sein."
		r.Result = "true"
	} else {
		// Programm ist installiert
		r.ResultRaw = res
		r.Result = "false"
	}
	return
}

func (mh *MethodHandler) IsNotInstalledValidate(params ParameterMap) error {

	return nil
}
