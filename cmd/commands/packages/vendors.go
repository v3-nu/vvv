package packages

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type Vendor struct {
	Name   string
	Prefix string

	// Default: -y
	ConfirmOption string

	Packages string

	// Default: {{ .Prefix }} install {{ .ConfirmOption }} {{ .Packages }}
	ListInstalled string

	// Default: {{ .Prefix }} {{ .ConfirmOption }} update {{ .Packages }}
	Update string

	// Default: {{ .Prefix }} {{ .ConfirmOption }} upgrade {{ .Packages }}
	Upgrade string

	// Default: {{ .Prefix }} {{ .ConfirmOption }} upgrade
	UpgradeAll string

	// Default: {{ .Prefix }} {{ .ConfirmOption }} install {{ .Packages }}
	Install string

	// Default: {{ .Prefix }} {{ .ConfirmOption }} remove {{ .Packages }}
	Remove string

	// Default: {{ .Prefix }} {{ .ConfirmOption }} search {{ .Packages }}
	Search string

	// Default: {{ .Prefix }} {{ .ConfirmOption }} info {{ .Packages }}
	Info string
}

func (v Vendor) TemplateString(tpl string) string {
	tmpl, err := template.New("command").Parse(tpl)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, v)

	return buf.String()
}

func (v Vendor) ListInstalledCommand() string {
	Command := "{{ .Prefix }} install {{ .ConfirmOption }} {{ .Packages }}"
	if v.ListInstalled != "" {
		Command = v.ListInstalled
	}

	return Command
}

func (v Vendor) UpdateCommand() string {
	Command := "{{ .Prefix }} {{ .ConfirmOption }} update {{ .Packages }}"
	if v.Update != "" {
		Command = v.Update
	}

	return Command
}

func (v Vendor) UpgradeCommand() string {
	Command := "{{ .Prefix }} {{ .ConfirmOption }} upgrade {{ .Packages }}"
	if v.Upgrade != "" {
		Command = v.Upgrade
	}

	return Command
}

func (v Vendor) UpgradeAllCommand() string {
	Command := "{{ .Prefix }} {{ .ConfirmOption }} upgrade"
	if v.UpgradeAll != "" {
		Command = v.UpgradeAll
	}

	return Command
}

func (v Vendor) InstallCommand() string {
	Command := "{{ .Prefix }} {{ .ConfirmOption }} install {{ .Packages }}"
	if v.Install != "" {
		Command = v.Install
	}

	return Command
}

func (v Vendor) RemoveCommand() string {
	Command := "{{ .Prefix }} {{ .ConfirmOption }} remove {{ .Packages }}"
	if v.Remove != "" {
		Command = v.Remove
	}

	return Command
}

func (v Vendor) SearchCommand() string {
	Command := "{{ .Prefix }} {{ .ConfirmOption }} search {{ .Packages }}"
	if v.Search != "" {
		Command = v.Search
	}

	return Command
}

func (v Vendor) InfoCommand() string {
	Command := "{{ .Prefix }} {{ .ConfirmOption }} info {{ .Packages }}"
	if v.Info != "" {
		Command = v.Info
	}

	return Command
}

var VendorMap = map[string]Vendor{
	"upt": {
		Name:   "upt",
		Prefix: "upt",
	},
	"apk": {
		Name:          "apk",
		Prefix:        "apk",
		ConfirmOption: "",
		Install:       "{{ .Prefix }} add {{ .Packages }}",
		Remove:        "{{ .Prefix }} del {{ .Packages }}",
		ListInstalled: "{{ .Prefix }} list --installed",
	},
	"apt": {
		Name:          "apt",
		Prefix:        "apt",
		ConfirmOption: "-y",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} install --only-upgrade {{ .Packages }}",
		Info:          "{{ .Prefix }} show {{ .Packages }}",
		Update:        "apt update",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} dist-upgrade",
		ListInstalled: "{{ .Prefix }} list --installed",
	},
	"brew": {
		Name:          "brew",
		Prefix:        "brew",
		ConfirmOption: "",
		Remove:        "{{ .Prefix }} {{ .ConfirmOption }} uninstall {{ .Packages }}",
		ListInstalled: "{{ .Prefix }} list",
	},
	"cards": {
		Name:          "cards",
		Prefix:        "cards",
		ConfirmOption: "",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} install --upgrade {{ .Packages }}",
		Update:        "{{ .Prefix }} {{ .ConfirmOption }} sync",
		ListInstalled: "{{ .Prefix }} list",
	},
	"choco": {
		Name:          "choco",
		Prefix:        "choco",
		ConfirmOption: "-y",
		Remove:        "{{ .Prefix }} {{ .ConfirmOption }}  uninstall {{ .Packages }}",
		Update:        "echo 'Not available'",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} upgrade all",
		ListInstalled: "{{ .Prefix }} list",
	},
	"dnf": {
		Name:          "dnf",
		Prefix:        "dnf",
		ConfirmOption: "-y",
		Update:        "{{ .Prefix }} check-update",
		ListInstalled: "{{ .Prefix }} list --installed",
	},
	"emerge": {
		Name:          "emerge",
		Prefix:        "emerge",
		ConfirmOption: "",
		Install:       "{{ .Prefix }} {{ .Packages }}",
		Remove:        "{{ .Prefix }} --depclean {{ .Packages }}",
		Upgrade:       "{{ .Prefix }} --update {{ .Packages }}",
		Search:        "{{ .Prefix }} --search {{ .Packages }}",
		Info:          "{{ .Prefix }} --info {{ .Packages }}",
		Update:        "{{ .Prefix }} --sync",
		UpgradeAll:    "{{ .Prefix }} -vuDN @world",
		ListInstalled: "qlist -Iv",
	},
	"eopkg": {
		Name:          "eopkg",
		Prefix:        "eopkg",
		ConfirmOption: "-y",
		Update:        "{{ .Prefix }} {{ .ConfirmOption }} update-repo",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} upgrade",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} list-installed",
	},
	"flatpak": {
		Name:          "flatpak",
		Prefix:        "flatpak",
		ConfirmOption: "-y",
		Remove:        "{{ .Prefix }} {{ .ConfirmOption }} uninstall {{ .Packages }}",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} update {{ .Packages }}",
		Update:        "",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} update",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} list",
	},
	"guix": {
		Name:          "guix",
		Prefix:        "guix",
		ConfirmOption: "",
		Info:          "{{ .Prefix }} {{ .ConfirmOption }} show {{ .Packages }}",
		Update:        "{{ .Prefix }} {{ .ConfirmOption }} refresh",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} package --list-installed",
	},
	"nix-env": {
		Name:          "nix-env",
		Prefix:        "nix-env",
		ConfirmOption: "",
		Install:       "{{ .Prefix }} {{ .ConfirmOption }} --install {{ .Packages }}",
		Remove:        "{{ .Prefix }} {{ .ConfirmOption }} --uninstall {{ .Packages }}",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} --upgrade {{ .Packages }}",
		Search:        "{{ .Prefix }} {{ .ConfirmOption }} --qaP {{ .Packages }}",
		Info:          "{{ .Prefix }} {{ .ConfirmOption }} --qa --description {{ .Packages }}",
		Update:        "nix-channel --update",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} --upgrade {{ .Packages }}",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} --query --installed",
	},
	"opkg": {
		Name:          "opkg",
		Prefix:        "opkg",
		ConfirmOption: "",
		Search:        "{{ .Prefix }} {{ .ConfirmOption }} find {{ .Packages }}",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} list-installed {{ .Packages }}",
	},
	"pacman": {
		Name:          "pacman",
		Prefix:        "pacman",
		ConfirmOption: "--noconfirm",
		Install:       "{{ .Prefix }} {{ .ConfirmOption }} -S {{ .Packages }}",
		Remove:        "{{ .Prefix }} {{ .ConfirmOption }} -R -s {{ .Packages }}",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} -S {{ .Packages }}",
		Search:        "{{ .Prefix }} {{ .ConfirmOption }} -S -s {{ .Packages }}",
		Info:          "{{ .Prefix }} {{ .ConfirmOption }} -S -i {{ .Packages }}",
		Update:        "{{ .Prefix }} {{ .ConfirmOption }} -S -y",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} -S -y -u",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} -Q",
	},
	"pkg": {
		Name:          "pkg",
		Prefix:        "pkg",
		ConfirmOption: "--yes",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} install {{ .Packages }}",
		ListInstalled: "{{ .Prefix }} info --all",
	},
	"pkg(termux)": {
		Name:          "pkg(termux)",
		Prefix:        "pkg",
		ConfirmOption: "--yes",
		Remove:        "{{ .Prefix }} {{ .ConfirmOption }} uninstall {{ .Packages }}",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} {{ .ConfirmOption }} install {{ .Packages }}",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} list-installed",
		Info:          "{{ .Prefix }} {{ .ConfirmOption }} show {{ .Packages }}",
	},
	"pkgman": {
		Name:          "pkgman",
		Prefix:        "pkgman",
		ConfirmOption: "-y",

		Remove:        "{{ .Prefix }} {{ .ConfirmOption }} uninstall {{ .Packages }}",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} update {{ .Packages }}",
		Info:          "echo 'Not available'",
		Update:        "{{ .Prefix }} {{ .ConfirmOption }} refresh",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} update",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} search --installed-only --all",
	},
	"prt-get": {
		Name:          "prt-get",
		Prefix:        "prt-get",
		ConfirmOption: "",
		Update:        "ports -u",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} sysup",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} listinst",
	},
	"scoop": {
		Name:          "scoop",
		Prefix:        "scoop",
		ConfirmOption: "",
		Remove:        "{{ .Prefix }} {{ .ConfirmOption }} uninstall {{ .Packages }}",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} update {{ .Packages }}",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} update *",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} list",
	},
	"slackpkg": {
		Name:          "slackpkg",
		Prefix:        "slackpkg",
		ConfirmOption: "",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} upgrade-all",
		ListInstalled: "ls -1 /var/log/packages",
	},
	"snap": {
		Name:          "snap",
		Prefix:        "snap",
		ConfirmOption: "",
		Install:       "{{ .Prefix }} {{ .ConfirmOption }} install --classic {{ .Packages }}",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} refresh {{ .Packages }}",
		Search:        "{{ .Prefix }} find {{ .Packages }}",
		Info:          "{{ .Prefix }} info {{ .Packages }}",
		Update:        "echo 'Not available'",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} refresh",
		ListInstalled: "{{ .Prefix }} list",
	},
	"urpm": {
		Name:          "urpm",
		Prefix:        "urpm",
		ConfirmOption: "-y",
		Install:       "urpmi {{ .ConfirmOption }} {{ .Packages }}",
		Remove:        "urpme {{ .ConfirmOption }} {{ .Packages }}",
		Upgrade:       "urpmi {{ .ConfirmOption }} {{ .Packages }}",
		Search:        "urpmq {{ .ConfirmOption }} {{ .Packages }}",
		Info:          "urpmq -i {{ .Packages }}",
		Update:        "urpmi.update -a",
		UpgradeAll:    "urpmi --auto-update",
		ListInstalled: "rpm --query --all",
	},
	"winget": {
		Name:          "winget",
		Prefix:        "winget",
		ConfirmOption: "",
		Remove:        "{{ .Prefix }} {{ .ConfirmOption }} uninstall {{ .Packages }}",
		Info:          "{{ .Prefix }} {{ .ConfirmOption }} show {{ .Packages }}",
		Update:        "echo 'Not available'",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} upgrade --all",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} list",
	},
	"xbps": {
		Name:          "xbps",
		Prefix:        "xbps",
		ConfirmOption: "-y",
		Install:       "{{ .Prefix }}-install {{ .ConfirmOption }} {{ .Packages }}",
		Remove:        "{{ .Prefix }}-remove {{ .ConfirmOption }} {{ .Packages }}",
		Upgrade:       "{{ .Prefix }}-install --update {{ .ConfirmOption }} {{ .Packages }}",
		Search:        "{{ .Prefix }}-query -Rs {{ .Packages }}",
		Info:          "{{ .Prefix }}-query -RS {{ .Packages }}",
		Update:        "{{ .Prefix }}-install {{ .ConfirmOption }} --sync",
		UpgradeAll:    "{{ .Prefix }}-install {{ .ConfirmOption }} --update",
		ListInstalled: "{{ .Prefix }}-query --list-pkgs",
	},
	"yum": {
		Name:          "yum",
		Prefix:        "yum",
		ConfirmOption: "-y",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} update {{ .Packages }}",
		Update:        "{{ .Prefix }} {{ .ConfirmOption }} check-update",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} update",
		ListInstalled: "{{ .Prefix }} {{ .ConfirmOption }} list --installed",
	},
	"zypper": {
		Name:          "zypper",
		Prefix:        "zypper",
		ConfirmOption: "-y",
		Upgrade:       "{{ .Prefix }} {{ .ConfirmOption }} update {{ .Packages }}",
		Update:        "{{ .Prefix }} {{ .ConfirmOption }} refresh",
		UpgradeAll:    "{{ .Prefix }} {{ .ConfirmOption }} update",
		ListInstalled: "zypper search --installed-only",
	},
}

var OsVendorMap = map[string][]string{
	"windows":    {"scoop", "choco", "winget"},
	"macos":      {"brew", "port"},
	"ubuntu":     {"apt"},
	"debian":     {"apt"},
	"linuxmint":  {"apt"},
	"pop":        {"apt"},
	"deepin":     {"apt"},
	"elementary": {"apt"},
	"kali":       {"apt"},
	"raspbian":   {"apt"},
	"aosc":       {"apt"},
	"zorin":      {"apt"},
	"antix":      {"apt"},
	"devuan":     {"apt"},
	"bodhi":      {"apt"},
	"lxle":       {"apt"},
	"sparky":     {"apt"},
	// dnf
	"fedora":    {"dnf", "yum"},
	"redhat":    {"dnf", "yum"},
	"rhel":      {"dnf", "yum"},
	"amzn":      {"dnf", "yum"},
	"ol":        {"dnf", "yum"},
	"almalinux": {"dnf", "yum"},
	"rocky":     {"dnf", "yum"},
	"oubes":     {"dnf", "yum"},
	"centos":    {"dnf", "yum"},
	"qubes":     {"dnf", "yum"},
	"eurolinux": {"dnf", "yum"},
	// pacman
	"arch":        {"pacman"},
	"manjaro":     {"pacman"},
	"endeavouros": {"pacman"},
	"arcolinux":   {"pacman"},
	"garuda":      {"pacman"},
	"antergos":    {"pacman"},
	"kaos":        {"pacman"},
	// apk
	"alpine":     {"apk"},
	"postmarket": {"apk"},
	// zypper
	"opensuse":            {"zypper"},
	"opensuse-leap":       {"zypper"},
	"opensuse-tumbleweed": {"zypper"},
	// nix
	"nixos": {"nix-env"},
	// emerge
	"gentoo": {"emerge"},
	"funtoo": {"emerge"},
	// xps
	"void": {"xbps"},
	// urpm
	"mageia": {"urpm"},
	// slackpkg
	"slackware": {"slackpkg"},
	// eopkg
	"solus": {"eopkg"},
	// opkg
	"openwrt": {"opkg"},
	// cards
	"nutyx": {"cards"},
	// prt-get
	"crux": {"prt-get"},
	// pkg
	"freebsd":  {"pkg"},
	"ghostbsd": {"pkg"},
	// pkg(termux)
	"android": {"pkg(termux)"},
	// pkgman
	"haiku": {"pkgman"},
}

var RegexOsDistroId = regexp.MustCompile(`(?m)^ID=(.*)$`)
var RegexOsDistroLike = regexp.MustCompile(`(?m)^ID_LIKE=(.*)$`)

func BestGuessOs() string {
	osname := runtime.GOOS

	if OsVendorMap[osname] != nil {
		return osname
	}

	if osname == "linux" {
		if _, err := os.Stat("/etc/os-release"); err == nil {
			content, err := os.ReadFile("/etc/os-release")
			if err == nil {
				match := RegexOsDistroId.FindStringSubmatch(string(content))
				if len(match) > 1 {
					osname = match[1]
					if OsVendorMap[osname] != nil {
						return osname
					}
				}

				match = RegexOsDistroLike.FindStringSubmatch(string(content))
				if len(match) > 1 {
					osname = match[1]
					if OsVendorMap[osname] != nil {
						return osname
					}
				}
			}
		}
	}

	log.Fatalf("Could not determine OS %s, please specify a package manager manually", osname)

	return ""
}

func ExecutableExists(executable string) bool {
	_, err := exec.LookPath(executable)
	return err == nil
}

func BestGuessPackageManager() string {
	osname := BestGuessOs()

	if OsVendorMap[osname] != nil {
		for _, vendor := range OsVendorMap[osname] {
			if ExecutableExists(vendor) {
				return vendor
			}
		}
	}

	log.Fatalf("Could not determine package manager for OS %s (tried %s), please specify a package manager manually", osname, strings.Join(OsVendorMap[osname], ", "))

	return ""
}
