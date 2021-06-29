package windows_utility_test

import (
	"fmt"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"testing"
)

func TestInt(t *testing.T) {
	// Typ: REG_DWORD, REG_QWORD
	res, err := util.RegQuery(`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Netlogon\Parameters`, "MaximumPasswordAge")
	if res == "30" {
		fmt.Println(res)
	} else {
		t.Error("Fehlgeschlagen:", res, err)
	}
}

func TestMultiString(t *testing.T) {
	// Typ: REG_MULTI_SZ
	res, err := util.RegQuery(`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\MSiSCSI`, "RequiredPrivileges")
	should := "SeAuditPrivilege,SeChangeNotifyPrivilege,SeCreateGlobalPrivilege,SeCreatePermanentPrivilege,SeImpersonatePrivilege,SeTcbPrivilege,SeLoadDriverPrivilege,"
	if err != static.ERROR_VALUE_NOT_FOUND && res == should {
		fmt.Println(res)
	} else {
		t.Error("Fehlgeschlagen:", res, err)
	}
}

func TestString(t *testing.T) {
	// Typ: REG_SZ, REG_EXPAND_SZ
	res, err := util.RegQuery(`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\MSiSCSI`, "Group")
	should := "iSCSI"
	if err != static.ERROR_VALUE_NOT_FOUND && res == should {
		fmt.Println(res)
	} else {
		t.Error("Fehlgeschlagen:", res, err)
	}
}
