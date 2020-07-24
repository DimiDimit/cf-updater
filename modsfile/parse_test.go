package modsfile

import (
	"regexp"
	"strings"
	"testing"

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
	expIDs []int, expExcls []*regexp.Regexp, expVersion string, expErr bool) {
	resIDs, resExcls, resVersion, resErr := Parse(strings.NewReader(file))
	if !expErr && resErr != nil {
		t.Errorf("Expected success, got %v", resErr)
	} else if expErr && resErr == nil {
		t.Error("Expected error, but got success")
	}
	if diff := cmp.Diff(expIDs, resIDs); diff != "" {
		t.Errorf("IDs are different:\n%v", diff)
	}
	if diff := cmp.Diff(expExcls, resExcls, cmp.Comparer(regexpCompare)); diff != "" {
		t.Errorf("Exclusions are different:\n%v", diff)
	}
	if diff := cmp.Diff(expVersion, resVersion); diff != "" {
		t.Errorf("Versions are different:\n%v", diff)
	}
}

func TestIDs(t *testing.T) {
	test(t,
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
		nil, []*regexp.Regexp{
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

func TestMixed(t *testing.T) {
	test(t,
		`version 1.12.2

		 # Not dependencies of my mod
		 # jei
		 238222

		 # Dependencies of my mod
		 # shadowfacts-forgelin
		 248453
		 # dimitrodam-test
		 321466
		 
		 # We want to keep OptiFine and Computronics.
		 exclude ^OptiFine.*\.jar$
		 exclude ^Computronics.*\.jar$
		 
		 # Thermal mods
		 # cofh-core
		 69162
		 # cofh-world
		 271384
		 # thermal-foundation
		 222880
		 # thermal-expansion
		 69163`,
		[]int{
			238222, 248453, 321466, 69162, 271384, 222880, 69163,
		}, []*regexp.Regexp{
			compile(t, "^OptiFine.*\\.jar$"),
			compile(t, "^Computronics.*\\.jar$"),
		}, "1.12.2", false)
}
