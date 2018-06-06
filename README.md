# Haute couture [![GoDoc](https://godoc.org/github.com/baltimore-sun-data/haute-couture?status.svg)](https://godoc.org/github.com/baltimore-sun-data/haute-couture) [![Go Report Card](https://goreportcard.com/badge/github.com/baltimore-sun-data/haute-couture)](https://goreportcard.com/report/github.com/baltimore-sun-data/haute-couture)

Haute couture looks through your CSS and static HTML to ensure that there are no out-of-date styles. It checks that all class and ID names in the specified CSS file exist in some HTML file. Styles that don't match are output to a file for review by the developer.

## Installation

First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```bash
GOBIN="$(pwd)" GOPATH="$(mktemp -d)" go get github.com/baltimore-sun-data/haute-couture
```

## Usage

```bash
$ haute-couture -h
Haute couture looks through your CSS and static HTML to ensure that there are
no out-of-date styles.

Usage: haute-couture [options]

Options:
  -css string
        CSS file to match against
  -exclude value
        regexp for sub-directories to exclude (default "^\.")
  -html-dir string
        directory to search for HTML files (default "public")
  -include value
        regexp for HTML files to process (default "\.html?$")
  -output string
        file to save any found extra CSS identifiers in (default "extra-css.txt")
```

## Endorsements

> Well, the name is dynamite.

— [Gregabit](https://www.reddit.com/r/golang/comments/8os2wp/haute_couture_looks_through_your_css_and_static/e068dxl/)
