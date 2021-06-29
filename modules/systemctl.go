// +build linux

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
)

func (mh *MethodHandler) SystemctlInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "Systemctl",
		ModuleDescription:   "Systemctl überprüft, ob die angegebene Unitdatei aktiviert ist.",
		ModuleAlias:         []string{"systemctl"},
		ModuleCompatibility: []string{"linux"},
		InputParams: ParameterSyntaxMap{
			"unitdatei": ParameterSyntax{
				ParamName:        "unitdatei",
				ParamAlias:       []string{"unitname"},
				ParamDescription: "Name der Unitdatei",
			},
		},
	}
}

// Systemctl überprüft, ob die angegebene Unitdatei aktiviert ist.
func (mh *MethodHandler) Systemctl(params ParameterMap) (r ModuleResult) {
	res, err := util.ExecCommand("systemctl is-enabled " + params["unitdatei"])
	if err != nil {
		r.Err = err
		return
	}

	r.ResultRaw = res
	r.Artifacts = append(r.Artifacts, Artifact{Name: "systemctl is-enabled " + params["unitdatei"], Value: r.ResultRaw})
	r.Result = r.ResultRaw
	return
}

func (mh *MethodHandler) SystemctlValidate(params ParameterMap) error {

	return nil
}
