package modsfile

import (
	"regexp"
	"strings"
	"testing"

	"github.com/DimitrodAM/cf-updater/v3/twitchapi"
	"github.com/google/go-cmp/cmp"
)

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
	expMods map[int]Mod, expExcls []*regexp.Regexp, expVersion string, expErr bool) {
	resMods, resExcls, resVersion, resErr := Parse(strings.NewReader(file))
	if !expErr && resErr != nil {
		t.Errorf("Expected success, got %v", resErr)
	} else if expErr && resErr == nil {
		t.Error("Expected error, but got success")
	}
	if diff := cmp.Diff(expMods, resMods); diff != "" {
		t.Errorf("Mods are different:\n%v", diff)
	}
	if diff := cmp.Diff(expExcls, resExcls, cmp.Comparer(regexpCompare)); diff != "" {
		t.Errorf("Exclusions are different:\n%v", diff)
	}
	if diff := cmp.Diff(expVersion, resVersion); diff != "" {
		t.Errorf("Versions are different:\n%v", diff)
	}
}

func testN(t *testing.T, file string,
	expIDs []int, expExcls []*regexp.Regexp, expVersion string, expErr bool) {
	expMods := make(map[int]Mod)
	for _, id := range expIDs {
		expMods[id] = Mod{-1, DefaultReleaseType}
	}
	test(t, file, expMods, expExcls, expVersion, expErr)
}

func TestIDs(t *testing.T) {
	testN(t,
		`version 1.12.2

		 # jei
		 238222

		 # shadowfacts-forgelin
		 248453
		 # dimitrodam-test
		 321466`,
		[]int{238222, 248453, 321466}, nil, "1.12.2", false)
}

func TestExcludes(t *testing.T) {
	test(t,
		`version 1.12.2

		 exclude ^OptiFine.*\.jar$
		 exclude ^Computronics.*\.jar$`,
		map[int]Mod{}, []*regexp.Regexp{
			compile(t, "^OptiFine.*\\.jar$"),
			compile(t, "^Computronics.*\\.jar$"),
		}, "1.12.2", false)
}

func TestNonNumeric(t *testing.T) {
	test(t,
		`version 1.12.2

		 cofhcore
		 https://somedifferentprefix`,
		nil, nil, "", true)
}

func TestComments(t *testing.T) {
	test(t,
		`version 1.12.2

		 #comment

		 #
		 ## comment with space`,
		map[int]Mod{}, nil, "1.12.2", false)
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

		 238222
		 238222`, nil, nil, "", true)
}

func TestVersion(t *testing.T) {
	test(t,
		"238222", nil, nil, "", true)
	test(t,
		`version 1.12.2

		 238222

		 version 1.12.2`, nil, nil, "", true)
}

func TestModVersion(t *testing.T) {
	test(t,
		`version 1.12.2

		 292785 2639533`, map[int]Mod{
			292785: {2639533, DefaultReleaseType},
		}, nil, "1.12.2", false)
}

func TestReleaseType(t *testing.T) {
	test(t,
		`version 1.12.2

		 69162 release
		 239286 beta
		 238222 alpha`, map[int]Mod{
			69162:  {-1, twitchapi.ReleaseTypes["release"]},
			239286: {-1, twitchapi.ReleaseTypes["beta"]},
			238222: {-1, twitchapi.ReleaseTypes["alpha"]},
		}, nil, "1.12.2", false)
	test(t,
		`version 1.12.2

		69162 re1ease
		239286 Beta
		238222 ALPHA`, nil, nil, "", true)
}

func TestMixed(t *testing.T) {
	test(t,
		`version 1.12.2

		 ## Not dependencies of my mod
		 # jei
		 238222

		 ## Dependencies of my mod
		 # shadowfacts-forgelin
		 248453
		 # dimitrodam-test
		 321466
		 
		 ## We want to keep OptiFine and Computronics.
		 exclude ^OptiFine.*\.jar$
		 exclude ^Computronics.*\.jar$
		 
		 ## Thermal mods
		 # cofh-core
		 69162
		 # cofh-world
		 271384
		 # thermal-foundation
		 222880
		 # thermal-expansion
		 69163
		 
		 ## Miscellaneous mods
		 # vanillafix 1.0.10-99
		 292785 2639533
		 # cyclic
     239286 beta`,
		map[int]Mod{
			238222: {-1, DefaultReleaseType},
			248453: {-1, DefaultReleaseType},
			321466: {-1, DefaultReleaseType},
			69162:  {-1, DefaultReleaseType},
			271384: {-1, DefaultReleaseType},
			222880: {-1, DefaultReleaseType},
			69163:  {-1, DefaultReleaseType},
			292785: {2639533, DefaultReleaseType},
			239286: {-1, twitchapi.ReleaseTypes["beta"]},
		}, []*regexp.Regexp{
			compile(t, "^OptiFine.*\\.jar$"),
			compile(t, "^Computronics.*\\.jar$"),
		}, "1.12.2", false)
}
