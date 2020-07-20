package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DimitrodAM/cf-updater/cfwidget"
	"github.com/DimitrodAM/cf-updater/modsfile"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
)

// Prefix is the URL prefix of the mods.
const Prefix = "https://www.curseforge.com/minecraft/mc-mods/"

var empty struct{}

func run() error {
	dir := flag.String("dir", ".", "The directory where the mods are located")
	flag.Parse()

	if err := os.Chdir(*dir); err != nil {
		return errors.Wrap(err, "error entering mods directory")
	}

	urls, excls, version, err := modsfile.ParseFile(Prefix, "mods.txt")
	if err != nil {
		return err
	}

	fmt.Println("… Fetching info about the mods...")
	bar := progressbar.Default(int64(len(urls)))
	client, downloads := resty.New(), make(map[string]*cfwidget.ModInfo)
	for _, url := range urls {
		info, err := cfwidget.GetModInfo(client, url, version)
		if err != nil {
			return err
		}
		downloads[info.Download.ActualName()] = info
		bar.Add(1)
	}

	fmt.Println("🗑 Deleting old mods...")
	files, err := filepath.Glob("*.jar")
	if err != nil {
		return err
	}
	bar = progressbar.Default(int64(len(files)))
	remaining := make(map[string]struct{})
	for _, file := range files {
		kept, err := func() (kept bool, err error) {
			defer bar.Add(1)
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
	}

	fmt.Println("⟳ Synchronizing mods...")
	bar = progressbar.Default(int64(len(downloads)))
	for name, info := range downloads {
		if err := func() error {
			defer bar.Add(1)
			if _, ok := remaining[name]; ok {
				fmt.Println("→", info.Title, "is up to date.")
				return nil
			}
			fmt.Printf("⤓ Downloading %v...\n", info.Title)
			err := browser.OpenURL(cfwidget.GetDownloadURL(info.Download.URL))
			if err != nil {
				return err
			}
			return nil
		}(); err != nil {
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
