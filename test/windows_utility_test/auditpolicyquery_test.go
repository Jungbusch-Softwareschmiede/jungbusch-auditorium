package windows_utility

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"testing"
)

func TestAuditPolicyQuery(t *testing.T) {
	result := handler.AuditPolicyQuery(models.ParameterMap{"guid": "{0CCE9214-69AE-11D9-BED3-505054503030}"})

	if result.Err != nil {
		t.Error(result.Err)
	}

	// Basiert auf Systemsprache
	if result.Result != "Erfolg und Fehler" && result.Result != "Success and Failure" {
		t.Error("Fehlgeschlagen:", result.Result)
	}
}
