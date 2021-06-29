package modulecontroller

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"github.com/pkg/errors"
	"runtime"
	s "strings"
)

// Versucht das Betriebssystem und seine Version zu bestimmen. Es wird nicht zwischen den unterschiedlichen
// Ausführungen von einzelnen Windows-Versionen (Home/Pro) etc. unterschieden.
func GetOS() (string, error) {
	var os string
	var err error

	Info(SeperateTitle("OS-Detector"))
	Info("Der OS-Detector wurde gestartet.")

	switch runtime.GOOS {
	case "windows":
		suffix := []string{"home", "pro", "enterprise", "enterpriseevaluation", "education", "essentials", "essentialsevaluation", "standard", "evaluation", "standardevaluation", "datacenter", "datacenterevaluation", "mobile", "foundation", "rt", "web", "ultimate", "professional", "homepremium", "homebasic", "starter", "business"}
		name, err := util.RegQuery(`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "ProductName")
		if err != nil {
			Err(err.Error())
			return "", err
		}
		name = s.ToLower(util.CompressString(name))
		Debug("Betriebssystem (roh): " + name)

		os = util.RemoveFromString(name, suffix)

	case "linux":
		name, err := util.ExecCommand("sed -n 3p /etc/os-release")
		if err != nil {
			Err(err.Error())
			return "", err
		}
		name = s.Trim(name, "\n")
		Debug("Betriebssystem (roh): " + name)
		os = s.ReplaceAll(name[s.LastIndex(name, "=")+1:], "\"", "")

	case "darwin":
		version, err := util.ExecCommand("sw_vers -productVersion")
		Debug("Betriebssystem (roh): MACOS " + version)
		if err != nil {
			Err(err.Error())
			return "", err
		}
		os = "macos" + version[0:s.LastIndex(version, ".")]

	default:
		Debug("Die von GO festgelegte Architektur ist " + runtime.GOOS + ". Unterstützt werden nur \"windows\", \"linux\" und \"darwin\".")
		err = errors.New("Das Betriebsystem wird nicht unterstützt.")
	}

	os = s.ToLower(os)
	if os != "" {
		Info("Der OS-Detector wurde mit folgendem Ergebnis beendet: " + os)
	}

	return os, err
}
