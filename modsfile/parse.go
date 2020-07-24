package modsfile

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/elliotchance/orderedmap"
	"github.com/pkg/errors"
)

var empty struct{}

// Parse returns a slice of IDs and a slice of exclusions.
func Parse(file io.Reader) (idsSlice []int, excls []*regexp.Regexp, version string, err error) {
	ids := orderedmap.NewOrderedMap()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		default:
			id, err := strconv.Atoi(line)
			if err != nil {
				return nil, nil, "", errors.New("invalid syntax: " + line)
			}
			if _, ok := ids.Get(id); ok {
				return nil, nil, "", errors.New("duplicated ID: " + line)
			}
			ids.Set(id, empty)

		case strings.HasPrefix(line, "exclude "):
			regex, err := regexp.Compile(strings.TrimSpace(strings.TrimPrefix(line, "exclude")))
			if err != nil {
				return nil, nil, "", errors.Wrap(err, "mods file exclude syntax error")
			}
			excls = append(excls, regex)

		case strings.HasPrefix(line, "version "):
			if version != "" {
				return nil, nil, "", errors.New("duplicated version statement: " + line)
			}
			version = strings.TrimSpace(strings.TrimPrefix(line, "version"))

		case line == "" || strings.HasPrefix(line, "#"):
			// Ignore empty lines and comments.
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, "", errors.Wrap(err, "error reading mods file")
	}

	if version == "" {
		return nil, nil, "", errors.New("version statement missing")
	}

	for id := ids.Front(); id != nil; id = id.Next() {
		idsSlice = append(idsSlice, id.Key.(int))
	}
	return
}

// ParseFile opens the file called fileName and calls Parse.
func ParseFile(fileName string) ([]int, []*regexp.Regexp, string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "error opening mods file")
	}
	defer file.Close()
	return Parse(file)
}
