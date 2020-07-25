package twitchapi

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type files struct {
	Files []File
}

func findFileByID(files []File, id int) (tfile *File, seen bool) {
	for _, file := range files {
		if file.ID == id {
			tfile = &file
			seen = true
			break
		}
	}
	return
}

// LatestDownloadForVersion returns the latest download for a version or an error if no such download exists.
// It prioritizes releases over betas over alphas (though that is subject to change because of the Twitch API).
func (info *ModInfo) LatestDownloadForVersion(client *resty.Client, version string) (*File, error) {
	var id int
	seen := false
	for _, file := range info.GameVersionLatestFiles {
		if file.GameVersion == version {
			id = file.ProjectFileID
			seen = true
			break
		}
	}
	if !seen {
		return nil, fmt.Errorf("couldn't find a download of %v for version %v", info.Name, version)
	}

	file, seen := findFileByID(info.LatestFiles, id)
	if !seen {
		// This actually happens relatively often.
		resp, err := client.R().
			SetResult(files{}.Files).
			Get(fmt.Sprintf("https://%v/api/v2/addon/%v/files", Endpoint, info.ID))
		if err != nil {
			return nil, errors.Wrap(err, "error fetching downloads for "+info.Name)
		}
		files := resp.Result().(*[]File)
		file, seen = findFileByID(*files, id)
		if !seen {
			// Nothing we can do about it now.
			return nil, fmt.Errorf("couldn't find a download with ID %v for %v", id, info.Name)
		}
	}
	return file, nil
}

// ActualName returns the actual name of a mod, because CurseForge replaces spaces with pluses.
func (file *File) ActualName() string {
	return strings.ReplaceAll(file.FileName, " ", "+")
}
