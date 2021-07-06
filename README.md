# Jungbusch-Auditorium

Dieses Repository enthält den Quellcode für das Jungbusch-Auditorium. Das Jungbusch-Auditorium ist ein Framework zum Erstellen Modularer System-Audits.

## Download der Binaries

Siehe [Releases](https://github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/releases)

## Quickstart

1. Binary herunterladen: [Download](https://github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/releases/latest/)

2. Im Pfad der Executable eine Datei mit dem Namen `audit.jba` erstellen

3. Eine Audit-Konfiguration einfügen: 

<details>
  <summary>Windows (Aufklappen)</summary>
  <p>
  
	/*
		Autor: Jungbusch Softwareschmiede
		Date: 06.07.2021
		Version: 1.0
		Anmerkungen: keine
		Vorlage: CIS Microsoft Windows 10 Enterprise (Release 20H2 or older) Benchmark | v1.10.0 - 27-01-2021
	*/
	{
		stepid: "1"
		desc: "Ensure 'Interactive logon: Do not require CTRL+ALT+DEL' is set to 'Disabled' (Automated)"
		module: "RegistryQuery"

		key: "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System"
		value: "DisableCAD"
		passed: if("%result% == '0'")
	},
	
	// Check if the Group Policy Template 'AdmPwd.admx' and the language file 'AdmPwd.adml' is present
	{
		stepid: "2"
		desc: "Check if template 'AdmPwd.admx/adml' is present."
		module: "IsGPTemplatePresent"
		templateName: "AdmPwd.admx/adml"
		
		passed: if("%result%")
		%templatePresent% = %passed%

		{
			condition: if("%templatePresent%")
			stepid: "2.1"
			desc: "Ensure 'Password Settings: Password Length' is set to 'Enabled: 15 ormore' (Automated)"
			module: "RegistryQuery"

			key: "HKEY_LOCAL_MACHINE\SOFTWARE\Policies\Microsoft Services\AdmPwd"
			value: "PasswordLengt"
			passed: if("%result% >= '15'")
		},
	},
	
  </p>
</details>

<details>
  <summary>Linux (Aufklappen)</summary>
  
  <p>
  
	/*
		Autor: Jungbusch Softwareschmiede
		Date: 06.07.2021
		Version: 1.0
		Anmerkungen: Keine
		Vorlage: CIS Red Hat Enterprise Linux 8 Benchmark | v1.0.0 09-30-2019
	*/
	{
		stepid: "1"
		desc: "Ensure mounting of cramfs filesystems is disabled"
		module: "Modprobe"
		name: "cramfs"
		passed: if("%result%")
	},
	{
		stepid: "2"
		desc: "Ensure /tmp is configured"
		module: "CheckPartition"
		grep: "\s/tmp\s"
		passed: if("%result%.includes('tmpfs on /tmp')")
	},
	
</details>

4. (Optional) Eine Konfigurationsdatei erstellen via Commandline-Parameter `-createDefault`. Es wird eine config.ini-Datei mit den Default-Werten erstellt.

5. Die Executeable per Commandline (ohne weitere Commandline-Parameter) ausführen.

## Übersicht der Jungbusch-Repositories

Siehe [Jungbusch-Overview](https://github.com/Jungbusch-Softwareschmiede/jungbusch-overview)

## Handbuch

Siehe [Jungbusch-Manual](https://github.com/Jungbusch-Softwareschmiede/jungbusch-manual)

## Dokumentation

Siehe [Jungbusch-Documentation](https://github.com/Jungbusch-Softwareschmiede/jungbusch-documentation)

## Roadmap

1. Übersetzen der Doku/des Handbuchs in Englisch

2. Vereinheitlichen und Auslagern der Konfiguration in ein Go-Module

3. Die Möglichkeit, Variablennamen in der Audit-Konfiguration zu escapen, da sonst Windows-Umgebungspfade als Variablen erkannt werden

4. Variablen generell (Parser sollte mehr Logik bzgl. der Variablen übernehmen, Interpreter so wenig wie möglich)

5. Spezifisch: Überarbeiten der Logik zur Verwendung von Variablen in Parametern (Relatiert: Punkt 3, 4)

## About

Dieses Projekt wurde im Rahmen des Projektsemesters im Studiengang [Cybersecurity](https://www.hs-mannheim.de/studieninteressierte/unsere-studiengaenge/bachelorstudiengaenge/cyber-security.html) an der [Hochschule Mannheim](https://www.hs-mannheim.de/) im Zeitraum von 03/2021 - 07/2021 entwickelt.

### Mitwirkende 

[Christian Höfig](https://github.com/cookieChrissi) 

[Tim Philipp](https://github.com/TimPhi) 
 
[Marius Schmalz](https://github.com/ByteSizedMarius) 
 
[Felix Klör](https://github.com/prefixFelix) 
 
[Tobias Nöth](https://github.com/Tobias01101110) 
 
[Lukas Hagmaier](https://github.com/Lucky-180) 

