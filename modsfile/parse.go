package modsfile

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/DimiDimit/cf-updater/v3/twitchapi"
	"github.com/pkg/errors"
)

// Mod represents additional information about a mod.
type Mod struct {
	ModVersion  int
	ReleaseType int
}

// DefaultReleaseType is the release type used when one isn't specified.
var DefaultReleaseType = twitchapi.ReleaseTypes["release"]

var fieldSeparator = regexp.MustCompile(`\s+`)

// Parse returns a map of Mods and a slice of exclusions.
func Parse(file io.Reader) (mods map[int]Mod, excls []*regexp.Regexp, versions []string, err error) {
	mods = make(map[int]Mod)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		default:
			fields := fieldSeparator.Split(line, -1)
			id, err := strconv.Atoi(fields[0])
			if err != nil {
				return nil, nil, nil, errors.New("invalid syntax: " + line)
			}
			if _, ok := mods[id]; ok {
				return nil, nil, nil, errors.New("duplicated ID: " + line)
			}

			nfields, modVersion, releaseType := len(fields), -1, DefaultReleaseType
			if nfields >= 2 {
				cf := fields[1]
				ver, err := strconv.Atoi(cf)
				if err == nil {
					modVersion = ver
				} else {
					rt, ok := twitchapi.ReleaseTypes[cf]
					if !ok {
						return nil, nil, nil, errors.New("unknown release type: " + cf)
					}
					releaseType = rt
				}
			}
			mods[id] = Mod{modVersion, releaseType}

		case strings.HasPrefix(line, "exclude "):
			regex, err := regexp.Compile(strings.TrimSpace(strings.TrimPrefix(line, "exclude")))
			if err != nil {
				return nil, nil, nil, errors.Wrap(err, "mods file exclude syntax error")
			}
			excls = append(excls, regex)

		case strings.HasPrefix(line, "version "):
			if versions != nil {
				return nil, nil, nil, errors.New("duplicated version statement: " + line)
			}
			for _, version := range fieldSeparator.Split(strings.TrimPrefix(line, "version "), -1) {
				versions = append(versions, strings.TrimSpace(version))
			}

		case line == "" || strings.HasPrefix(line, "#"):
			// Ignore empty lines and comments.
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, nil, errors.Wrap(err, "error reading mods file")
	}

	if versions == nil {
		return nil, nil, nil, errors.New("version statement missing")
	}

	return
}

// ParseFile opens the file named fileName and calls Parse.
func ParseFile(fileName string) (map[int]Mod, []*regexp.Regexp, []string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error opening mods file")
	}
	defer file.Close()
	return Parse(file)
}
