# Humanize

Convert big numbers to human-readable.

## Background

As a data engineer, I encounter a lot of big numbers in the course of my work. I got tired of
squinting and counting digits to tell if something was in the billions or trillions, so I wrote
this to replace the numbers inline.

## Usage

```
Usage of humanize:
  -bin
    	use base-2 divisors instead of base-10
  -min float
    	minimum absolute value to humanize (default 1000)
```

```
$ echo "Woah, check out this big number: 123456789012345" | humanize
Woah, check out this big number: 123T
$ echo "Woah, check out this big number: 123456789012345" | humanize -bin
Woah, check out this big number: 112Ti
```
