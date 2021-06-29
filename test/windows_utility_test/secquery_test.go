package windows_utility

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/modules"
	"os"
	"testing"
)

var (
	gopath  = os.Getenv("gopath") + `\src\github.com\Jungbusch-Softwareschmiede\jungbusch-auditorium\`
	handler = modules.MethodHandler{}
)

func TestDump(t *testing.T) {
	result := handler.DumpSecuritySettings(models.ParameterMap{"path": gopath + `test\testdata\windows_utility_testdata\`})
	if result.Err != nil {
		t.Error(result.Err)
	}
}

func TestSecQuery(t *testing.T) {
	result := handler.SecuritySettingsQuery(models.ParameterMap{"path": gopath + `test\testdata\windows_utility_testdata\\secedit.cfg`, "valueName": "PasswordComplexity"})

	if result.Err != nil {
		t.Error(result.Err)
	}

	if result.Result != "0" {
		t.Error("Fehlgeschlagen:", result.Result)
	}
}

func TestSecQueryQuotationMarks(t *testing.T) {
	result := handler.SecuritySettingsQuery(models.ParameterMap{"path": gopath + `test\testdata\windows_utility_testdata\secedit.cfg`, "valueName": "NewAdministratorName"})

	if result.Err != nil {
		t.Error(result.Err)
	}

	if result.Result != "Administrator" {
		t.Error("Fehlgeschlagen:", result.Result)
	}
}

func TestSecQuerySID(t *testing.T) {
	result := handler.SecuritySettingsQuery(models.ParameterMap{"path": gopath + `test\testdata\windows_utility_testdata\secedit.cfg`, "valueName": "SeBackupPrivilege"})

	// SID: Administratoren, Sicherungsoperatoren
	if result.Result != "*S-1-5-32-544,*S-1-5-32-551" {
		t.Error("Fehlgeschlagen:", result.Result)
	}
}
