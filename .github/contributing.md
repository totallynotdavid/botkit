# Contributing

Read [docs/architecture.md](../docs/architecture.md) for the package map and the
boundaries. [AGENTS.md](../AGENTS.md) holds the rules every change follows.

## Set up

Follow [Install](../readme.md#install) in the README. `mise install` installs
the versions pinned in [`mise.toml`](../mise.toml), including Zig and Node. Zig
is the C compiler for `go test -race`, which needs cgo, so a clean host needs no
`gcc`; Node runs oxfmt, the Markdown formatter. The latex tests also need
`prlimit` from util-linux on `PATH`.

## Tasks

| Task                         | Does                                                                            |
| ---------------------------- | ------------------------------------------------------------------------------- |
| `mise run install`           | Downloads the Go modules.                                                       |
| `mise run lint`              | Checks Go formatting and lint rules, and the Markdown.                          |
| `mise run fix`               | Formats the Go code and the documents, applies lint fixes.                      |
| `mise run test`              | Runs `go test -race ./...` with Zig as the C compiler.                          |
| `mise run build`             | Builds `bin/botkit-flow` and `bin/botkit-latex`.                                |
| `mise run check`             | Runs `lint`, `test`, `build`, and `check:env-example`. Changes no tracked file. |
| `mise run install:bin`       | Builds and installs both bots into `$PREFIX/bin`.                               |
| `mise run install:units`     | Installs the systemd units into `$UNIT_DIR`.                                    |
| `mise run env:example`       | Rewrites `.env.example` from the settings the bots read.                        |
| `mise run check:env-example` | Checks `.env.example` against the settings the bots read.                       |
| `mise run dev:flow`          | Runs the flow bot from source, loading `.env`.                                  |
| `mise run dev:latex`         | Runs the latex bot from source, loading `.env`.                                 |

`install:bin` installs into `$PREFIX/bin` (default `/usr/local/bin`) and
`install:units` into `$UNIT_DIR` (default `/etc/systemd/system`).

Run `mise run check` before you open a pull request. GitHub Actions runs
`mise run install` and `mise run check`.

## Tests

Tests run the real apps over SQLite and the in-memory transport in
`internal/whatsapp/fake`, and never connect to WhatsApp. The latex tests run the
real `typst`.

## Settings

A bot's settings come from the code that reads them. After you add or change
one, give it an explanation in
[`internal/envfile`](../internal/envfile/envfile.go), run
`mise run env:example`, and update
[docs/configuration.md](../docs/configuration.md). `mise run check:env-example`
fails while `.env.example` is out of date.

## Documentation

oxfmt formats Markdown to 80 columns, as [`.oxfmtrc.json`](../.oxfmtrc.json)
sets. `mise run fix` rewrites the files; `mise run lint` reports the ones that
differ.
