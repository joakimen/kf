# kf

known files

Manage files that you visit somewhat regularly.

## Description

Manages the files you visit regularly, providing commands for adding, removing, and listing candidates in a configuration file.

## Configuration

The configuration file `~/.config/kf/config` should contain a list of files you want to manage. Each line should contain a file path, such as `/Users/jason/.bashrc`, or `~/.zshrc`.

## Installation

```sh
go install github.com/joakimen/kf@latest
```

## Usage

### List

List all known files

```sh
kf list
```

### Add

Add a file to the list of known files

```bash
kf add ~/.zshrc
```
