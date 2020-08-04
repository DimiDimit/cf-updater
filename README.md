# CurseForge Updater

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/mod/github.com/DimitrodAM/cf-updater/v3)
![Build](https://github.com/DimitrodAM/cf-updater/workflows/Build/badge.svg)
![Test and Lint](https://github.com/DimitrodAM/cf-updater/workflows/Test%20and%20Lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/DimitrodAM/cf-updater)](https://goreportcard.com/report/github.com/DimitrodAM/cf-updater)

A tool to update [CurseForge](https://www.curseforge.com) mods, written in [Go](https://golang.org). Currently only tested with [Minecraft](https://www.curseforge.com/minecraft/mc-mods).

## Compatibility

### Compatibility with v2

Version 3's mods file **is** compatible with version 2! But note that the API is **not**.

### Compatibility with v1

Version 2 is **not** compatible with version 1! You must **refactor your mods file** or you'll get syntax errors!

## Installation and upgrading

### Installing on Windows

To install it on Windows, download it from the [Releases](https://www.github.com/DimitrodAM/cf-updater/releases) page and save it into your mods folder (e.g. for Minecraft it's `.minecraft/mods` or `.minecraft/mods/<version>`).

### Installing with `go get`

First, [install Go](https://golang.org/doc/install).
Then run:

```sh
go get -u github.com/DimitrodAM/cf-updater
```

### Upgrading

The program has **no** built-in update checker or updater, so you should **watch the repository** in “Releases only” mode to get notified of updates.

The upgrade process is the same as the installation.

## Setup

**Backup your mods before doing anything else!**

Unfortunately, some setup is first required. You must put the IDs of all of your mods in a file called `mods.txt`. This can be tedious if you've got a lot of mods, but I don't see a way around it. If you have an idea, please submit an [issue](https://www.github.com/DimitrodAM/cf-updater/issues) or a [pull request](https://www.github.com/DimitrodAM/cf-updater/pulls)!

**Treat this file as the single source of mods**, because the tool will **delete** any mods not in the mods file! Files with any other extension (e.g. `.bak`) do not count as mods and will **not** be deleted.

To find the ID of a mod, look at the `Project ID` in the `About Project` panel on the right of its CurseForge page.

### Example

Here's an example mods file:

```
version 1.12.2

## Some mods
# jei
238222
# shadowfacts-forgelin
248453
# dimitrodam-test
321466
```

`version` is the version of the game that the mods are for and is **required**. Lines starting with `#` are comments.

To use **another modloader like Fabric**, you should specify multiple versions (case-sensitive):

```
version 1.16.1 Fabric
```

### Versions and types

Sometimes you may want to keep a mod on a certain version or use a different release type (e.g. alpha, beta or release). You can do so like this:

```
version 1.12.2

# cyclic
239286 beta
# vanillafix 1.0.10-99
292785 2639533
```

Version IDs can be found in their URLs.

### Mods not on CurseForge

Some mods aren't on CurseForge. They should be downloaded manually and specified with `exclude`:

```
exclude ^OptiFine.*\.jar$
exclude ^Computronics.*\.jar$
```

These support regexes for updating mods manually without having to edit the mods file (you only really need to remember `.*` and `^$`), but because of that you should escape dots with a backslash (`\.`).

As stated above, **mods that aren't specified or `exclude`d will be deleted**!

## Usage

Now that the preparations are complete, simply run `cf-updater` in a terminal in the mods folder or double-click the executable. For more usage options, run `cf-updater -h`.

If you want to use an option every time without having to open a terminal on Windows, you should create a `.cmd` file. On Linux, you should create a shell script.

## Why does it delete all mods not in the mods file?

That's because it's impossible to tell if you've got an older version of the mod that should be deleted or if that is just another mod. And because Forge doesn't allow multiple versions of the same mod at the same time, the old one must be deleted. So in the end, this is the only possible solution (again, suggestions are welcome!).
