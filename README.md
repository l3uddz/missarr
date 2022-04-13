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

missarr offers [pre-compiled binaries](https://github.com/l3uddz/missarr/releases/latest) for Linux, macOS and Windows for each official release.

Example install (installing to /opt/missarr on linux amd64):
````
cd /opt
mkdir missarr && cd missarr
curl -fLvo missarr https://github.com/jolbol1/missarr/releases/download/v1.2.0/missarr_v1.2.0_linux_amd64
chmod +x missarr
````

Alternatively, you can build the Missarr binary yourself.
To build missarr on your system, make sure:

1. Your machine runs Linux, macOS or Windows
2. You have [Go](https://golang.org/doc/install) installed (1.16 or later preferred)
3. Clone this repository and cd into it from the terminal
4. Run `make build` from the terminal

You should now have a binary with the name `missarr` in the appropriate dist subdirectory of the project.

If you need to debug certain Missarr behaviour, either add the `-v` flag for debug mode or the `-vv` flag for trace mode to get even more details about internal behaviour.

### Full config file

```yaml
sonarr:
  url: https://sonarr.your-domain.com
  api_key: your_api_key
radarr:
  url: https://radarr.your-domain.com
  api_key: your_api_key
```

If you are experiencing timeouts while retrieving data from your PVR, you can add the `timeout` config option, which currently defaults to `90` (seconds).

You can place this config file in the same folder as the missarr binary as `config.yml`

Alternatively, you can place this config file in `~/.config/missarr/config.yml`

### Example commands

Update to latest version: `missarr --update`

Search for 10 seasons: `missarr sonarr --limit 10`

Search for 10 seasons (without updating seasons cache) `missarr sonarr --limit 10 --skip-refresh`

Search for 10 movies: `missarr radarr --limit 10`

Search for 10 movies (without updating movies cache) `missarr radarr --limit 10 --skip-refresh`

Search for 10 movies with cutoff unmet: `missarr radarr --limit 10 --cutoff`


## Donate

If you find this project helpful, feel free to make a small donation:

- [Revolut](https://revolut.me/l3uddz): Credit Cards, Apple Pay, Google Pay

- [Paypal: l3uddz@gmail.com](https://www.paypal.me/l3uddz)

- [GitHub Sponsor](https://github.com/sponsors/l3uddz): GitHub matches contributions for first 12 months.

- BTC: 3CiHME1HZQsNNcDL6BArG7PbZLa8zUUgjL
