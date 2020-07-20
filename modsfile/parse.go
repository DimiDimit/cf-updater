package modsfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/elliotchance/orderedmap"
	"github.com/pkg/errors"
)

var empty struct{}

// Parse returns a slice of URLs and a slice of exclusions.
func Parse(prefix string, file io.Reader) (
	urlsSlice []string, excls []*regexp.Regexp, version string, err error) {
	urls := orderedmap.NewOrderedMap()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		default:
			line = prefix + line
			fallthrough

		case strings.HasPrefix(line, "https://"):
			if !strings.HasPrefix(line, prefix) {
				fmt.Printf("Warning: URL doesn't start with \"%v\": \"%v\"\n", prefix, line)
			}
			if _, ok := urls.Get(line); ok {
				return nil, nil, "", errors.New("duplicated URL: " + line)
			}
			urls.Set(line, empty)

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

	for url := urls.Front(); url != nil; url = url.Next() {
		urlsSlice = append(urlsSlice, url.Key.(string))
	}
	return
}

// ParseFile opens the file called fileName and calls Parse.
func ParseFile(prefix, fileName string) ([]string, []*regexp.Regexp, string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "error opening mods file")
	}
	defer file.Close()
	return Parse(prefix, file)
}
