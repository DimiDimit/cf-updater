package twitchapi

import (
	"strings"
)

// ActualName returns the actual name of a mod, because CurseForge replaces spaces with pluses.
func (file *File) ActualName() string {
	return strings.ReplaceAll(file.FileName, " ", "+")
}
