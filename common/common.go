package common

import "regexp"

// FieldSeparator is what fields in the mods file are separated by.
var FieldSeparator = regexp.MustCompile(`\s+`)
