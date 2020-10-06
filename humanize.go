package humanize

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

// numToHuman returns a replacement function used by rexexp.ReplaceAllStringFunc
// to replace numbers in a string with more human-readable values
func numToHuman(isBinary bool, minValue float64) func(string) string {
	var suffixes = []string{"K", "M", "G", "T", "P", "E", "Z", "Y"}

	var base = 10.0      // float because it reduces the number of casts we have to do below
	var expdivisor = 3.0 // each successive suffix represents a Pow(base, expdivisor) increase

	if isBinary {
		base = 2
		expdivisor = 10
	}

	return func(numstr string) string {
		// we parse as a float64 because it allows us to handle much larger values than uint64
		if num, err := strconv.ParseFloat(numstr, 64); err == nil {
			// don't bother if it's small enough
			if num < minValue {
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
			if isBinary {
				suffix = suffix + "i"
			}

			return fmt.Sprintf("%d%s", shortnum, suffix)
		} else {
			// if we failed to parse the number for whatever reason, just return the original
			return numstr
		}
	}

}

// humanize takes a Reader and configuration flags and replaces any numbers it reads
// with more human-readable versions
func humanize(reader io.Reader, isBinary bool, minValue float64) <-chan string {
	output := make(chan string)

	go func() {
		defer close(output)

		numre := regexp.MustCompile(`\d+`)
		bufreader := bufio.NewReader(reader)
		for {
			line, err := bufreader.ReadString('\n')
			output <- fmt.Sprint(numre.ReplaceAllStringFunc(line, numToHuman(isBinary, minValue)))
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	return output
}

func main() {
	binaryFlagPtr := flag.Bool("bin", false, "use base-2 divisors instead of base-10")
	minValuePtr := flag.Float64("min", 1000, "minimum absolute value to humanize")
	flag.Parse()

	output := humanize(os.Stdin, *binaryFlagPtr, *minValuePtr)
	for line := range output {
		fmt.Print(line)
	}
}
