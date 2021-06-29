// +build linux

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
	"regexp"
)

func (mh *MethodHandler) AuditctlInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:                 "Auditctl",
		ModuleDescription:          "Mit Auditctl können die Kernel-Audit-Regeln ausgelesen werden",
		ModuleAlias:                []string{"auditctl"},
		ModuleCompatibility:        []string{"linux"},
		RequiresElevatedPrivileges: true,
		InputParams: ParameterSyntaxMap{
			"grep": ParameterSyntax{
				ParamName:        "grep",
				ParamAlias:       []string{"name"},
				ParamDescription: "Suchbegriff, entspricht Pipen des Outputs in grep",
			},
		},
	}
}

// Auditctl liefert die zurzeit geladenen Auditregeln
func (mh *MethodHandler) Auditctl(params ParameterMap) (r ModuleResult) {
	res, err := util.ExecCommand("auditctl -l")
	if err != nil {
		r.Err = err
		return
	}

	r.ResultRaw = res

	r.Result = mh.Grep(ParameterMap{
		"input": r.ResultRaw,
		"grep":  params["grep"],
	}).Result

	r.Artifacts = append(r.Artifacts, Artifact{
		Name:  "auditctl -l | grep " + params["grep"],
		Value: r.Result,
	})

	return
}

func (mh *MethodHandler) AuditctlValidate(params ParameterMap) error {

	_, err := regexp.Compile(params["grep"])
	if err != nil {
		return errors.New("Modul: Auditctl - " + err.Error())
	}

	return nil
}
