# kf

[![ci](https://github.com/joakimen/kf/actions/workflows/ci.yml/badge.svg)](https://github.com/joakimen/kf/actions/workflows/ci.yml) [![GoDoc](https://godoc.org/github.com/joakimen/kf?status.svg)](https://godoc.org/github.com/joakimen/kf) [![Go Report Card](https://goreportcard.com/badge/github.com/joakimen/kf)](https://goreportcard.com/report/github.com/joakimen/kf)

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

List all known files.

```sh
kf list
```

### Add

Add a file to the list of known files.

```bash
kf add ~/.zshrc
```

### Remove

Remove a file from the list of known files.

```bash
kf forget ~/.zshrc
```

### Show config file

Show the path to the configuration file.

```bash
kf config
```
