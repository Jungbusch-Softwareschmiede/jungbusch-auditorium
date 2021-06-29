// +build linux

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
)

func (mh *MethodHandler) AwkScriptInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "AwkScript",
		ModuleDescription:   "Awk ist eine Skriptsprache zum Editieren und Analysieren von Texten. AwkScript führt ein Skript auf die Input-Datei aus",
		ModuleAlias:         []string{"awk", "awkScript", "awk_script"},
		ModuleCompatibility: []string{"linux"},
		InputParams: ParameterSyntaxMap{
			"input": ParameterSyntax{
				ParamName:        "input",
				ParamDescription: "String, auf den das Skript angewendet wird",
			},
			"awkscript": ParameterSyntax{
				ParamName:        "awkscript",
				ParamAlias:       []string{"script"},
				ParamDescription: "Awk-Script, welches ausgeführt werden soll",
			},
			"separator": ParameterSyntax{
				ParamName:        "separator",
				IsOptional:       true,
				ParamDescription: "Legt den Field-Separator fest",
			},
		},
	}
}

// Awk ist eine Skriptsprache zum Editieren und Analysieren von Texten. AwkScript führt ein Awk-Script, auf die Input-Datei aus
func (mh *MethodHandler) AwkScript(params ParameterMap) (r ModuleResult) {
	var separator string
	if params["separator"] != "" {
		separator = "-F" + params["separator"]
	}
	res, err := util.ExecCommand("awk " + separator + " " + params["awkscript"])
	if err != nil {
		r.Err = err
		return
	}

	r.ResultRaw = res
	r.Artifacts = append(r.Artifacts, Artifact{Name: "awk " + params["awkscript"] + params["input"], Value: r.ResultRaw})
	r.Result = r.ResultRaw
	return
}

func (mh *MethodHandler) AwkScriptValidate(params ParameterMap) error {
	// ?
	return nil
}
