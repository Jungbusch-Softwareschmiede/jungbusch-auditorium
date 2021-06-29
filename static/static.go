// Im Package static werden Konstanten und "Konstanten-ähnliche" Variablen gesammelt. Dies erlaubt das leichte anpassen
// vieler Konstanten und dient der Übersichtlichkeit.
package static

import (
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"os"
)

var (
	// Die Liste der verfügbaren Module. Wird ein neues Modul implementiert, muss es hier hinzugefügt werden.
	Modules = []string{
		// All
		"ExecuteCommand",
		"FileContent",
		"IsFile",
		"Grep",
		"Permissions",
		"Script",
		// Linux
		"Auditctl",
		"AwkScript",
		"BashScript",
		"CheckPartition",
		"Modprobe",
		"NftListRuleset",
		"Sshd",
		"Stat",
		"Sysctl",
		"Systemctl",
		// Linux / Darwin
		"IsInstalled",
		"IsNotInstalled",
		// Redhat
		"Authselect",
		// Windows
		"AuditPolicyQuery",
		"DumpSecuritySettings",
		"ExportInstalledSoftware",
		"GetAccountName",
		"GetUserSID",
		"GetWinEnv",
		"IsGPTemplatePresent",
		"RegistryQuery",
		"SecuritySettingsQuery",
	}

	// Legt alle Betriebssysteme fest, die sich in der Linux-Wildcard befinden.
	Linux = []string{
		"ubuntu",
		"debian",
		"centos",
		"rhel",     // RedHat
		"sles",     // Suse
		"sles_sap", // Suse SAP
	}

	// Legt alle Betriebssysteme fest, die sich in der Windows-Wildcard befinden.
	Windows = []string{
		"windows10",
		"windowsserver2012",
		"windowsserver2016",
		"windowsserver2019",
		"windows8",
		"windows7",
		"windowsxp",
	}

	// Legt alle Betriebssysteme fest, die sich in der Darwin-Wildcard befinden.
	Darwin = []string{
		"macos11.0",
		"macos10.15",
		"macos10.14",
	}

	// Über diesen String wird in der Jungbusch-Auditorium Global auf das Ergebnis des OS-Detectors zugegriffen.
	OperatingSystem = ""

	// Über diesen String wird in der Jungbusch-Auditorium Global auf den Pfad des
	// temporären Ordners zugegriffen. (Nur Windows)
	TempPath = ""

	// In diesem Bool wird gespeichert, ob das Jungbusch-Auditorium root, bzw. Administrator-Rechte hat.
	HasElevatedPrivileges = false

	// In dieser Variable wird gespeichert, mit welchen Permissions das Jungbusch-Auditorium Dateien erstellt.
	// Dieser Wert ändert sich, wenn das Auditorium Root-Rechte hat. Siehe: privilege.HasRootPrivileges
	CREATE_FILE_PERMISSIONS os.FileMode = 0644

	// In dieser Variable wird gespeichert, mit welchen Permissions das Jungbusch-Auditorium Ordner erstellt.
	// Dieser Wert ändert sich, wenn das Auditorium Root-Rechte hat. Siehe: privilege.HasRootPrivileges
	CREATE_DIRECTORY_PERMISSIONS os.FileMode = 0755

	// Speichert den Error, der geworfen wird, wenn ein Registry-Key nicht gefunden wurde. (Windows)
	ERROR_KEY_NOT_FOUND = errors.New("Der angegebene Key konnte nicht gefunden werden!")

	// Speichert den Error, der geworfen wird, wenn eine Registry-Value nicht gefunden wurde. (Windows)
	ERROR_VALUE_NOT_FOUND = errors.New("Die angegebene Value konnte nicht gefunden werden!")

	ProgressBar *progressbar.ProgressBar
)

const (
	CURRENT_VERSION = "V1.0 (Letzte Änderung: 29.06.2021)"
	PATH_SEPERATOR  = string(os.PathSeparator)

	// Setzt Default-Werte der Programmkonfiguration
	EXPECTED_CONFIG_PATH      string = `./config.ini`
	EXPECTED_AUDIT_PATH       string = `./audit.jba`
	DEFAULT_OUTPUT_PATH       string = `./`
	DEFAULT_LOG_VERBOSITY     int    = 3
	DEFAULT_CONSOLE_VERBOSITY int    = 2

	// Legt fest, welche Methodennamen das Jungbusch-Auditorium in Modulen erwartet.
	// Werden diese Werte angepasst, müssen auch alle bestehenden Module angepasst werden.
	MODULE_INITIALIZER_SUFFIX = "Init"
	MODULE_EXECUTOR_SUFFIX    = ""
	MODULE_VALIDATOR_SUFFIX   = "Validate"

	// Legt den Name der zu erstellenden Log-Datei fest.
	LOG_NAME = "Jungbusch-Auditorium"

	// Legt den Name des zu erstellenden Output-Ordners fest.
	OUTPUT_FOLDER_NAME      = "jba-result"
	OUTPUT_TIMESTAMP_FORMAT = "02-01-06T150405" // Siehe Time.Format

	// Legt fest, welche Strings von den Modulen als Anzeige des Fortschritts ausgegeben werden
	PASSED_OUTPUT       = "Test '%v' PASSED"
	NOTPASSED_OUTPUT    = "Test '%v' NOT PASSED"
	UNSUCCESSFUL_OUTPUT = "Test '%v' UNSUCCESSFUL"
	NOTEXECUTED_OUTPUT  = "Test '%v' not executed"

	// Reg-Ex
	REG_VARIABLE_NAME = "^[A-Za-z0-9_]+$"    // Zeichensatz der in Variablen erlaubt ist
	REG_VARIABLE      = "%([a-zA-Z0-9_]*?)%" // Anhand von diesem Zeichnsatz werden Variablen identifiziert
	REG_MODULE_NAME   = "^[A-Za-z0-9_]+$"    // Zeichensatz der in Modul-Namen/Parameter-Namen oder Aliasen enthalten sein darf

	// Parser Errors
	COMMENTBLOCK_NEVER_CLOSED                = "Ein Kommentarblock wurde nicht ordnungsgemäß geschlossen."
	AUDIT_CONFIG_EMPTY                       = "Die angegebene Audit-Konfigurationsdatei ist leer."
	INVALID_LINE                             = "Die Zeile ist ungültig. Möglicher Grund: "
	MODULE_MISSING_OPENING_BRACKET           = "Die öffnende geschweifte Klammer fehlt."
	MODULE_MISSING_CLOSING_BRACKET           = "Die schließende geschweifte Klammer fehlt."
	MODULE_MISSING_CLOSING_BRACKET_COMMA     = "Die schließende geschweifte Klammer ist ungültig. Möglicherweise fehlt das Komma."
	MODULE_EMPTY                             = "Der Audit-Schritt ist leer."
	MODULE_INVALID_EXPRESSION                = "Ungültiger Ausdruck."
	VARIABLE_INVALID                         = "Der Syntax der Variable ist nicht gültig. Möglicher Grund: "
	VARIABLE_INVALID_NAME                    = "Der Name der Variable ist nicht gültig. Möglicher Grund: "
	VARIABLE_INVALID_VALUE                   = "Der Syntax des Werts der Variable ist nicht gültig."
	MODULE_INVALID_MODULE_NAME               = "Der Name des Moduls ist nicht gültig. Möglicher Grund: "
	MODULE_VALUE_INVALID_MULTILINE           = "Der Wert des Multiline-Parameters ist nicht gültig. Möglicher Grund: "
	MODULE_VALUE_INVALID                     = "Der Wert des Parameters ist nicht gültig. Möglicher Grund: "
	MODULE_VALUE_ALREADY_SET                 = "Dieser Wert darf pro Audit-Schritt nur einmalig gesetzt werden."
	MODULE_MISSING_PARAMETER                 = "Einem Audit-Schritt fehlen folgende zwingend zu setzenden Parameter: "
	MODULE_NOT_FOUND                         = "Das angegebene Modul wurde nicht gefunden. Mögliche Gründe: Das Modul wurde nicht als kompatibel markiert oder das Modul ist nicht mit der aktuellen Architektur kompatibel."
	MODULE_REQUIRED_MODULE_PARAMETER_NOT_SET = "Folgender Modul-Parameter muss zwingend angegeben werden: "
	MODULE_NO_TEXT_BEHIND_MULTILINE          = "Hinter den schließenden Backticks von Multiline-Parametern ist kein Text (Kommentare o.ä.) erlaubt."
	PARAMETER_EMPTY                          = "Der Wert des Parameters ist leer."
	MODULE_INVALID_PARAMETER                 = "Der Parameter ist für dieses Modul nicht definiert."
	GLOBAL_MODULES_NOT_ALLOWED_NESTED        = "Globale Audit-Schritte dürfen keine verschachtelten Audit-Schritte haben. In folgendem Schritt ist ein Fehler aufgetreten: (ID): "
	GLOBAL_MODULE_NESTED                     = "Globale Audit-Schritte dürfen nicht verschachtelt sein. In folgendem Schritt ist ein Fehler aufgetreten: (ID): "
	GLOBAL_VARIABLE_READONLY                 = "Globale Variablen dürfen nur in globalen Audit-Schritten überschrieben werden. In folgendem Schritt ist ein Fehler aufgetreten (ID): "

	// Parser Errors: Gründe
	TOO_MANY_QUOTATIONMARKS                  = "Zu viele Anführungszeichen."
	TOO_MANY_TICKS                           = "Zu viele Backticks."
	MISSING_TICKS                            = "Fehlende Backticks."
	MISSING_QUOTATIONMARKS                   = "Fehlende Anführungszeichen."
	VARIABLE_MISSING_PERCENTAGE              = "Der Variablenname muss von Prozentzeichen umschlossen sein."
	INVALID_CHARACTERS                       = "Es wurden nicht-erlaubte Zeichen verwendet."
	INVALID_SYNTAX_OR_MISSING_QUOTATIONMARKS = "Der Syntax der Bedingung ist nicht valide oder fehlende Anführungszeichen."
	INVALID_VALUE_FOR_BOOL                   = "Der Boolean-Wert ist ungültig. (true/false)"
	ID_ALREADY_USED                          = "IDs müssen einzigartig sein."
	ENVIRONMENT_VARIABLES_READONLY           = "Umgebungsvariablen dürfen nicht überschrieben werden."
	VARIABLE_ALREADY_SET                     = "Diese Variable wurde in diesem Modul bereits gesetzt."
	ROW_IS_END_OF_MODULE                     = " Anmerkung: Die angegebene Zeile ist die schließende Klammer des betreffenden Moduls."
)
