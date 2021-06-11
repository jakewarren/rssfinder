# rssfinder
[![CircleCI](https://circleci.com/gh/jakewarren/rssfinder.svg?style=shield)](https://circleci.com/gh/jakewarren/rssfinder)
[![GitHub release](http://img.shields.io/github/release/jakewarren/rssfinder.svg?style=flat-square)](https://github.com/jakewarren/rssfinder/releases])
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/jakewarren/rssfinder/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakewarren/rssfinder)](https://goreportcard.com/report/github.com/jakewarren/rssfinder)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=shields)](http://makeapullrequest.com)
> Finds RSS feeds for a given website

## Install
### Option 1: Binary

Download the latest release from [https://github.com/jakewarren/rssfinder/releases/latest](https://github.com/jakewarren/rssfinder/releases/latest)

### Option 2: From source

```
go get github.com/jakewarren/rssfinder
```

### Option 3: From gobinaries.com

```
curl -sf https://gobinaries.com/jakewarren/rssfinder | sh
```

## Usage

```
❯ rssfinder -h
Usage: rssfinder [flags] <URL>

Flags:
  -f, --fuzzer    enables the fuzzer module
  -h, --help      display help
  -s, --scraper   enables the scraper module
  -v, --verbose   enable trace logging
  -V, --version   display version information

URL: https://github.com/jakewarren/rssfinder
```

## Changes

All notable changes to this project will be documented in the [changelog].

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

## License

MIT © 2021 Jake Warren

[changelog]: https://github.com/jakewarren/rssfinder/blob/master/CHANGELOG.md
