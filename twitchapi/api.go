package twitchapi

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// Endpoint is the Twitch API endpoint domain.
const Endpoint = "addons-ecs.forgesvc.net"

// ReleaseTypes is a map of human-readable release types to their IDs.
var ReleaseTypes = map[string]int{
	"release": 1,
	"beta":    2,
	"alpha":   3,
}

// ModInfo represents information about a mod, as returned from the Twitch API.
type ModInfo struct {
	ID      int
	Name    string
	Authors []struct {
		Name              string
		URL               string
		ProjectID         int
		ID                int
		ProjectTitleID    int
		ProjectTitleTitle string
		UserID            int
		TwitchID          int
	}
	Attachments []struct {
		ID           int
		ProjectID    int
		Description  string
		IsDefault    bool
		ThumbnailURL string
		Title        string
		URL          string
		Status       int
	}
	WebsiteURL    string
	GameID        int
	Summary       string
	DefaultFileID int
	DownloadCount float32
	LatestFiles   []File
	Categories    []struct {
		CategoryID int
		Name       string
		URL        string
		AvatarURL  string
		ParentID   int
		RootID     int
		ProjectID  int
		AvatarID   int
		GameID     int
	}
	Status            int
	PrimaryCategoryID int
	CategorySection   struct {
		ID                      int
		GameID                  int
		Name                    string
		PackageType             int
		Path                    string
		InitialInclusionPattern string
		ExtraIncludePattern     string
		GameCategoryID          int
	}
	Slug                   string
	GameVersionLatestFiles []struct {
		GameVersion       string
		ProjectFileID     int
		ProjectFileName   string
		FileType          int
		GameVersionFlavor string
	}
	IsFeatured         bool
	PopularityScore    float32
	GamePopularityRank int
	PrimaryLanguage    string
	GameSlug           string
	GameName           string
	PortalName         string
	DateModified       time.Time
	DateCreated        time.Time
	DateReleased       time.Time
	IsAvailable        bool
	IsExperimental     bool
}

// File represents information about a file, as returned by the Twitch API.
type File struct {
	ID              int
	DisplayName     string
	FileName        string
	FileDate        time.Time
	FileLength      int
	ReleaseType     int
	FileStatus      int
	DownloadURL     string
	IsAlternate     bool
	AlternateFileID int
	Dependencies    []struct {
		ID      int
		AddonID int
		Type    int
		FileID  int
	}
	IsAvailable bool
	Modules     []struct {
		FolderName  string
		Fingerprint int
		Type        int
	}
	PackageFingerprint  int
	GameVersion         []string
	SortableGameVersion []struct {
		GameVersionPadded      string
		GameVersion            string
		GameVersionReleaseDate time.Time
		GameVersionName        string
	}
	InstallMetadata            string
	Changelog                  string
	HasInstallScript           bool
	IsCompatibleWithClient     bool
	CategorySectionPackageType int
	RestrictProjectFileAccess  int
	ProjectStatus              int
	RenderCacheID              int
	FileLegacyMappingID        int
	ProjectID                  int
	ParentProjectFileID        int
	ParentFileLegacyMappingID  int
	FileTypeID                 int
	ExposeAsAlternative        bool
	PackageFingerprintID       int
	GameVersionDateReleased    time.Time
	GameVersionMappingID       int
	GameVersionID              int
	GameID                     int
	IsServerPack               bool
	ServerPackFileID           int
	GameVersionFlavor          string
}

type modInfos struct {
	ModInfos []*ModInfo
}

// GetModInfo returns the info for a mod by its ID.
// See https://twitchappapi.docs.apiary.io/#/reference/0/get-addon-info for more information.
func GetModInfo(client *resty.Client, modID int) (*ModInfo, error) {
	info, err := client.R().
		SetResult(ModInfo{}).
		Get(fmt.Sprintf("https://%v/api/v2/addon/%v", Endpoint, modID))
	if err != nil {
		return nil, errors.Wrap(err, "error fetching mod info for "+string(modID))
	}
	return info.Result().(*ModInfo), nil
}

// GetMultipleMods returns the info for multiple mods by their IDs.
// This is more efficient then calling GetModInfo in a loop.
// See https://twitchappapi.docs.apiary.io/#/reference/0/get-multiple-addons for more information.
func GetMultipleMods(client *resty.Client, modIDs []int) (*[]*ModInfo, error) {
	info, err := client.R().
		SetBody(modIDs).
		SetResult(modInfos{}.ModInfos).
		Post(fmt.Sprintf("https://%v/api/v2/addon", Endpoint))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error fetching mod info for %v", modIDs))
	}
	return info.Result().(*[]*ModInfo), nil
}
