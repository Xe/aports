package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const packageTemplate = `
# Contributor: {{ .Owner }}
# Maintainer: {{ .Owner }}
_gopkg={{ .Gopkg }}
_goimport={{ .GoImport }}

pkgname={{ .GopkgName }}
pkgver={{ .Version }}
pkgrel=0
pkgdesc="{{ .Desc }}"
url="{{ .URL }}"
arch="{{ .Arch }}"
license="{{ .License }}"
depends="go gopath {{ .RuntimeDepends }}"
makedepends="{{ .MakeDepends }}"
install=""
source="$_goimport-$pkgver.tar.gz::{{ .TarballURL }}"
builddir="$srcdir/"

build() {
        cd "$builddir"
}

package() {
        cd "$builddir"
        mkdir -p "$pkgdir"/usr/share/gopath/src/"$_gopkg"
        cp -rf "$_goimport"-"$pkgver"/* "$pkgdir"/usr/share/gopath/src/"$_gopkg"
}

check() {
        cd "$builddir"
        export GOPATH="$builddir"

        mkdir -p src/"$_gopkg"
        cp -rf $_goimport-"$pkgver"/* src/"$_gopkg"
        go test -v "$_gopkg"/...
}
`

var (
	owner          = flag.String("owner", "", "maintainer and contributer of this package")
	goPkg          = flag.String("gopkg", "", "the go package")
	goImport       = flag.String("goimport", "", "the short name the package is imported as")
	version        = flag.String("version", "", "the version of the package")
	desc           = flag.String("desc", "", "single sentence about the package")
	pkgURL         = flag.String("url", "", "url for the homepage of the package")
	pkgArch        = flag.String("arch", "noarch", "package architecture")
	pkgLicense     = flag.String("license", "", "package license")
	runtimeDepends = flag.String("runtime-deps", "", "space separated list of extra packages to add at runtime")
	makeDepends    = flag.String("make-deps", "", "space separated list of extra packages to add at build time")
	tarballURL     = flag.String("tarball-url", "", "url for the tarball download")
)

func validateFlags() {
	lastFlag := ""
	switch {
	case *owner == "":
		lastFlag = "owner"
	case *goPkg == "":
		lastFlag = "gopkg"
	case *goImport == "":
		lastFlag = "goimport"
	case *version == "":
		lastFlag = "version"
	case *desc == "":
		lastFlag = "desc"
	case *pkgURL == "":
		lastFlag = "url"
	case *pkgLicense == "":
		lastFlag = "license"
	case *tarballURL == "":
		lastFlag = "tarball-url"
	}

	switch lastFlag {
	case "":
		return
	default:
	}

	log.Fatalf("please set %s", lastFlag)
}

func main() {
	flag.Parse()
	validateFlags()

	pkgName := makePackageName(*goPkg)

	err := os.Mkdir(pkgName, 0700)
	if err != nil {
		log.Fatal("making package directory", err)
	}

	fout, err := os.Create("./" + filepath.Join(pkgName, "APKBUILD"))
	if err != nil {
		log.Fatal("creating output file", err)
	}

	data := struct {
		Owner          string
		Gopkg          string
		GoImport       string
		GopkgName      string
		Version        string
		Desc           string
		URL            string
		Arch           string
		License        string
		RuntimeDepends string
		MakeDepends    string
		TarballURL     string
	}{
		Owner:          *owner,
		Gopkg:          *goPkg,
		GoImport:       *goImport,
		GopkgName:      pkgName,
		Version:        *version,
		Desc:           *desc,
		URL:            *pkgURL,
		Arch:           *pkgArch,
		License:        *pkgLicense,
		RuntimeDepends: *runtimeDepends,
		MakeDepends:    *makeDepends,
		TarballURL:     *tarballURL,
	}

	t := template.Must(template.New("APKBUILD").Parse(packageTemplate))
	err = t.Execute(fout, data)
	if err != nil {
		log.Fatal("writing template", err)
	}
}

func makePackageName(goPkg string) string {
	pkgName := strings.Replace(goPkg, ".", "-", -1)
	pkgName = strings.Replace(pkgName, "/", "-", -1)
	pkgName = "go-" + strings.ToLower(pkgName)

	return pkgName
}
