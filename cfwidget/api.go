package cfwidget

import (
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// Endpoint is the CFWidget API endpoint domain.
const Endpoint = "api.cfwidget.com"

// ModInfo represents information about a mod, as returned from CFWidget.
type ModInfo struct {
	ID      int
	Title   string
	Summary string
	Game    string
	Type    string
	URLs    struct {
		CurseForge string
		Project    string
	}
	Files     []File
	Downloads struct {
		// Monthly is always 0 because of API changes.
		Monthly int
		Total   int
	}
	Thumbnail   string
	Categories  []string
	CreatedAt   time.Time `json:"created_at,string"`
	Description string
	LastFetch   time.Time `json:"last_fetch,string"`
	Download    File
}

// File represents information about a file, as returned by CFWidget.
type File struct {
	ID         int
	URL        string
	Display    string
	Name       string
	Type       string
	Version    string
	FileSize   int
	Versions   []string
	Downloads  int
	UploadedAt time.Time `json:"uploaded_at,string"`
}

// GetModInfo returns the info for a mod by its URL.
// See https://cfwidget.com/#documentation:about for more information.
func GetModInfo(client *resty.Client, modURL string, gameVersion string) (*ModInfo, error) {
	apiURL, err := url.Parse(modURL)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing mod URL")
	}
	apiURL.Host = Endpoint
	apiURL.RawQuery = "version=" + gameVersion
	info, err := client.R().
		SetResult(ModInfo{}).
		Get(apiURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "error fetching mod info for "+modURL)
	}
	return info.Result().(*ModInfo), nil
}
