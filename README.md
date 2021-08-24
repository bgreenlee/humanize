# Humanize

Convert big numbers to human-readable.

## Background

As a data engineer, I encounter a lot of big numbers in the course of my work. I got tired of
squinting and counting digits to tell if something was in the billions or trillions, so I wrote
this to replace the numbers inline.

## Install

```
go get github.com/bgreenlee/humanize
```

If you don't have Go installed, [do that](https://golang.org/doc/install) first.

The `humanize` binary will be installed in `$GOPATH/bin`, which is likely `~/go/bin`, so make sure that's in your `PATH`.

## Usage

```
Usage of humanize:
  -binary
        use base-2 divisors instead of base-10
  -min float
        minimum absolute value to humanize (default 1000)
  -preserve
        preserve formatting when replacing
  -version
        print version and exit
```

```
$ echo "Woah, check out this big number: 123456789012345" | humanize
Woah, check out this big number: 123T
$ echo "Woah, check out this big number: 123456789012345" | humanize -binary
Woah, check out this big number: 112Ti
```

## Tests

![Go](https://github.com/bgreenlee/humanize/workflows/Go/badge.svg)
