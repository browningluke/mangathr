# 📦 MangathrV2
[![Go Report Card](https://goreportcard.com/badge/github.com/browningluke/mangathrV2)](https://goreportcard.com/report/github.com/browningluke/mangathrV2)
![docker build](https://github.com/browningluke/mangathrV2/actions/workflows/docker-publish.yml/badge.svg)

**mangathrV2** is a command-line program to download Manga chapters from numerous online platforms (See [Sources](#sources)). It bundles each chapter with metadata in ComicInfo format, or a number of others (see [Metadata Agents](#metadata-agents)). It supports monitoring a source for new chapters of registered Manga.

It has an older version ([mangathr](https://github.com/browningluke/mangathr)), written in Typescript.

```
$ mangathrv2 <COMMAND> -s SOURCE QUERY
```

- [Installation](#installation)
- [Usage](#usage)
- [Options](#options)
- [Configuration](#configuration)
- [Sources](#sources)
- [Databases](#databases)
- [Metadata Agents](#metadata-agents)

## Installation

### Go install

If Go is installed, install it with:
```
$ go install github.com/browningluke/mangathrV2/cmd/mangathrV2@latest
```

### Pre-build binary

The pre-built binary files for most common UNIX OS/ARCH can be found in the latest [release](https://github.com/browningluke/mangathrV2/releases). No external dependencies should be required.

### Docker (via Docker hub / ghcr.io)

An official docker image can be found on [Dockerhub](https://hub.docker.com/r/browningluke/mangathrv2) and [ghcr.io](https://github.com/browningluke/mangathrV2/pkgs/container/mangathrv2).

Get started with [here](docker/production/README.md) with an example [docker-compose.yml](docker/production/docker-compose.yml) file.

### Docker (via source)

The docker image can be built from source by following the instructions [here](docker/build/README.md).

### Building manually

Clone the repo with:
```
$ git clone https://github.com/browningluke/mangathrV2.git
$ cd mangathrV2
```

Install all Go module dependencies with:
```
$ go mod download && go mod verify
```

Build executable with:
```
$ go build -v -o path/to/place/bin ./cmd/mangathrV2
```

Make binary executable and move into path with:
```
$ chmod +x path/to/place/bin
$ cp path/to/place/bin /usr/local/bin/mangathrv2
```

## Usage

```
Usage:  mangathrv2 [OPTIONS] COMMAND

📦 A CLI utility for downloading Manga & metadata.

Commands:
  completion  Generate the autocompletion script for the specified shell
  download    Download chapters from source
  register    Register chapters to database
  update      Check for new chapters

Options:
      --config string      config file (default is $XDG_CONFIG_HOME/mangathrv2/config)
  -h, --help               help for mangathrv2
  -l, --log-level string   Set the logging level ("debug"|"info"|"warn"|"error"|"off") (default "off")

Use "mangathrv2 [command] --help" for more information about a command.
```

## Options

```
    --config string          Path to config file (default is $XDG_CONFIG_HOME/mangathrv2/config)
-l, --log-level string       Set the logging level ("debug"|"info"|"warn"|"error"|"off") (default "off")
-h, --help                   Print this text
    --version                Print program version and exit
```

### Download / Register

```
    --dry-run            Disables downloads & writes to the database
-s, --source string      Source to search query on (REQUIRED)
```

### Update

```
    --dry-run            Disables downloads & writes to the database
```

## Configuration

The default configuration path is `$XDG_CONFIG_HOME/mangathrv2/config`, which on most machines should be `~/.config/mangathrv2/config`.

```yaml
logLevel: "off"                  # One of: ("debug"|"info"|"warn"|"error"|"off") 
database:
  driver: "sqlite"               # One of the supported database drivers
  uri: "examples/db.sqlite"      # Relative/absolute path to database (or URI for non-file dbs)
downloader:
  dryRun: false                  # Disables downloads & writes to the database
  simultaneousPages: 2           # Number of pages to download at once
  pageRetries: 5                 # Number of time to retry a page before failing
  delay:
    page: "100ms"                # Time to sleep between each (batch of) page(s)
    chapter: "200ms"             # Time to sleep between each chapter when downloading
    updateChapter: "2s"          # Time to sleep between each chapter when checking for updates
  output:
    path: './manga'              # Location to place downloaded chapters
    updatePath: './new-manga'    # Location of new chapters (defaults to value of path)
    zip: true                    # Whether to zip chapters into CBZ archives
    filenameTemplate: "..."      # See following section on templates
  metadata:
    agent: "comicinfo"           # One of the supported metadata formats 
    location: "internal"         # Whether to place metadata file inside or next to archive
sources:
  ...                            # See relevant source for configuration options
```


### Templating

mangathrV2 uses a custom (simplified) templating structure. It should be powerful enough to cover most desired changes of positioning and addition/removal of variables from chapter names.

The following is the default template used, and the filename it produces:
```
"{num:3} - Chapter {num}{title: - <.>}{groups: [<.>]}"

| 001 - Chapter 1 - This is an example chapter [Group 1, Group 2].cbz
| ...
```

The format is: `{var: - text to substitute <.> into - } `

The variable specified following the leading bracket `{` replaces the `<.>`.

Available variables are (verbatim):
```
- num      # Number               (e.g. 2.5)
- lang     # ISO 2 Language code  (e.g fr)
- title    # Title (no 'chapter') (e.g. This is an example chapter)
- groups   # Scanlation groups    (e.g. Group 1, Group 2)
```

- `num` can have x number of zeros padded to it by appending `:x` to the variable name
- `lang`, `title` and `groups` may sometimes be empty (depending on the source), if this is the case, nothing within `{...}` brackets is substituted.





## Sources

| Source   | URL                   | Scraper | Account Sync |
|----------|-----------------------|---------|--------------|
| Mangadex | https://mangadex.org/ |    ✓    |     WIP      |


### Mangadex

#### Configuration

```yaml
mangadex:
  filenameTemplate: "..."                              # See config template section (overrides global template)
  ratingFilter: ["safe", "suggestive", "erotica"]      # Show only these ratings when searching
  languageFilter: ["en", "fr"]                         # Include chapters with these languages
  dataSaver: false                                     # Use Mangadex's 'data saver' page size
```

## Databases

The following are the currently supported database drivers. More may be added as interest permits.

- sqlite3

## Metadata Agents

- ComicInfo
- JSON (WIP)

### ComicInfo

This agent will write a `ComicInfo.xml` file to the location selected in the configuration. It has been tested and confirmed to work with [Komga](https://github.com/gotson/komga) v0.157.0-master. More info about the ComicInfo format can be found [here](https://github.com/anansi-project/comicinfo).

The following is an example of a ComicInfo.xml file, and the metadata ingested:

```xml
<?xml version="1.0" encoding="UTF-8"?>
 <ComicInfo>
   <Title>Chapter 1 - This is an example chapter</Title>
   <Number>1</Number>
   <Web>https://mangadex.org/chapter/{CHAPTER-ID}</Web>
   <Year>2022</Year>
   <Month>07</Month>
   <Day>22</Day>
   <Editor>Group 1, Group 2</Editor>
   <PageCount>42</PageCount>
 </ComicInfo>
```
