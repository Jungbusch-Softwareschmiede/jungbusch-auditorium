// +build linux

package modules

import (
	"errors"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"regexp"
)

func (mh *MethodHandler) SshdInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "Sshd",
		ModuleDescription:   "Sshd liest die Informationen aus der sshd Config-Datei aus ",
		ModuleAlias:         []string{"sshd"},
		ModuleCompatibility: []string{"linux"},
		InputParams: ParameterSyntaxMap{
			"grep": ParameterSyntax{
				ParamName:        "grep",
				IsOptional:       true,
				ParamDescription: "Optionaler Suchbegriff, entspricht dem Pipen des Outputs in grep",
			},
		},
	}
}

// Sshd liest die Informationen aus der sshd Config-Datei aus
func (mh *MethodHandler) Sshd(params ParameterMap) (r ModuleResult) {
	res, err := util.ExecCommand("sshd -T")
	if err != nil {
		r.Err = err
		return
	}

	r.ResultRaw = res
	r.Result = r.ResultRaw
	r.Artifacts = append(r.Artifacts, Artifact{
		Name:  "sshd -T",
		Value: r.ResultRaw,
	})
	if params["grep"] != "" {
		r.Result = mh.Grep(ParameterMap{
			"input": r.ResultRaw,
			"grep":  params["grep"],
		}).Result
	}

	return
}

func (mh *MethodHandler) SshdValidate(params ParameterMap) error {
	_, err := regexp.Compile(params["grep"])
	if err != nil {
		return errors.New("Modul: Sshd - " + err.Error())
	}

	return nil
}
