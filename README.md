# CertGuard

[![GoReportCard example](https://goreportcard.com/badge/github.com/pimg/certguard)](https://goreportcard.com/report/github.com/pimg/certguard) ![CI tests](https://github.com/pimg/certguard/actions/workflows/build.yml/badge.svg)

A Terminal User Interface (TUI) for inspecting Certificate Revocation Lists (CRL's)

With CertGuard it is currently possible to:
- download new CRL files to the local cache directory
- browse locally downloaded CRL files
- inspect entries in a CRL file

![demo](demo.gif)

## File locations
CertGuard uses two file locations:
- `~/.cache/certguard` for the file cache where CRL files are stored 
- `~/.local/share/certguard` for the `debug.log` file

## States
CertGuard TUI is built with [BubbleTea](https://github.com/charmbracelet/bubbletea/tree/master) using the [Elm architecture](https://guide.elm-lang.org/architecture/).
Different screens are built using different states. Below is a statemachine depicting the state model of CertGuard:

![states](states.svg)

## Development
A MAKE file has been included for convenience:
- `make run` builds and run the `certguard` application in `debug` mode
- `make test` runs all unit tests
- `make lint` runs the linter 

Since a TUI application cannot log to `stdout` a `debug.log` file is used for debug logging. It is located at: `~/.local/share/certguard/debug.log`
