// +build linux

package modules

import (
	"errors"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"regexp"
)

func (mh *MethodHandler) CheckPartitionInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "CheckPartition",
		ModuleDescription:   "CheckPartition gibt alle eingehängten Datenträger aus. Ist der grep-Parameter gesetzt, werden auch Zeilen zurückgegeben, in denen das Pattern gefunden wurde.",
		ModuleAlias:         []string{"mount"},
		ModuleCompatibility: []string{"linux"},
		InputParams: ParameterSyntaxMap{
			"grep": ParameterSyntax{
				ParamName:        "grep",
				IsOptional:       true,
				ParamDescription: "Optionaler Suchbegriff, entspricht dem Pipen des Outputs in grep",
			},
			"vgrep": ParameterSyntax{
				ParamName:        "vgrep",
				IsOptional:       true,
				ParamDescription: "Optionaler Suchbegriff, entspricht dem Pipen des Outputs in grep -v",
			},
		},
	}
}

// CheckPartition gibt alle eingehängten Datenträger aus.
// Ist der grep-Parameter gesetzt, werden auch Zeilen zurückgegeben, in denen
// das Pattern gefunden wurde.
func (mh *MethodHandler) CheckPartition(params ParameterMap) (r ModuleResult) {
	res, err := util.ExecCommand("mount")
	if err != nil {
		r.Err = err
		return
	}

	r.ResultRaw = res
	r.Result = r.ResultRaw
	r.Artifacts = append(r.Artifacts, Artifact{
		Name:  "mount",
		Value: r.ResultRaw,
	})
	if params["grep"] != "" {
		r.Result = mh.Grep(ParameterMap{
			"input": r.ResultRaw,
			"grep":  params["grep"],
		}).Result
	}

	if params["vgrep"] != "" {
		r.Result = mh.Grep(ParameterMap{
			"input": r.Result,
			"grep":  "-v " + params["grep"],
		}).Result
	}

	return
}

func (mh *MethodHandler) CheckPartitionValidate(params ParameterMap) error {
	_, err := regexp.Compile(params["grep"])
	if err != nil {
		return errors.New("Modul: CheckPartition - " + err.Error())
	}

	_, err = regexp.Compile(params["vgrep"])
	if err != nil {
		return errors.New("Modul: CheckPartition - " + err.Error())
	}

	return nil
}
