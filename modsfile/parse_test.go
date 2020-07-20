package modsfile

import (
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const prefix = "https://www.curseforge.com/minecraft/mc-mods/"

func compile(t *testing.T, regex string) *regexp.Regexp {
	compiled, err := regexp.Compile(regex)
	if err != nil {
		t.Errorf("Error compiling regexp: %v", err)
	}
	return compiled
}

func regexpCompare(x, y *regexp.Regexp) bool {
	return x.String() == y.String()
}

func test(t *testing.T, file string,
	expURLs []string, expExcls []*regexp.Regexp, expVersion string, expErr bool) {
	resURLs, resExcls, resVersion, resErr := Parse(prefix, strings.NewReader(file))
	if !expErr && resErr != nil {
		t.Errorf("Expected success, got %v", resErr)
	} else if expErr && resErr == nil {
		t.Error("Expected error, but got success")
	}
	if diff := cmp.Diff(expURLs, resURLs); diff != "" {
		t.Errorf("URLs are different:\n%v", diff)
	}
	if diff := cmp.Diff(expExcls, resExcls, cmp.Comparer(regexpCompare)); diff != "" {
		t.Errorf("Exclusions are different:\n%v", diff)
	}
	if diff := cmp.Diff(expVersion, resVersion); diff != "" {
		t.Errorf("Versions are different:\n%v", diff)
	}
}

func TestURLs(t *testing.T) {
	test(t,
		`version 1.12.2

		 https://www.curseforge.com/minecraft/mc-mods/jei

		 https://www.curseforge.com/minecraft/mc-mods/shadowfacts-forgelin
		 https://www.curseforge.com/minecraft/mc-mods/dimitrodam-test`,
		[]string{
			"https://www.curseforge.com/minecraft/mc-mods/jei",
			"https://www.curseforge.com/minecraft/mc-mods/shadowfacts-forgelin",
			"https://www.curseforge.com/minecraft/mc-mods/dimitrodam-test",
		}, nil, "1.12.2", false)
}

func TestExcludes(t *testing.T) {
	test(t,
		`version 1.12.2

		 exclude ^OptiFine.*\.jar$
		 exclude ^Computronics.*\.jar$`,
		nil, []*regexp.Regexp{
			compile(t, "^OptiFine.*\\.jar$"),
			compile(t, "^Computronics.*\\.jar$"),
		}, "1.12.2", false)
}

func TestShortSyntax(t *testing.T) {
	test(t,
		`version 1.12.2

		 jei
		 shadowfacts-forgelin
		 dimitrodam-test`,
		[]string{
			"https://www.curseforge.com/minecraft/mc-mods/jei",
			"https://www.curseforge.com/minecraft/mc-mods/shadowfacts-forgelin",
			"https://www.curseforge.com/minecraft/mc-mods/dimitrodam-test",
		}, nil, "1.12.2", false)
}

func TestDifferentPrefixes(t *testing.T) {
	test(t,
		`version 1.12.2

		 https://mods.curse.com/mc-mods/minecraft/cofhcore
		 https://somedifferentprefix`,
		[]string{
			"https://mods.curse.com/mc-mods/minecraft/cofhcore",
			"https://somedifferentprefix",
		}, nil, "1.12.2", false)
}

func TestComments(t *testing.T) {
	test(t,
		`version 1.12.2

		 #comment

		 #
		 # comment with space`,
		nil, nil, "1.12.2", false)
}

func TestRegexpErrors(t *testing.T) {
	tests := map[string]string{
		"parentheses": "(",
		"quantifier":  "?",
		"multiple":    "(?",
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			test(t, "version 1.12.2\nexclude "+tc, nil, nil, "", true)
		})
	}
}

func TestDuplication(t *testing.T) {
	test(t,
		`version 1.12.2

		 https://www.curseforge.com/minecraft/mc-mods/jei
		 https://www.curseforge.com/minecraft/mc-mods/jei`, nil, nil, "", true)
	test(t,
		`version 1.12.2

		 https://www.curseforge.com/minecraft/mc-mods/jei
		 jei`, nil, nil, "", true)
}

func TestVersion(t *testing.T) {
	test(t,
		"https://www.curseforge.com/minecraft/mc-mods/jei", nil, nil, "", true)
	test(t,
		`version 1.12.2

		 https://www.curseforge.com/minecraft/mc-mods/jei

		 version 1.12.2`, nil, nil, "", true)
}

func TestMixed(t *testing.T) {
	test(t,
		`version 1.12.2

		 # Not dependencies of my mod
		 https://www.curseforge.com/minecraft/mc-mods/jei

		 # Dependencies of my mod
		 https://www.curseforge.com/minecraft/mc-mods/shadowfacts-forgelin
		 https://www.curseforge.com/minecraft/mc-mods/dimitrodam-test
		 
		 # We want to keep OptiFine and Computronics.
		 exclude ^OptiFine.*\.jar$
		 exclude ^Computronics.*\.jar$
		 
		 # Thermal mods
		 cofh-core
		 cofh-world
		 thermal-foundation
		 thermal-expansion`,
		[]string{
			"https://www.curseforge.com/minecraft/mc-mods/jei",
			"https://www.curseforge.com/minecraft/mc-mods/shadowfacts-forgelin",
			"https://www.curseforge.com/minecraft/mc-mods/dimitrodam-test",
			"https://www.curseforge.com/minecraft/mc-mods/cofh-core",
			"https://www.curseforge.com/minecraft/mc-mods/cofh-world",
			"https://www.curseforge.com/minecraft/mc-mods/thermal-foundation",
			"https://www.curseforge.com/minecraft/mc-mods/thermal-expansion",
		}, []*regexp.Regexp{
			compile(t, "^OptiFine.*\\.jar$"),
			compile(t, "^Computronics.*\\.jar$"),
		}, "1.12.2", false)
}
