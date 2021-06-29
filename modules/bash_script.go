// +build linux

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"os/exec"
)

func (mh *MethodHandler) BashScriptInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "BashScript",
		ModuleDescription:   "BashScript führt das verlinkte Bash-Script aus. Als Ergebnis wird der Output des Skripts erhalten.",
		ModuleAlias:         []string{"bash_script", "bashScript"},
		ModuleCompatibility: []string{"linux"},
		InputParams: ParameterSyntaxMap{
			"script": ParameterSyntax{
				ParamName:        "script",
				ParamAlias:       []string{},
				ParamDescription: "Das auszuführende Bash-Script",
			},
		},
	}
}

// ExecuteCommand führt den übergebenen Befehl aus und speichert das Ergebnis in einem String.
func (mh *MethodHandler) BashScript(params ParameterMap) (r ModuleResult) {
	script, err := util.GetAbsolutePath(params["script"])
	if err != nil {
		r.Err = err
		return
	}
	out, err := exec.Command("bash", script).CombinedOutput()
	if err != nil {
		r.Err = err
		return
	}

	r.ResultRaw = string(out)
	r.Result = r.ResultRaw

	r.Artifacts = append(r.Artifacts, Artifact{Name: params["script"], Value: r.ResultRaw})

	return
}
