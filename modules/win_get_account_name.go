// +build windows

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
	s "strings"
)

func (mh *MethodHandler) GetAccountNameInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "GetAccountName",
		ModuleDescription:   "GetAccountName gibt den Name des Gast oder Lokalen Accounts aus.",
		ModuleAlias:         []string{"getaccountname"},
		ModuleCompatibility: []string{"windows"},
		InputParams: ParameterSyntaxMap{
			"type": ParameterSyntax{
				ParamName:        "type",
				ParamAlias:       []string{"typ"},
				ParamDescription: "Accounttyp dessen Name abgefragt werden möchte. (Guest/User)",
			},
		},
	}
}

func (mh *MethodHandler) GetAccountName(params ParameterMap) (r ModuleResult) {
	switch s.ToLower(params["type"]) {
	case "guest":
		// Liste aller Nutzer mit Name und SID
		allUsers, err := util.ExecCommand("Get-WmiObject -Class Win32_UserAccount -Filter LocalAccount='True' | Select-Object Name, SID")
		if err != nil {
			r.Err = err
			return
		}

		r.Artifacts = append(r.Artifacts, Artifact{
			Name:  "Get-WmiObject -Class Win32_UserAccount -Filter LocalAccount='True' | Select-Object Name, SID",
			Value: "Liste aller Nutzer mit Name und SID:\n" + allUsers,
		})

		r.ResultRaw = allUsers

		lineSlice := s.Split(allUsers, "\n")
		for n := range lineSlice {
			if s.Contains(lineSlice[n], " ") {
				line := s.TrimSpace(lineSlice[n])

				back := s.TrimSpace(line[s.LastIndex(line, " "):])
				front := s.TrimSpace(line[:s.LastIndex(line, " ")])
				// Das Konto, welches auf 501 endet ist das Gast Konto
				if s.HasSuffix(back, "501") {
					r.Result = front
				}
			}
		}

	case "user":
		// Name des akutuellen lokalen Kontos mit hilfe von whoami finden
		username, err := util.ExecCommand("cmd.exe /c \"echo %USERNAME%\"")

		if err != nil {
			r.Err = err
			return
		}

		r.Artifacts = append(r.Artifacts, Artifact{
			Name:  "cmd.exe /c \"echo %USERNAME%\"",
			Value: "Name des aktuellen Users:\n" + username,
		})

		r.ResultRaw = username

		r.Result = s.TrimSpace(username)
	}
	return
}

func (mh *MethodHandler) GetAccountNameValidate(params ParameterMap) error {
	if s.ToLower(params["type"]) != "guest" && s.ToLower(params["type"]) != "user" {
		return errors.New("Der Accounttyp ist falsch geschrieben oder wird nicht unterstützt. Bitte 'Guest' oder 'User' verwenden.")
	}
	return nil
}
