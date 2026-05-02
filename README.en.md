# mdview

mdview is a terminal-native Markdown pager.

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
git clone https://github.com/YamasouA/mdview.git
cd mdview
make build
```

The binary is created at `bin/mdview`.

```sh
./bin/mdview README.md
```

Pass multiple files to open them as separate tabs.

```sh
./bin/mdview README.md README.en.md
```

## Usage

Pass the Markdown file you want to open.

```sh
go run . README.md
```

```sh
go run . README.md README.en.md
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
| `Enter` | Confirm search, open the link shown in the status bar, or jump to the selected heading in the heading list |
| `Esc` | Close search input or link list |
| `n` | Go to the next search match |
| `N` | Go to the previous search match |
| `H` | Open the heading list |
| `]h` | Go to the next heading |
| `[h` | Go to the previous heading |
| `o` | Open the link list |
| `1`-`9` | Open the corresponding link in the link list |
| `t` | Open a link in a new tab, or switch to an existing tab |
| `b` | Go back in history |
| `f` | Go forward in history |
| `gt` | Go to the next tab |
| `gT` | Go to the previous tab |
| `x` | Close the current tab |
| `q` | Quit |
| `Ctrl-c` | Quit |

## Link Navigation

mdview extracts inline Markdown links.

You can use the [link test page](docs/next.md) to try link navigation.

```md
[Next](docs/next.md)
```

Relative paths are resolved from the directory of the currently open Markdown file.

`Enter` opens the link shown in the status bar in the current tab. In the document body, the line containing the Enter target is marked with `enter>`.

`o` opens the link list. In the link list, use `j` / `k` to select a link and `Enter` to open it in the current tab. `t` opens the status-bar link or the selected link-list item in a new tab. If the target file is already open, mdview switches to that tab instead of creating a duplicate.

`gt` / `gT` switch to the next or previous tab. `x` closes the current tab and selects the tab to the right when one exists, otherwise the tab to the left. The last remaining tab cannot be closed.

External URLs are not opened in the current MVP.

## Heading Navigation

mdview extracts ATX headings (`#`, `##`, ...) and lets you move through the Markdown structure.

Press `H` to open the heading list, use `j` / `k` to select a heading, and press `Enter` to jump. In the document view, `]h` jumps to the next heading and `[h` jumps to the previous heading.

## Live Reload

When an open Markdown file changes, mdview automatically reloads it.

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
| `make build` | Build `bin/mdview` |
| `make run FILE=...` | Run mdview with the specified file |
| `make check` | Run `fmt-check`, `lint`, and `test` |

## Design

mdview is built around Bubble Tea's MVU model.

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
