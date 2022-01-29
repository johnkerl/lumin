// ================================================================
// lumin highlights matches to a specified pattern in the specified files.
// This is like grep with --color, except it shows all lines, not just matching
// lines.
// ================================================================

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/johnkerl/lumin/internal/pkg/argf"
	"github.com/johnkerl/lumin/pkg/colors"
)

const ENV_COLOR_NAME = "LUMIN_MATCH_COLOR"

// Default escape sequence to start colorization
var highlightStartString = colors.MakeANSIEscapesFromNameUnconditionally("196")

// ----------------------------------------------------------------
func usage(ostream *os.File, exitCode int) {
	fmt.Fprintf(ostream,
		`Usage: %s [options] {pattern} [zero or more filenames]
Highlights matches to {pattern} in the specified files.
If zero filenames are specified, standard input is read.
This is like grep with --color, except it shows all lines, not just
matching lines.

Options:
-w                     Restrict matches to word boundaries.
-i                     Allow for case-insensitive matches.
-c|--color {name}      Use {name} to highlight matches -- see -l/-n for choices.
                       Example names: red, yellow, green, orchid, 9, 11, 2, 170.
                       You can also set the %s environment variable if you like.
-l|--list-color-codes  Show available color codes 0..255.
-n|--list-color-names  Show available color names (aliases for the 0..255 codes).
--                     Signify end of options, so next argument is the pattern.
                       E.g. to search for "-x" in file foo.txt, use "lumin -- -x foo.txt".
-h|--help              Print this messsage.
`,
		os.Args[0], ENV_COLOR_NAME,
	)
	os.Exit(exitCode)
}

func main() {
	// Set defaults for options
	matchOnWordBoundary := false
	caseInsensitive := false
	envColorName := os.Getenv(ENV_COLOR_NAME)
	if envColorName != "" {
		ok := setColor(envColorName)
		if !ok {
			fmt.Fprintf(os.Stderr, "%s: color \"%s\" not found.\n", os.Args[0], envColorName)
			fmt.Fprintf(os.Stderr, "See %s -h for help.\n", os.Args[0])
			os.Exit(1)
		}
	}

	os.Args = getoptify(os.Args) // lumin -iw -> lumin -i -w

	// Parse command-line options
	argi := 1
	argc := len(os.Args)

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := os.Args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++
		if opt == "--" {
			break // Let people search for things starting with a dash via "lumin -- -x filename.txt"
		}

		if opt == "-h" || opt == "--help" {
			usage(os.Stdout, 0)

		} else if opt == "-w" {
			matchOnWordBoundary = true

		} else if opt == "-i" {
			caseInsensitive = true

		} else if opt == "-l" || opt == "--list-color-codes" {
			colors.ListColorCodes()
			os.Exit(0)

		} else if opt == "-n" || opt == "--list-color-names" {
			colors.ListColorNames()
			os.Exit(0)

		} else if opt == "-c" || opt == "--color" {
			if argi >= argc {
				fmt.Fprintf(os.Stderr, "%s: option %s requires an argument.\n", os.Args[0], opt)
				fmt.Fprintf(os.Stderr, "See %s -h for help.\n", os.Args[0])
				os.Exit(1)
			}
			colorName := os.Args[argi]
			ok := setColor(colorName)
			if !ok {
				fmt.Fprintf(os.Stderr, "%s: color \"%s\" not found.\n", os.Args[0], colorName)
				fmt.Fprintf(os.Stderr, "See %s -h for help.\n", os.Args[0])
				os.Exit(1)
			}
			argi++

		} else {
			fmt.Fprintf(os.Stderr, "%s: Unrecognized option \"%s\".\n", os.Args[0], opt)
			fmt.Fprintf(os.Stderr, "See %s -h for help.\n", os.Args[0])
			os.Exit(1)
		}
	}

	if argi >= argc {
		fmt.Fprintf(os.Stderr, "%s: need search pattern as argument.\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "See %s -h for help.\n", os.Args[0])
		os.Exit(1)
	}
	pattern := os.Args[argi]
	filenames := os.Args[argi+1:]

	if matchOnWordBoundary {
		pattern = "\\b" + pattern + "\\b"
	}

	if caseInsensitive {
		pattern = "(?i)" + pattern
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
	}

	istream, err := argf.Open(filenames)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = luminStream(regex, istream)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func setColor(name string) bool {
    escape, ok := colors.MakeANSIEscapesFromName(name)
    if ok {
        highlightStartString = escape
    }
    return ok
}

func colorize(input string) string {
	return highlightStartString + input + colors.DefaultColorString
}

// ----------------------------------------------------------------
// getoptify expands "-xyz" into "-x -y -z" while leaving "--xyz" intact. This
// is a keystroke-saver for the user.
//
// Secondly, we split "--foo=bar" into "--foo" and "bar".
func getoptify(inargs []string) []string {
	expandRegex := regexp.MustCompile("^-[a-zA-Z0-9]+$")
	splitRegex := regexp.MustCompile("^--[^=]+=.+$")
	outargs := make([]string, 0)
	for _, inarg := range inargs {
		if expandRegex.MatchString(inarg) {
			for _, c := range inarg[1:] {
				outargs = append(outargs, "-"+string(c))
			}
		} else if splitRegex.MatchString(inarg) {
			pair := strings.SplitN(inarg, "=", 2)
			outargs = append(outargs, pair[0])
			outargs = append(outargs, pair[1])
		} else {
			outargs = append(outargs, inarg)
		}
	}
	return outargs
}

// ----------------------------------------------------------------
func luminStream(regex *regexp.Regexp, istream io.Reader) error {
	scanner := bufio.NewScanner(istream)

	for scanner.Scan() {
		line := scanner.Text()
		// This is how to do a chomp:
		line = strings.TrimRight(line, "\n")
		fmt.Println(luminLine(regex, line))
	}

	return nil
}

// ----------------------------------------------------------------
func luminLine(regex *regexp.Regexp, input string) string {

	matrix := regex.FindAllStringIndex(input, -1)
	// fmt.Printf("%+v\n", matrix)
	if matrix == nil || len(matrix) == 0 {
		return input
	}

	// The key is the Go library's regex.FindAllStringIndex.  It gives us start
	// (inclusive) and end (exclusive) indices for matches.
	//
	// Example: for pattern "foo" and input "abc foo def foo ghi" we'll have
	// matrix [[4 7] [12 15]] which indicates matches from positions 4-6 and
	// 12-14.  We simply need to print out:
	// *  0-3  "abc "  with default color
	// *  4-6  "foo"   with highlight color
	// *  7-11 " def " with default color
	// * 12-14 "foo"   with highlight color
	// * 15-18 " ghi"  with default color.
	//
	// Example: with pattern "f.*o" and input "abc foo def foo ghi" we'll have
	// matrix [[4 15]] so "foo def foo" will be highlighted.

	var buffer bytes.Buffer // Faster since os.Stdout is unbuffered
	nonMatchStartIndex := 0

	for _, startEnd := range matrix {
		buffer.WriteString(input[nonMatchStartIndex:startEnd[0]])
		buffer.WriteString(colorize(input[startEnd[0]:startEnd[1]]))
		nonMatchStartIndex = startEnd[1]
	}

	buffer.WriteString(input[nonMatchStartIndex:])

	return buffer.String()
}
