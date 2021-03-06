# Important!
**Considerably bad UX has befallen this project!** If you're playing Minecraft, use `mcmodmgr` _(under development, soon to be on GitHub)_ instead.

### History
I started this project with the goal of making an easy-to-use mod updater. But after switching to the Twitch API in [#1](https://github.com/DimiDimit/cf-updater/issues/1), which didn't have searching by slugs, I had to make a compromise and make users specify IDs, because my implementation was inflexible. As other features got added, I just threw them together *somehow.* Soon, you had to copy IDs, remember arguments and leave tons of comments—it became a mess! Now, with a clear vision, strict rules and going in the right direction where I had previously taken a wrong turn (e.g. choosing Rust instead of Go (not hating on Go, I just found I like Rust more), using TOML rather than making my own parser, etc.) I'm planning something much bigger, better, and focused (on Minecraft).

### The future of this project
OK, so what about `cf-updater`? Well, unlike `mcmodmgr`, it's *supposed to* work for any CurseForge game, not just Minecraft, and it still works (I think), so I'll continue maintaining it, fixing bugs and adding minor features, at least for now. I'm not planning on doing any serious updates though, so if you're willing to do anything big, instead of submitting a pull request you might as well fork it, tell me in an issue/discussion and I'll link to it here.

# CurseForge Updater

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/mod/github.com/DimiDimit/cf-updater/v3)
![Build](https://github.com/DimiDimit/cf-updater/workflows/Build/badge.svg)
![Test and Lint](https://github.com/DimiDimit/cf-updater/workflows/Test%20and%20Lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/DimiDimit/cf-updater)](https://goreportcard.com/report/github.com/DimiDimit/cf-updater)

A tool to update [CurseForge](https://www.curseforge.com) mods, written in [Go](https://golang.org). Currently only tested with [Minecraft](https://www.curseforge.com/minecraft/mc-mods).

## Compatibility

### Compatibility with v2

Version 3's mods file **is** compatible with version 2! But note that the API is **not**.

### Compatibility with v1

Version 2 is **not** compatible with version 1! You must **refactor your mods file** or you'll get syntax errors!

## Installation and upgrading

### Installing on Windows

To install it on Windows, download it from the [Releases](https://www.github.com/DimiDimit/cf-updater/releases) page and save it into your mods folder (e.g. for Minecraft it's `.minecraft/mods` or `.minecraft/mods/<version>`).

### Installing with `go get`

First, [install Go](https://golang.org/doc/install).
Then run:

```sh
go get -u github.com/DimiDimit/cf-updater
```

### Upgrading

The program has **no** built-in update checker or updater, so you should **watch the repository** in “Releases only” mode to get notified of updates.

The upgrade process is the same as the installation.

## Setup

**Backup your mods before doing anything else!**

Unfortunately, some setup is first required. You must put the IDs of all of your mods in a file called `mods.txt`. This can be tedious if you've got a lot of mods, but I don't see a way around it. If you have an idea, please submit an [issue](https://www.github.com/DimiDimit/cf-updater/issues) or a [pull request](https://www.github.com/DimiDimit/cf-updater/pulls)!

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
