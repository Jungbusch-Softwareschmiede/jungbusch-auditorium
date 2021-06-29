package acutil

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
)

// Diese Methode enthält in statischen Objekten die allgemeingültigen Parameter der Auditkonfiguration und deren Syntax.
func parameterSyntax() []ParameterSyntax {
	return []ParameterSyntax{ //
		{
			ParamName:        "stepid",
			ParamDescription: "Mit diesem Parameter muss für jeden Schritt eine einzigartige ID festgelegt werden.",
			ParamAlias:       []string{"stepidentification", "id", "identifikation"},
		},
		{
			ParamName:        "module",
			ParamDescription: "Mit diesem Parameter wird festgelegt, welches Modul ausgeführt werden soll.",
			ParamAlias:       []string{"mod", "modul"},
		},
		{
			ParamName:        "description",
			ParamDescription: "Mit diesem Parameter kann eine Beschreibung für den auszuführenden Schritt definiert werden.",
			ParamAlias:       []string{"desc", "beschreibung", "definition"},
		},
		{
			ParamName:        "condition",
			ParamDescription: "Mit diesem Parameter kann eine Bedingung festgelegt werden. Der aktuelle Schritt wird nur ausgeführt, wenn diese zutrifft.",
			ParamAlias:       []string{"cond", "bedingung"},
		},
		{
			ParamName:        "passed",
			ParamDescription: "Mit diesem Parameter kann eine Bedingung festgelegt werden. Der Schritt wird als Passed markiert, wenn diese zutrifft.",
			ParamAlias:       []string{"bestanden", "erfolgreich"},
		},
		{
			ParamName:        "requireselevatedprivileges",
			ParamDescription: "Mit diesem Parameter kann festgelegt werden, dass ein Modul Administrator, bzw. Root-Rechte benötigt.",
			ParamAlias:       []string{"requiresadmin", "requiresroot", "requiresprivileges", "rootonly", "adminonly"},
		},
		{
			ParamName:        "print",
			ParamDescription: "Mit diesem Parameter kann eine Konsolen-, sowie eine zusätzliche Ausgabe im Result gemacht werden.",
			ParamAlias:       []string{"ausgeben"},
		},
	}
}

// Diese Methode enthält den Syntax einer Bedingung.
func ifSyntax() []ParameterSyntax {
	return []ParameterSyntax{ //
		{
			ParamName:        "if",
			ParamDescription: "Markiert den Start einer Bedingung.",
			ParamAlias:       []string{"wenn"},
		},
	}
}
