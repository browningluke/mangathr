# 📦 mangathr
[![Go Report Card](https://goreportcard.com/badge/github.com/browningluke/mangathr)](https://goreportcard.com/report/github.com/browningluke/mangathr)
![docker build](https://github.com/browningluke/mangathr/actions/workflows/docker-publish.yml/badge.svg)

**mangathr** is a command-line program to download Manga chapters from numerous online platforms (See [Sources](#sources)). It bundles each chapter with metadata in ComicInfo format, or a number of others (see [Metadata Agents](#metadata-agents)). It supports monitoring a source for new chapters of registered Manga.

It has an older version ([mangathr-legacy](https://github.com/browningluke/mangathr-legacy)), written in Typescript.

```
$ mangathr <COMMAND> -s SOURCE QUERY
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
$ go install github.com/browningluke/mangathr/v2/cmd/mangathr@latest
```

### Pre-build binary

The pre-built binary files for most common UNIX OS/ARCH can be found in the latest [release](https://github.com/browningluke/mangathr/releases). No external dependencies should be required.

### Docker (via Docker hub / ghcr.io)

An official docker image can be found on [Dockerhub](https://hub.docker.com/r/browningluke/mangathr) and [ghcr.io](https://github.com/browningluke/mangathr/pkgs/container/mangathr).

Get started with [here](docker/production/README.md) with an example [docker-compose.yml](docker/production/docker-compose.yml) file.

### Docker (via source)

The docker image can be built from source by following the instructions [here](docker/build/README.md).

### Building manually

Clone the repo with:
```
$ git clone https://github.com/browningluke/mangathr.git
$ cd mangathr
```

Install all Go module dependencies with:
```
$ go mod download && go mod verify
```

Build executable with:
```
$ go build -v -o path/to/place/bin ./cmd/mangathr
```

Make binary executable and move into path with:
```
$ chmod +x path/to/place/bin
$ cp path/to/place/bin /usr/local/bin/mangathr
```

## Usage

```
Usage:  mangathr [OPTIONS] COMMAND

📦 A CLI utility for downloading Manga & metadata.

Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Manage the config file
  download    Download chapters from source
  manage      Manage series registered in database
  register    Register chapters to database
  update      Check for new chapters
  version     Print the version number of mangathr

Options:
      --config string      config file (default is $XDG_CONFIG_HOME/mangathr/config)
  -h, --help               help for mangathr
  -l, --log-level string   Set the logging level ("debug"|"info"|"warn"|"error"|"off") (default "off")
      --override strings   Override config values (e.g. database.driver=postgres,database.postgres.user=example)

Use "mangathr [command] --help" for more information about a command.
```

## Options

### Download / Register

```
    --dry-run            Disables downloads & writes to the database
-s, --source string      Source to search query on (REQUIRED)
```

### Update

```
    --dry-run            Disables downloads & writes to the database
```

### Configuration Overrides

Configuration options can be specified in the command line to temporarily override the config file's values (either set by the user, or the default values). This may be useful when switching databases, or running in unusual environments, such as Kubernetes.

This option uses the same format (`key.subkey=value`) as [Helm's](https://helm.sh/) `--set`. More details about how the values get converted to YAML can be found on [Helm's website](https://helm.sh/docs/intro/using_helm/#the-format-and-limitations-of---set).

The following command will connect to PostgreSQL with the provided details, regardless of the values in the config file.
```shell
    mangathr \
      --override database.driver=postgres \
      --override database.postgres.host=127.0.0.1 \
      --override database.postgres.user=postgres \
      --override database.postgres.password=mangathr \
      --override database.postgres.dbName=mangathr \
      update
```

## Configuration

The default configuration path is `$XDG_CONFIG_HOME/mangathr/config`, which on most machines should be `~/.config/mangathr/config`.

```yaml
logLevel: "off"                  # One of: ("debug"|"info"|"warn"|"error"|"off") 
database:
  driver: "sqlite"               # One of the supported database drivers (sqlite, postgres, etc.)
  sqlite:
    ...                          # See database section for configuration options
  postgres:
    ...                          # See database section for configuration options
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

mangathr uses a custom (simplified) templating structure. It should be powerful enough to cover most desired changes of positioning and addition/removal of variables from chapter names.

The following is the default template used, and the filename it produces:
```
"{num:3} - {source:<.>} - {manga:<.>} - Chapter {num}{title: - <.>}{groups: [<.>]}"

| 001 - Mangadex - Berserk - Chapter 1 - This is an example chapter [Group 1, Group 2].cbz
| ...
```

The format is: `{var: - text to substitute <.> into - } `

The variable specified following the leading bracket `{` replaces the `<.>`.

Available variables are (verbatim):
```
- source   # Source name          (e.g. Mangadex, Cubari)
- manga    # Manga title          (e.g. Berserk)
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
| Mangadex | https://mangadex.org/ |    ✓    | WIP          |
| Cubari   | https://cubari.moe/   |    ✓    | ✗            |


### Mangadex

#### Configuration

```yaml
mangadex:
  filenameTemplate: "..."           # See config template section (overrides global template)
  ratingFilter:                     # Show only these ratings when searching
    - "safe"
    - "suggestive"
    - "erotica"
  languageFilter: ["en", "fr"]      # Include chapters with these languages
  dataSaver: false                  # Use Mangadex's 'data saver' page size
  groups:
    include:                        # Only download chapters from <Example A> and <Example B>
      - "Example A"
      - "Example B"
    exclude:                        # Do not download chapters from <Example C>
      - "Example C"
```

### Cubari

#### Configuration

```yaml
cubari:
  filenameTemplate: "..."      # See config template section (overrides global template)
  groups:
    include:                   # Only download chapters from <Example A> and <Example B>
      - "Example A"
      - "Example B"
    exclude:                   # Do not download chapters from <Example C>
      - "Example C"
```


## Databases

| Database   | Supported |
|------------|:---------:|
| SQLite3    |     ✓     |
| PostgreSQL |     ✓     |
| MySQL      |     ✗     |
| MongoDB    |     ✗     |

### SQLite3

```yaml
sqlite:
  path: ~/mangathr/db.sqlite    # default: ~/.config/mangathr/db.sqlite (/config/db.sqlite if in container)
```

### PostgreSQL

```yaml
postgres:
  host: 127.0.0.1             # default: 127.0.0.1
  port: 5432                  # default: 5432
  user: mangathr              # default: mangathr
  password: PASSWORD          # required
  dbName: mangathr            # default: mangathr
  sslMode: disable            # default: disable
  
  # (advanced usage) Appends extra options to connection string in format `key=value`
  opts: ''                    # default: ''
 ```

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
