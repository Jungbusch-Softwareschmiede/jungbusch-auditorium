// +build linux

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
	"regexp"
)

func (mh *MethodHandler) NftListRulesetInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:                 "NftListRuleset",
		ModuleDescription:          "Mit NftListRuleset lässt sich das Ruleset von nftables analysieren",
		ModuleAlias:                []string{"nftListRuleset", "nft_list_ruleset"},
		ModuleCompatibility:        []string{"linux"},
		RequiresElevatedPrivileges: true,
		InputParams: ParameterSyntaxMap{
			"awk": ParameterSyntax{
				ParamName:        "awk",
				ParamDescription: "Optionales Skript",
				IsOptional:       true,
			},
			"grep": ParameterSyntax{
				ParamName:        "grep",
				ParamDescription: "Suchbegriff, entspricht dem Pipen des Outputs in grep",
			},
		},
	}
}

// NftListRuleset listet das Ruleset von nftables auf
func (mh *MethodHandler) NftListRuleset(params ParameterMap) (r ModuleResult) {
	res, err := util.ExecCommand("nft list ruleset")
	if err != nil {
		r.Err = err
		return
	}

	r.ResultRaw = res
	r.Artifacts = append(r.Artifacts, Artifact{
		Name:  "nft list ruleset",
		Value: r.ResultRaw,
	})

	if params["awk"] != "" {
		r.Result = mh.AwkScript(ParameterMap{
			"input":     r.ResultRaw,
			"awkscript": params["awk"],
		}).Result

		r.Result = mh.Grep(ParameterMap{
			"input": r.Result,
			"grep":  params["grep"],
		}).Result
	} else {
		r.Result = mh.Grep(ParameterMap{
			"input": r.ResultRaw,
			"grep":  params["grep"],
		}).Result
	}

	return
}

func (mh *MethodHandler) NftListRulesetValidate(params ParameterMap) error {
	_, err := regexp.Compile(params["grep"])
	if err != nil {
		return errors.New("Modul: NftListRuleset - " + err.Error())
	}

	return nil
}
