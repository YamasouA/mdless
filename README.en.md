# mdless

mdless is a terminal-native Markdown pager.

It is a small tool for reading Markdown directly in the terminal, without a browser or a server.

## Concept

- Terminal-only
- Read-only
- Vim-like controls
- Markdown link navigation
- Lightweight single binary

## Current Status

The MVP currently supports:

- Rendering Markdown files
- Scrolling with `j` / `k` and related keys
- Searching with `/` and highlighting search matches
- Extracting inline Markdown links
- Opening relative links
- Back/forward navigation with `b` / `f`
- Opening links in new tabs
- Live reload when open files change

Session persistence is planned for a future release.

## Requirements

- Go 1.26.1 or later
- make

## Installation

Clone this repository and build the binary.

```sh
git clone https://github.com/YamasouA/mdless.git
cd mdless
make build
```

The binary is created at `bin/mdless`.

```sh
./bin/mdless README.md
```

## Usage

Pass the Markdown file you want to open.

```sh
go run . README.md
```

You can also run it through the Makefile.

```sh
make run FILE=README.md
```

## Controls

| Key | Action |
| --- | --- |
| `j` | Scroll down one line |
| `k` | Scroll up one line |
| `Ctrl-d` | Scroll down half a page |
| `Ctrl-u` | Scroll up half a page |
| `gg` | Go to the top |
| `G` | Go to the bottom |
| `/` | Start search input |
| `Enter` | Confirm search, or open the link shown in the status bar |
| `Esc` | Close search input or link list |
| `n` | Go to the next search match |
| `N` | Go to the previous search match |
| `o` | Open the link list |
| `1`-`9` | Open the corresponding link in the link list |
| `t` | Open a link in a new tab |
| `b` | Go back in history |
| `f` | Go forward in history |
| `gt` | Go to the next tab |
| `gT` | Go to the previous tab |
| `x` | Close the current tab |
| `q` | Quit |
| `Ctrl-c` | Quit |

## Link Navigation

mdless extracts inline Markdown links.

You can use the [link test page](docs/next.md) to try link navigation.

```md
[Next](docs/next.md)
```

Relative paths are resolved from the directory of the currently open Markdown file.

External URLs are not opened in the current MVP.

## Live Reload

When an open Markdown file changes, mdless automatically reloads it.

If the same file is open in multiple tabs, all matching tabs are updated together. Each tab keeps its scroll position after reload.

## Development

The project uses a Makefile as its task runner.

```sh
make help
make fmt
make lint
make test
make build
make check
```

Main targets:

| Target | Description |
| --- | --- |
| `make fmt` | Format Go files with `gofmt` |
| `make fmt-check` | Check whether any Go files are unformatted |
| `make lint` | Run `go vet ./...` |
| `make test` | Run `go test ./...` |
| `make build` | Build `bin/mdless` |
| `make run FILE=...` | Run mdless with the specified file |
| `make check` | Run `fmt-check`, `lint`, and `test` |

## Design

mdless is built around Bubble Tea's MVU model.

- `cmd`: CLI
- `internal/app`: TUI Model / Update / View
- `internal/render`: Markdown rendering and link extraction
- `internal/nav`: Links and history
- `internal/watch`: Live reload file watcher
- `internal/ui`: lipgloss styles

Markdown rendering uses `glamour`, the TUI uses `bubbletea`, layout and styling use `lipgloss`, and the CLI uses `cobra`.

## Roadmap

- Session persistence
- Status bar improvements
- Link selection UI improvements
