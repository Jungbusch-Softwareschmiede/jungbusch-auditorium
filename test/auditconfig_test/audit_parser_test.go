package auditconfig_test

import (
	"fmt"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/parser"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/syntaxchecker"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/modulecontroller"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	s "strings"
	"testing"
)

//
//
//
// *=*=*=*=*=*=*=*=*=*=*=*=
//
// Positive Tests
//
// *=*=*=*=*=*=*=*=*=*=*=*=
//
//
//

func TestParseAuditConfigurationGlobalValid(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestGlobalVariable", func(t *testing.T) {
		m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/global.jba")
		checkMe(t, err)
		fmt.Println(err)
		if !m[0].IsGlobal {
			t.Errorf("Das Modul sollte global sein, ist es aber nicht!")
		}
	})
	t.Run("TestMulitpleGlobalVariables", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_global.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestNonGlobalUseGlobal", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/non_global_use_global.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationCaseSensitivity(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestVariableCaseSenitivity1", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/case_sensitivity1.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestVariableCaseSenitivity2", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/case_sensitivity2.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestVariableCaseSenitivityScript1", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/case_sensitivity_script1.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestVariableCaseSenitivityScript2", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/case_sensitivity_script2.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationAlias(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestNameAlias", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/name_alias.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestParamAlias", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/param_alias.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestIfAlias", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/if_alias.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationRequireselevatedprivilegesValid(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestRequireselevatedprivilegesTrue1", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/requireselevatedprivileges_true1.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestRequireselevatedprivilegesTrue2", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/requireselevatedprivileges_true2.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestRequireselevatedprivilegesTrue3", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/requireselevatedprivileges_true3.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestRequireselevatedprivilegesFalse1", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/requireselevatedprivileges_false1.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestRequireselevatedprivilegesFalse2", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/requireselevatedprivileges_false2.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestRequireselevatedprivilegesFalse3", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/requireselevatedprivileges_false3.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationPrint(t *testing.T) {
	loadedModules := initMe()
	_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/print.jba")
	checkMe(t, err)
	fmt.Println(err)
}

func TestParseAuditConfigurationDescription(t *testing.T) {
	loadedModules := initMe()
	_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/description.jba")
	checkMe(t, err)
	fmt.Println(err)
}

func TestParseAuditConfigurationBackticks(t *testing.T) {
	loadedModules := initMe()
	_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/backticks.jba")
	checkMe(t, err)
	fmt.Println(err)
}

func TestParseAuditConfigurationMultilineValid(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestMultiValid", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multi_valid.jba")
		checkMe(t, err)

	})
	t.Run("TestMultiValidOneLine", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multi_valid_one_line.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationBasic(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestBasicModule", func(t *testing.T) {
		m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/basic_module.jba")
		checkMe(t, err)
		fmt.Println(err)

		n := recursivelyTestNestedModules(m, t, 0)
		if n != 1 {
			t.Errorf("Falsche Anzahl Module: %v gefunden, sollte 1 sein", n)
		}
	})
	t.Run("TestBasicModuleNested", func(t *testing.T) {
		m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/basic_module_nested.jba")
		checkMe(t, err)
		fmt.Println(err)

		n := recursivelyTestNestedModules(m, t, 0)
		if n != 2 {
			t.Errorf("Falsche Anzahl Module: %v gefunden, sollte 2 sein", n)
		}
	})
	t.Run("TestMultipleBasicModules", func(t *testing.T) {
		m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_basic_modules.jba")
		checkMe(t, err)
		fmt.Println(err)

		n := recursivelyTestNestedModules(m, t, 0)
		if n != 3 {
			t.Errorf("Falsche Anzahl Module: %v gefunden, sollte 3 sein", n)
		}
	})
	t.Run("TestMultipleBasicModulesNested", func(t *testing.T) {
		m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_basic_modules_nested.jba")
		checkMe(t, err)
		fmt.Println(err)

		n := recursivelyTestNestedModules(m, t, 0)
		if n != 10 {
			t.Errorf("Falsche Anzahl Module: %v gefunden, sollte 10 sein", n)
		}
	})
}

func TestParseAuditConfigurationWhitespace(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestLeadingWhitespace", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/leading_whitespace.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestEmptyLines", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/empty_lines.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationKeywordSequence(t *testing.T) {
	loadedModules := initMe()

	_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/keyword_sequence.jba")
	checkMe(t, err)
	fmt.Println(err)
}

func TestParseAuditConfigurationMissingCondition(t *testing.T) {
	loadedModules := initMe()

	_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_condition.jba")
	checkMe(t, err)
	fmt.Println(err)
}

func TestParseAuditConfigurationSymbolInVariable(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestColonInValue", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/colon_in_value.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestBraceInValue", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/brace_in_value.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestRoundBracketInValue", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/round_bracket_in_value.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationSymbolInParam(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestEqualSignInParam", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/equal_sign_in_param.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestBraceInParam", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/brace_in_param.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestRoundBracketInParam", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/round_bracket_in_param.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationSlashesAndCommentsInParams(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestSlashesInParam", func(t *testing.T) {
		m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/slashes_in_param.jba")
		checkMe(t, err)
		fmt.Println(err)
		if m[0].ModuleParameters["command"] != "te//st" {
			t.Errorf("Parameter wurde falsch geparsed")
		}
	})
	t.Run("TestSlashesInParamWithComment", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/slashes_in_param_with_comment.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestSlashesInParamWithQuotationMarksInComment", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/slashes_in_param_with_quotation_mark_in_comment.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestCommentAfterParam", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/comment_after_param.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestCommentAfterParamWithQuotiationMarks", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/comment_after_param_with_quotation_mark.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestEmptyParamWithComment", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/empty_param_with_comment.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestCommentAfterBrace", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/comment_after_brace.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationSlashesAndCommentsInVariables(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestSlashesInValue", func(t *testing.T) {
		m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/slashes_in_value.jba")
		checkMe(t, err)
		fmt.Println(err)
		if m[0].Variables["%test%"].Value != "te//st" {
			t.Errorf("Variable wurde falsch geparsed")
		}
	})
	t.Run("TestSlashesInValueWithComment", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/slashes_in_value_with_comment.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestSlashesInValueWithQuotationMarksInComment", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/slashes_in_value_with_quotation_mark_in_comment.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestCommentAfterValue", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/comment_after_value.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
	t.Run("TestCommentAfterValueWithQuotiationMarks", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/comment_after_value_with_quotation_mark.jba")
		checkMe(t, err)
		fmt.Println(err)
	})
}

//
//
//
// *=*=*=*=*=*=*=*=*=*=*=*=
//
// Negative Tests
//
// *=*=*=*=*=*=*=*=*=*=*=*=
//
//
//

func TestParseAuditConfigurationGlobalInvalid(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestGlobalWrongOrder", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/global_wrong_order.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleGlobalWrongOrder1", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_global_wrong_order1.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleGlobalWrongOrder2", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_global_wrong_order2.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestNestedGlobal1", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/nested_global1.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestNestedGlobal2", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/nested_global2.jba")
		checkMeNegated(t, err)
	})
	t.Run("TestGlobalOverwrite", func(t *testing.T) {
		m, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/global_overwrite.jba")
		fmt.Println(m[1].IsGlobal)
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationAliasNeg(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestAdditionalAlias", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/additional_alias.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationRequireselevatedprivilegesInvalid(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestRequireselevatedprivilegesInvalid1", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/requireselevatedprivileges_invalid1.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestRequireselevatedprivilegesInvalid2", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/requireselevatedprivileges_invalid2.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationPrintInvalid(t *testing.T) {
	loadedModules := initMe()
	_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/print_invalid.jba")
	checkMeNegated(t, err)
	fmt.Println(err)
}

func TestParseAuditConfigurationDescriptionInvalid(t *testing.T) {
	loadedModules := initMe()
	_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/description_invalid.jba")
	checkMeNegated(t, err)
	fmt.Println(err)
}

func TestParseAuditConfigurationMultipleParams(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestMultipleParams1", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_params1.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleParams2", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_params2.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleParams3", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_params3.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleParams4", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_params4.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleParams5", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_params5.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleParamsTrue", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_param_true.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleParamsFalse", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_param_false.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationMultilineInvalid(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestMultiInvalidTextAfterEnd", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multi_invalid_text_after_end.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultiInvalidOneLineAfterEnd", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multi_invalid_one_line_text_after_end.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultiInvalidQuotesInsteadOfTicks", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multi_invalid_quotes_instead_of_ticks.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultiInvalidEmpty", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multi_invalid_empty.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultiInvalidEOF", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multi_invalid_eof.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationEmptyFile(t *testing.T) {
	loadedModules := initMe()

	_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/empty_file.jba")
	checkMeNegated(t, err)
	fmt.Println(err)
}

func TestParseAuditConfigurationEmptyModule(t *testing.T) {
	loadedModules := initMe()

	_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/empty_module.jba")
	checkMeNegated(t, err)
	fmt.Println(err)
}

func TestParseAuditConfigurationBrace(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestOpeningBraceOnly", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/opening_brace_only.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingBrace", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_brace.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationWhitespaceInNames(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestWhitespaceInParamName", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/whitespace_in_param_name.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestWhitespaceInModuleName", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/whitespace_in_module_name.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationComma(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestMissingComma", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_comma.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestAdditionalComma", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/additional_comma.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationRedeclaredVariables(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestRedeclaredVariablesSameModule1", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/redeclared_variable_same_module1.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestRedeclaredVariablesSameModule2", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/redeclared_variable_same_module2.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationPercentSign(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestMissingBeginningPercentSign", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_beginning_percent_sign.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingEndingPercentSign", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_ending_percent_sign.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingPercentSigns", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_percent_signs.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationVariableDeclaration(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestMissingVariableName", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_variable_name.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingVariable", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_variable.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingEqualSign", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_equal_sign.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingVariableValue", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_variable_value.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingVariableAndEqualSign", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_variable_equal_sign.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestEqualSignOnly", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/equal_sign_only.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestVariableNameOnly", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/variable_name_only.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestInvalidVariableDeclaration1", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/invalid_variable_declaration1.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestInvalidVariableDeclaration2", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/invalid_variable_declaration2.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestInvalidVariableDeclaration3", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/invalid_variable_declaration3.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationParameterDeclaration(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestMissingParameter", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_parameter.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingParameterValue", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_parameter_value.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestParameterValueOnly", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/parameter_value_only.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestParameterNameOnly", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/parameter_name_only.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestInvalidParameterValue", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/invalid_parameter_value.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingColon", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_colon.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestColonOnly", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/colon_only.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingModuleName", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_module_name.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleModuleDeclaration", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_module_declaration.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationQuotationMarks(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestOneQuotationMarkOnly", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/one_quotation_mark_only.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingOpeningQuotationMark", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_opening_quotation_mark.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingClosingQuotationMark", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_closing_quotation_mark.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingQuotationMarks", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_quotation_marks.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationStepID(t *testing.T) {
	loadedModules := initMe()

	t.Run("TestMissingStepID", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_stepid.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMissingStepIDNested", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/missing_stepid_nested.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleStepID", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_stepid.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleStepIDNested", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_stepid_nested.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
	t.Run("TestMultipleStepIDDeclaration", func(t *testing.T) {
		_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/multiple_stepid_declaration.jba")
		checkMeNegated(t, err)
		fmt.Println(err)
	})
}

func TestParseAuditConfigurationOverwrittenEnvVariable(t *testing.T) {
	loadedModules := initMe()

	_, _, _, err := executeMe(loadedModules, gopath+"/test/testdata/parser_testdata/overwrite_env_variable.jba")
	checkMeNegated(t, err)
	fmt.Println(err)
}

//
//
//
// *=*=*=*=*=*=*=*=*=*=*=*=
//
// Utility
//
// *=*=*=*=*=*=*=*=*=*=*=*=
//
//
//

func executeMe(loadedModules []ModuleSyntax, file string) ([]AuditModule, bool, int, error) {
	lines, err := util.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = syntaxchecker.Syntax(lines)
	fmt.Println(err)
	if err != nil {
		return nil, false, 0, err
	}

	return parser.Parse(lines, loadedModules)
}

func initMe() []ModuleSyntax {
	cs := ConfigStruct{OutputPath: gopath + "./ProgrammOutput"}
	err := logger.InitializeLogger(&cs, []LogMsg{})
	if err != nil {
		panic(err)
	}
	static.OperatingSystem = "windows10"
	loadedModules, err := modulecontroller.Initialize(false)
	if err != nil {
		panic(err)
	}
	return loadedModules
}

func checkMe(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Fehler beim Parsen: %v", err)
	}
}

func checkMeNegated(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Invalider Syntax wurde angenommen.")
	} else {
		if s.HasPrefix(err.Error(), "Die Audit-Datei ") {
			t.Errorf("Invalide Datei.")
		}
	}
}

func recursivelyTestNestedModules(auditModules []AuditModule, t *testing.T, n int) int {
	for i := range auditModules {
		if auditModules[i].ModuleName != "ExecuteCommand" {
			t.Errorf("Falscher Modul-Name: %v, sollte ExecuteCommand sein", auditModules[i].ModuleName)
		}
		for k, v := range auditModules[i].ModuleParameters {
			if v != "test" {
				t.Errorf("%v-Parameter falsch gesetzt: \"%v\" erhalten, sollte \"test\" sein", k, v)
			}
		}
		if len(auditModules[i].NestedModules) != 0 {
			n = recursivelyTestNestedModules(auditModules[i].NestedModules, t, n)
		}
		n++
	}
	return n
}
