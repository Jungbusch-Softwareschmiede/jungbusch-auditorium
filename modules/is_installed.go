package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"os/exec"
)

func (mh *MethodHandler) IsInstalledInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "IsInstalled",
		ModuleDescription:   "IsInstalled überprüft, ob das angegebene Package installiert ist.",
		ModuleAlias:         []string{"isInstalled", "is_installed", "installed"},
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

// IsInstalled liefert Informationen darüber, ob das angegeben Package installiert ist.
func (mh *MethodHandler) IsInstalled(params ParameterMap) (r ModuleResult) {
	res, err := exec.LookPath(params["package"])
	if err != nil {
		r.ResultRaw = res
		r.Result = "false"
	} else {
		// Programm ist installiert
		r.ResultRaw = params["package"] + " ist installiert."
		r.Result = "true"
	}

	return
}

func (mh *MethodHandler) IsInstalledValidate(params ParameterMap) error {

	return nil
}
