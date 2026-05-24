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
- [Hooks](#hooks)
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

| Source    | URL                               | Scraper | Account Sync |
|-----------|-----------------------------------|---------|--------------|
| Mangadex  | https://mangadex.org/             |    ✓    | WIP          |
| Cubari    | https://cubari.moe/               |    ✓    | ✗            |
| MangaPlus | https://mangaplus.shueisha.co.jp/ |    ✓    | ✗            |


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

### MangaPlus

#### Configuration

```yaml
mangaplus:
  filenameTemplate: "..."    # See config template section (overrides global template)
  language: 0                # ENGLISH=0, SPANISH=1, FRENCH=2, INDONESIAN=3, PORTUGUESE_BR=4, RUSSIAN=5, THAI=6
  imageQuality: "super_high" # One of: (super_high|high|low)
  split: "no"                # One of: (no|yes) — whether to split spreads
```

## Hooks

Hooks execute at lifecycle events during `download` and `update` runs. They can call webhooks or run local commands. All hook fields that accept a string payload use Go's [`text/template`](https://pkg.go.dev/text/template).

### Events

| Event | Fires when |
|---|---|
| `download.chapter` | A chapter is successfully downloaded (both `download` and `update` commands) |
| `update.success` | A manga is successfully checked/downloaded during an `update` run |
| `update.error` | An error occurs while processing a manga during an `update` run |

### Common options

All hook types share these fields:

| Field | Type | Default | Description |
|---|---|---|---|
| `name` | string | — | Identifies the hook in logs (required) |
| `on` | list | — | Events that trigger this hook (required) |
| `aggregate` | bool | `false` | Collect all events and fire once at end of run instead of immediately |
| `abortOnError` | bool | `false` | Propagate hook failure to abort the current operation |
| `fireIfEmpty` | bool | `false` | Fire even when no chapters were downloaded |

### Template context

When `aggregate: false`, templates receive a `HookContext`:

```
.Manga.Title        # Manga title
.Manga.Source       # Source name (e.g. "mangadex")
.Chapter.Num        # Chapter number string
.Chapter.Title      # Chapter title
.Chapter.Path       # Filesystem path where the chapter was saved
.Chapter.Lang       # Language code
.Chapter.Groups     # Comma-separated scanlation group names
.Chapter.Count      # Number of chapters downloaded in this invocation
.Error.Message      # Error message (only set for error events; nil otherwise)
.Event              # Event name that triggered this hook
```

When `aggregate: true`, templates receive an `AggregateHookContext`:

```
.Items[]            # Slice of HookContext, one per collected event
.ChapterCount       # Total chapters across all items
.ErrorCount         # Number of items with errors
.MangaCount         # Number of distinct manga titles
.Event              # Event name of the last collected item
```

### Discord

Posts a message to a Discord webhook. Set `embed: true` for a rich embed, or `false` for a plain text message.

```yaml
hooks:
  discord:
    - name: "My Discord"
      webhookURL: "https://discord.com/api/webhooks/..."
      abortOnError: false
      fireIfEmpty: false
      on:
        - update.success
      aggregate: true           # fire once with a summary at end of update run
      embed: true
      template:
        title: "mangathr update"
        description: |
          {{ range .Items }}• **{{ .Manga.Title }}**: {{ .Chapter.Count }} chapter(s)
          {{ end }}
        color: 5814783          # decimal colour value
        footer: "mangathr"
        # message: "..."        # used instead of title/description/color/footer when embed: false
```

### Webhook

Fires a templated HTTP request. The `body` and `headers` fields are `text/template` strings; `headers` must render to a JSON object.

```yaml
hooks:
  webhook:
    - name: "My API"
      webhookURL: "https://example.com/hook"
      requestType: POST         # GET / POST / PUT / PATCH (default: POST)
      successCode: 200
      abortOnError: false
      fireIfEmpty: false
      on:
        - download.chapter
      aggregate: false
      body: |
        {"manga": "{{ .Manga.Title }}", "chapter": "{{ .Chapter.Num }}", "path": "{{ .Chapter.Path }}"}
      headers: |
        {"Authorization": "Bearer TOKEN"}
```

### Subcommand

Runs a local command. Each entry in `args` and each value in `env` is a `text/template` string.

```yaml
hooks:
  subcommand:
    - name: "Post-process"
      command: "/usr/local/bin/post-process.sh"
      abortOnError: true        # abort download if script exits non-zero
      fireIfEmpty: false
      on:
        - download.chapter
      aggregate: false
      args:
        - "{{ .Chapter.Path }}"
      env:
        MANGA_TITLE: "{{ .Manga.Title }}"
        CHAPTER_NUM: "{{ .Chapter.Num }}"
```

## Databases

| Database   | Supported |
|------------|:---------:|
| SQLite3    |     ✓     |
| PostgreSQL |     ✓     |


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
