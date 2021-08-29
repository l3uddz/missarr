[![made-with-golang](https://img.shields.io/badge/Made%20with-Golang-blue.svg?style=flat-square)](https://golang.org/)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%203-blue.svg?style=flat-square)](https://github.com/l3uddz/missarr/blob/master/LICENSE.md)
[![Discord](https://img.shields.io/discord/381077432285003776.svg?colorB=177DC1&label=Discord&style=flat-square)](https://discord.io/cloudbox)
[![Donate](https://img.shields.io/badge/Donate-gray.svg?style=flat-square)](#donate)

# Missarr

missarr sends search requests for missing episodes to Sonarr

## Table of contents

- [Installing missarr](#installing-missarr)
- [Full config file](#full-config-file)
- [Example commands](#example-commands)
- [Donate](#donate)

## Installing missarr

missarr offers [pre-compiled binaries](https://github.com/l3uddz/missarr/releases/latest) for Linux, MacOS and Windows for each official release.

Alternatively, you can build the Missarr binary yourself.
To build missarr on your system, make sure:

1. Your machine runs Linux, macOS or Windows
2. You have [Go](https://golang.org/doc/install) installed (1.16 or later preferred)
3. Clone this repository and cd into it from the terminal
4. Run `make build` from the terminal

You should now have a binary with the name `missarr` in the appropriate dist sub-directory of the project.

If you need to debug certain Missarr behaviour, either add the `-v` flag for debug mode or the `-vv` flag for trace mode to get even more details about internal behaviour.

### Full config file

```yaml
sonarr:
  url: https://sonarr.your-domain.com
  api_key: your_api_key
```

### Example commands

Update to latest version: `missarr --update`

Search for 10 seasons: `missarr sonarr --limit 10`

Search for 10 seasons (without updating seasons cache) `missarr sonarr --limit 10 --skip-refresh`

## Donate

If you find this project helpful, feel free to make a small donation:

- [Monzo](https://monzo.me/today): Credit Cards, Apple Pay, Google Pay

- [Paypal: l3uddz@gmail.com](https://www.paypal.me/l3uddz)

- [GitHub Sponsor](https://github.com/sponsors/l3uddz): GitHub matches contributions for first 12 months.

- BTC: 3CiHME1HZQsNNcDL6BArG7PbZLa8zUUgjL