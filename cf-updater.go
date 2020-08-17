package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/DimiDimit/cf-updater/v3/modsfile"
	"github.com/DimiDimit/cf-updater/v3/twitchapi"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

var empty struct{}

const userAgent = "Mozilla/5.0 (X11; Linux x86_64) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36"

type download struct {
	Info     *twitchapi.ModInfo
	Download *twitchapi.File
}

func run() error {
	dir := flag.String("dir", ".", "The directory where the mods are located")
	modsfileLocF := flag.String("modsfile", "%dir/mods.txt", "The mods file location, use %dir/ to make it relative to the mods directory")
	showUpToDate := flag.Bool("u2d", false, "List up to date mods, useful for debugging")
	hideKeptBack := flag.Bool("hidekb", false, "Hide mods which are kept back from upgrading")
	flag.Parse()

	modsfileLoc := *modsfileLocF
	if strings.HasPrefix(modsfileLoc, "%dir/") {
		modsfileLoc = filepath.Join(*dir, strings.TrimPrefix(modsfileLoc, "%dir/"))
	}
	mods, excls, version, err := modsfile.ParseFile(modsfileLoc)
	if err != nil {
		return err
	}
	ids := make([]int, len(mods))
	{
		i := 0
		for k := range mods {
			ids[i] = k
			i++
		}
	}

	if err := os.Chdir(*dir); err != nil {
		return errors.Wrap(err, "error entering mods directory")
	}

	fmt.Println("… Fetching info about the mods...")
	client, downloads := resty.New(), make(map[string]download)
	client.SetHeader("User-Agent", userAgent)
	modInfos, err := twitchapi.GetMultipleMods(client, ids)
	if err != nil {
		return err
	}
	{
		var g errgroup.Group
		var downloadsm sync.Mutex
		bar := progressbar.Default(int64(len(*modInfos)))
		for _, info := range *modInfos {
			info := info
			g.Go(func() error {
				defer barInc(bar)
				mod := mods[info.ID]
				file, err := info.LatestDownload(client, version, mod.ReleaseType, mod.ModVersion)
				if err != nil {
					return err
				}
				downloadsm.Lock()
				downloads[file.ActualName()] = download{info, file}
				downloadsm.Unlock()
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
	}

	remaining := make(map[string]struct{})
	{
		fmt.Println("🗑 Deleting old mods...")
		files, err := filepath.Glob("*.jar")
		if err != nil {
			return err
		}
		bar := progressbar.Default(int64(len(files)))
		for _, file := range files {
			kept, err := func() (kept bool, err error) {
				kept = true
				if _, ok := downloads[file]; ok {
					return
				}
				for _, excl := range excls {
					if excl.MatchString(file) {
						return
					}
				}
				fmt.Println("Deleting", file)
				err = os.Remove(file)
				if err != nil {
					return false, err
				}
				return false, nil
			}()
			if err != nil {
				return err
			} else if kept {
				remaining[file] = empty
			}
			barInc(bar)
		}
	}

	{
		fmt.Println("⟳ Synchronizing mods...")
		var g errgroup.Group
		bar := progressbar.Default(int64(len(downloads)))
		for name, download := range downloads {
			name, download := name, download
			g.Go(func() error {
				defer barInc(bar)
				if _, ok := remaining[name]; ok {
					if mods[download.Info.ID].ModVersion != -1 {
						if !*hideKeptBack {
							fmt.Println("←", download.Info.Name, "has been kept back.")
						}
					} else if *showUpToDate {
						fmt.Println("→", download.Info.Name, "is up to date.")
					}
					return nil
				}
				fmt.Printf("⤓ Downloading %v...\n", download.Info.Name)
				url := download.Download.DownloadURL
				file, err := os.Create(download.Download.ActualName())
				if err != nil {
					return err
				}
				defer file.Close()
				resp, err := http.Get(url)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				_, err = io.Copy(file, resp.Body)
				if err != nil {
					return err
				}
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func barInc(bar *progressbar.ProgressBar) {
	_ = bar.Add(1)
}
