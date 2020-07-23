package cfwidget

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// DownloadURL returns a URL to a direct download of this file.
func (file *File) DownloadURL(parentID int, client *resty.Client) (string, error) {
	url, err := client.R().Get(fmt.Sprintf(
		"https://addons-ecs.forgesvc.net/api/v2/addon/%v/file/%v/download-url", parentID, file.ID))
	if err != nil {
		return "", errors.Wrap(err, "error getting download URL for mod")
	}
	return string(url.Body()), nil
}

// ActualName returns the actual name of a mod, because CurseForge replaces spaces with pluses.
func (file *File) ActualName() string {
	return strings.ReplaceAll(file.Name, " ", "+")
}
