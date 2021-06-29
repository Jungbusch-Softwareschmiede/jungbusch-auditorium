// +build linux

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
)

func (mh *MethodHandler) SysctlInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "Sysctl",
		ModuleDescription:   "Sysctl wird dazu verwendet, Kernelparameter zur Laufzeit zu ändern. Die verfügbaren Parameter sind unter /proc/sys/ aufgelistet. Für die sysctl-Unterstützung in Linux ist Procfs notwendig. Sie können sysctl sowohl zum Lesen als auch zum Schreiben von Sysctl-Daten verwenden.",
		ModuleAlias:         []string{"sysctl"},
		ModuleCompatibility: []string{"linux"},
		InputParams: ParameterSyntaxMap{
			"kernelparameter": ParameterSyntax{
				ParamName:        "kernelparameter",
				ParamAlias:       []string{"kernelparam"},
				ParamDescription: "Bezeichnet den Namen des Schlüssels, aus dem gelesen werden soll.",
			},
		},
	}
}

//sysctl  wird  dazu  verwendet,  Kernelparameter  zur  Laufzeit  zu ändern. Die verfügbaren
//Parameter sind unter /proc/sys/ aufgelistet. Für die  sysctl-Unterstützung  in  Linux  ist
//Procfs  notwendig.  Sie  können  sysctl  sowohl  zum  Lesen  als  auch  zum  Schreiben von
//Sysctl-Daten verwenden.
func (mh *MethodHandler) Sysctl(params ParameterMap) (r ModuleResult) {
	res_param1, err := util.ExecCommand("sysctl " + params["kernelparameter"])
	if err != nil {
		r.Err = err
		return
	}
	r.ResultRaw = res_param1
	r.Artifacts = append(r.Artifacts, Artifact{Name: "sysctl " + params["kernelparameter"], Value: r.ResultRaw})

	r.Result = r.ResultRaw
	return
}

func (mh *MethodHandler) SysctlValidate(params ParameterMap) error {

	return nil
}
