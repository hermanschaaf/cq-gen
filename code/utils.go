package code

import (
	"path/filepath"
	"regexp"
	"strings"
)

// PkgAndType takes a string in the form github.com/package/blah.Type and splits it into package and type
func PkgAndType(name string) (string, string) {
	parts := strings.Split(name, ".")
	if len(parts) == 1 {
		return "", name
	}

	return strings.Join(parts[:len(parts)-1], "."), parts[len(parts)-1]
}

// PkgName takes an import string in the form github.com/package/blah and returns only the name, e.g. blah
func PkgName(pkg string) string {
	parts := strings.Split(pkg, "/")
	if len(parts) == 1 {
		return pkg
	}
	return parts[len(parts)-1]
}

var modsRegex = regexp.MustCompile(`^(\*|\[\])*`)

// NormalizeVendor takes a qualified package path and turns it into normal one.
// eg .
// github.com/foo/vendor/github.com/99designs/gqlgen/graphql becomes
// github.com/99designs/gqlgen/graphql
func NormalizeVendor(pkg string) string {
	modifiers := modsRegex.FindAllString(pkg, 1)[0]
	pkg = strings.TrimPrefix(pkg, modifiers)
	parts := strings.Split(pkg, "/vendor/")
	return modifiers + parts[len(parts)-1]
}

var invalidPackageNameChar = regexp.MustCompile(`[^\w]`)

func SanitizePackageName(pkg string) string {
	return invalidPackageNameChar.ReplaceAllLiteralString(filepath.Base(pkg), "_")
}
