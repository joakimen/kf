# kf

> [!IMPORTANT]
> Superseded by https://github.com/joakimen/scriv

[![ci](https://github.com/joakimen/kf/actions/workflows/ci.yml/badge.svg)](https://github.com/joakimen/kf/actions/workflows/ci.yml)

Known files — manage the files you visit regularly.

`kf` keeps a list of file paths you return to often and provides commands to add,
list, pick, forget, and prune them.

## Installation

```sh
mise use -g github:joakimen/kf
```

Or build from source:

```sh
cargo install --git https://github.com/joakimen/kf
```

## Usage

```sh
kf add ~/.zshrc          # add a file
kf list                  # list known files
kf list --status         # list with ✓/✗ existence indicators
kf list --missing        # only files that no longer exist
kf pick                  # fuzzy-pick a file and print its path
kf forget ~/.zshrc       # remove a file (omit the path to pick interactively)
kf prune                 # drop entries whose files no longer exist
kf edit                  # open the config in $EDITOR
kf config                # print the config file path
```

Paths under `$HOME` are stored as `~/…`; relative paths are resolved against the
working directory.

## Configuration

Entries live one per line in the config file, resolved in this order:

1. `--config <path>`
2. `$XDG_CONFIG_HOME/kf/config`
3. `~/.config/kf/config`

Colored status output is suppressed when stdout is not a terminal or when
`NO_COLOR` is set.
