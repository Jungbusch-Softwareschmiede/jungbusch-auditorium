// +build linux

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
	"regexp"
)

func (mh *MethodHandler) AuthselectInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "Authselect",
		ModuleDescription:   "Mit Authselect lässt sich die Konfiguration des authselect-Profils überprüfen.",
		ModuleAlias:         []string{"authselect"},
		ModuleCompatibility: []string{"rhel"},
		InputParams: ParameterSyntaxMap{
			"grep": ParameterSyntax{
				ParamName:        "grep",
				IsOptional:       true,
				ParamDescription: "Optionaler Suchbegriff, entspricht dem Pipen des Outputs in grep",
			},
		},
	}
}

// ExecuteCommand führt den übergebenen Befehl aus und speichert das Ergebnis in einem String.
func (mh *MethodHandler) Authselect(params ParameterMap) (r ModuleResult) {
	res, err := util.ExecCommand("authselect current")
	if err != nil {
		r.Err = err
		return
	}

	r.ResultRaw = res
	r.Result = r.ResultRaw
	r.Artifacts = append(r.Artifacts, Artifact{
		Name:  "authselect current",
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

func (mh *MethodHandler) AuthselectValidate(params ParameterMap) error {
	_, err := regexp.Compile(params["grep"])
	if err != nil {
		return errors.New("Modul: Authselect - " + err.Error())
		// besser: return errors.Wrap(err, "Modul: Authselect")
	}

	return nil
}
