package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

var version = "v1.2.0"

type options struct {
	isBinary           bool
	minValue           float64
	preserveFormatting bool
}

// numToHuman returns a replacement function used by rexexp.ReplaceAllStringFunc
// to replace numbers in a string with more human-readable values
func numToHuman(opt options) func(string) string {
	var suffixes = []string{"K", "M", "G", "T", "P", "E", "Z", "Y"}

	var base = 10.0      // float because it reduces the number of casts we have to do below
	var expdivisor = 3.0 // each successive suffix represents a Pow(base, expdivisor) increase

	if opt.isBinary {
		base = 2
		expdivisor = 10
	}

	return func(numstr string) string {
		// we parse as a float64 because it allows us to handle much larger values than uint64
		if num, err := strconv.ParseFloat(numstr, 64); err == nil {
			// don't bother if it's small enough
			if num < opt.minValue {
				return numstr
			}

			exp := math.Log(num) / math.Log(base) // same as Log<base>(num)
			suffidx := int((exp / expdivisor) - 1)
			if suffidx >= len(suffixes) {
				suffidx = len(suffixes) - 1
			} else if suffidx < 0 {
				suffidx = 0
			}
			shortnum := int(math.Round(num / math.Pow(base, float64(suffidx+1)*expdivisor)))
			if shortnum < 0 {
				// oops, overflowed float64
				return numstr
			}

			var suffix = suffixes[suffidx]
			if opt.isBinary {
				suffix = suffix + "i"
			}

			if opt.preserveFormatting {
				// left-pad with spaces
				return fmt.Sprintf("%*d%s", len(numstr)-1, shortnum, suffix)
			} else {
				return fmt.Sprintf("%d%s", shortnum, suffix)
			}
		} else {
			// if we failed to parse the number for whatever reason, just return the original
			return numstr
		}
	}

}

// humanize takes a Reader and configuration flags and replaces any numbers it reads
// with more human-readable versions
func humanize(reader io.Reader, writer io.Writer, opt options) {
	numre := regexp.MustCompile(`\d+`)
	bufreader := bufio.NewReader(reader)
	//bufwriter := bufio.NewWriter(*writer)
	for {
		line, err := bufreader.ReadString('\n')
		fmt.Fprint(writer, numre.ReplaceAllStringFunc(line, numToHuman(opt)))
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	var opt options
	flag.BoolVar(&opt.isBinary, "binary", false, "use base-2 divisors instead of base-10")
	flag.Float64Var(&opt.minValue, "min", 1000, "minimum absolute value to humanize")
	flag.BoolVar(&opt.preserveFormatting, "preserve", false, "preserve formatting when replacing")
	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "print version and exit")

	flag.Parse()

	if versionFlag {
		fmt.Printf("humanize %s\n", version)
		os.Exit(0)
	}

	humanize(os.Stdin, os.Stdout, opt)
}
