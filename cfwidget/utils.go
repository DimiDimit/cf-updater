package cfwidget

import "strings"

// GetDownloadURL returns a URL to download a mod from a URL that points to the file.
func GetDownloadURL(url string) string {
	return strings.Replace(url, "files", "download", 1)
}

// ActualName returns the actual name of a mod, because CurseForge replaces spaces with pluses.
func (file *File) ActualName() string {
	return strings.ReplaceAll(file.Name, " ", "+")
}
