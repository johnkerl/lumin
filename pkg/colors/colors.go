// ================================================================
// ANSI 256-color support for lumin.
//
// Note: code-share with github.com/johnkerl/miller.
// ================================================================

package colors

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Escape sequence to end colorization
const DefaultColorString = "\033[0m"

// ================================================================
// Externally visible API

// ListColorCodes shows codes in the range 0..255.
func ListColorCodes() {
	fmt.Println("Available color codes:")
	for i := 0; i <= 255; i++ {
		fmt.Printf("%s%3d%s", makeColorString(i), i, DefaultColorString)
		if i%16 < 15 {
			fmt.Print(" ")
		} else {
			fmt.Print("\n")
		}
	}
}

// ListColorNames shows names for codes in the range 0..255.
func ListColorNames() {
	fmt.Println("Available color names:")
	fmt.Printf("%s%s%s\n", MakeANSIEscapesFromNameUnconditionally("bold"), "bold", DefaultColorString)
	fmt.Printf("%s%s%s\n", MakeANSIEscapesFromNameUnconditionally("underline"), "underline", DefaultColorString)
	fmt.Printf("%s%s%s\n", MakeANSIEscapesFromNameUnconditionally("reverse"), "reverse", DefaultColorString)
	for _, pair := range namesAndCodes {
		fmt.Printf(
			"%s%-20s%s %d\n",
			makeColorString(pair.code), pair.name, DefaultColorString, pair.code,
		)
	}
}

// MakeANSIEscapesFromName takes several pieces delimited by "-", e.g.  "bold-red" or "red-bold" or
// "underline-236".
func MakeANSIEscapesFromName(name string) (string, bool) {
	namePieces := strings.Split(name, "-")
	escapePieces := make([]string, len(namePieces))
	for i := range namePieces {
		escapePiece, ok := makeANSIEscapeFromName(namePieces[i])
		if !ok {
			return "", false
		}
		escapePieces[i] = escapePiece
	}
	return strings.Join(escapePieces, ""), true
}

// MakeANSIEscapesFromNameUnconditionally is for hard-coded defaults within this source file.
func MakeANSIEscapesFromNameUnconditionally(name string) string {
	escape, ok := MakeANSIEscapesFromName(name)
	if !ok {
		fmt.Fprintf(os.Stderr, "internal coding error in colorizer initialization.\n")
	}
	return escape
}

// makeANSIEscapeFromName operates on a single piece from makeANSIEscapesFromName, e.g. the "bold",
// "red", etc.
func makeANSIEscapeFromName(name string) (string, bool) {
	if name == "plain" {
		return "", true
	} else if name == "bold" || name == "bolded" {
		return boldString, true
	} else if name == "underline" || name == "underlined" {
		return underlineString, true
	} else if name == "reverse" || name == "reversed" {
		return reversedString, true
	} else {
		code, ok := makeColorCodeFromName(name)
		if ok {
			return makeColorString(code), true
		} else {
			return "", false
		}
	}
}

// makeColorString constructs an ANSI-16-color-mode escape sequence if arg is
// in 0..15, else ANSI-256-color-mode escape sequence
func makeColorString(i int) string {
	if 0 <= i && i <= 15 {
		return make16ColorString(i)
	} else {
		return make256ColorString(i)
	}
}

// makeColorCodeFromName looks up a named code, if available: e.g. "orchid"
// maps to 170.
func makeColorCodeFromName(name string) (int, bool) {
	// Things like "170"
	code, err := strconv.Atoi(name)
	if err == nil {
		return code, true
	}

	// Things like "orchid"
	for _, pair := range namesAndCodes {
		if pair.name == name {
			return pair.code, true
		}
	}
	return -1, false
}

// make16ColorString constructs an ANSI-16-color-mode escape sequence.
func make16ColorString(i int) string {
	i &= 15
	boldBit := (i >> 3) & 1
	colorBits := i & 7
	return fmt.Sprintf("\033[%d;%dm", boldBit, 30+colorBits)
}

// make256ColorString constructs an ANSI-256-color-mode escape sequence.
func make256ColorString(i int) string {
	i &= 255
	return fmt.Sprintf("\033[1;38;5;%dm", i&255)
}

// ----------------------------------------------------------------
// Name-to-code lookup table.
//
// Source: https://jonasjacek.github.io/colors/
//
// This is intentionally an array, not a map.
//
// A map would be more efficient for lookup but maps do not preserve
// insertion-ordering in Go -- so, doing 'list colors' would print them in a
// random order. To do nice 'list colors' we'd have to put them in this kind of
// array and sort anyway. Also, lookups are done exactly once at program start
// -- so, it's more efficient to make this an array with lookup by search, than
// to populate a hash-map (computing hashcodes) for a use-only-once hash lookup.

const boldString = "\u001b[1m"
const underlineString = "\u001b[4m"
const reversedString = "\u001b[7m"

type tNameAndCode struct {
	name string
	code int
}

var namesAndCodes = []tNameAndCode{
	{"aqua", 14},
	{"aquamarine1", 122},
	// {"aquamarine1", 86}, // duplicate
	{"aquamarine3", 79},
	{"black", 0},
	{"blue", 12},
	{"blue1", 21},
	{"blue3", 19},
	// {"blue3", 20}, // duplicate
	{"blueviolet", 57},
	{"cadetblue", 72},
	// {"cadetblue", 73}, // duplicate
	{"chartreuse1", 118},
	{"chartreuse2", 112},
	// {"chartreuse2", 82}, // duplicate
	{"chartreuse3", 70},
	// {"chartreuse3", 76}, // duplicate
	{"chartreuse4", 64},
	{"cornflowerblue", 69},
	{"cornsilk1", 230},
	{"cyan1", 51},
	{"cyan2", 50},
	{"cyan3", 43},
	{"darkblue", 18},
	{"darkcyan", 36},
	{"darkgoldenrod", 136},
	{"darkgreen", 22},
	{"darkkhaki", 143},
	{"darkmagenta", 90},
	// {"darkmagenta", 91}, // duplicate
	{"darkolivegreen1", 191},
	// {"darkolivegreen1", 192}, // duplicate
	{"darkolivegreen2", 155},
	{"darkolivegreen3", 107},
	// {"darkolivegreen3", 113}, // duplicate
	// {"darkolivegreen3", 149}, // duplicate
	{"darkorange", 208},
	{"darkorange3", 130},
	// {"darkorange3", 166}, // duplicate
	{"darkred", 52},
	// {"darkred", 88}, // duplicate
	{"darkseagreen", 108},
	{"darkseagreen1", 158},
	// {"darkseagreen1", 193}, // duplicate
	{"darkseagreen2", 151},
	// {"darkseagreen2", 157}, // duplicate
	{"darkseagreen3", 115},
	// {"darkseagreen3", 150}, // duplicate
	{"darkseagreen4", 65},
	// {"darkseagreen4", 71}, // duplicate
	{"darkslategray1", 123},
	{"darkslategray2", 87},
	{"darkslategray3", 116},
	{"darkturquoise", 44},
	{"darkviolet", 128},
	// {"darkviolet", 92}, // duplicate
	{"deeppink1", 198},
	// {"deeppink1", 199}, // duplicate
	{"deeppink2", 197},
	{"deeppink3", 161},
	// {"deeppink3", 162}, // duplicate
	{"deeppink4", 125},
	// {"deeppink4", 53}, // duplicate
	// {"deeppink4", 89}, // duplicate
	{"deepskyblue1", 39},
	{"deepskyblue2", 38},
	{"deepskyblue3", 31},
	// {"deepskyblue3", 32}, // duplicate
	{"deepskyblue4", 23},
	// {"deepskyblue4", 24}, // duplicate
	// {"deepskyblue4", 25}, // duplicate
	{"dodgerblue1", 33},
	{"dodgerblue2", 27},
	{"dodgerblue3", 26},
	{"fuchsia", 13},
	{"gold1", 220},
	{"gold3", 142},
	// {"gold3", 178}, // duplicate
	{"green", 2},
	{"green1", 46},
	{"green3", 34},
	// {"green3", 40}, // duplicate
	{"green4", 28},
	{"greenyellow", 154},
	{"grey", 8},
	{"grey0", 16},
	{"grey100", 231},
	{"grey11", 234},
	{"grey15", 235},
	{"grey19", 236},
	{"grey23", 237},
	{"grey27", 238},
	{"grey3", 232},
	{"grey30", 239},
	{"grey35", 240},
	{"grey37", 59},
	{"grey39", 241},
	{"grey42", 242},
	{"grey46", 243},
	{"grey50", 244},
	{"grey53", 102},
	{"grey54", 245},
	{"grey58", 246},
	{"grey62", 247},
	{"grey63", 139},
	{"grey66", 248},
	{"grey69", 145},
	{"grey7", 233},
	{"grey70", 249},
	{"grey74", 250},
	{"grey78", 251},
	{"grey82", 252},
	{"grey84", 188},
	{"grey85", 253},
	{"grey89", 254},
	{"grey93", 255},
	{"honeydew2", 194},
	{"hotpink", 205},
	// {"hotpink", 206}, // duplicate
	{"hotpink2", 169},
	{"hotpink3", 132},
	// {"hotpink3", 168}, // duplicate
	{"indianred", 131},
	// {"indianred", 167}, // duplicate
	{"indianred1", 203},
	// {"indianred1", 204}, // duplicate
	{"khaki1", 228},
	{"khaki3", 185},
	{"lightcoral", 210},
	{"lightcyan1", 195},
	{"lightcyan3", 152},
	{"lightgoldenrod1", 227},
	{"lightgoldenrod2", 186},
	// {"lightgoldenrod2", 221}, // duplicate
	// {"lightgoldenrod2", 222}, // duplicate
	{"lightgoldenrod3", 179},
	{"lightgreen", 119},
	// {"lightgreen", 120}, // duplicate
	{"lightpink1", 217},
	{"lightpink3", 174},
	{"lightpink4", 95},
	{"lightsalmon1", 216},
	{"lightsalmon3", 137},
	// {"lightsalmon3", 173}, // duplicate
	{"lightseagreen", 37},
	{"lightskyblue1", 153},
	{"lightskyblue3", 109},
	// {"lightskyblue3", 110}, // duplicate
	{"lightslateblue", 105},
	{"lightslategrey", 103},
	{"lightsteelblue", 147},
	{"lightsteelblue1", 189},
	{"lightsteelblue3", 146},
	{"lightyellow3", 187},
	{"lime", 10},
	{"magenta1", 201},
	{"magenta2", 165},
	//{"magenta2", 200}, // duplicate
	{"magenta3", 127},
	// {"magenta3", 163}, // duplicate
	// {"magenta3", 164}, // duplicate
	{"maroon", 1},
	{"mediumorchid", 134},
	{"mediumorchid1", 171},
	// {"mediumorchid1", 207}, // duplicate
	{"mediumorchid3", 133},
	{"mediumpurple", 104},
	{"mediumpurple1", 141},
	{"mediumpurple2", 135},
	// {"mediumpurple2", 140}, // duplicate
	{"mediumpurple3", 97},
	// {"mediumpurple3", 98}, // duplicate
	{"mediumpurple4", 60},
	{"mediumspringgreen", 49},
	{"mediumturquoise", 80},
	{"mediumvioletred", 126},
	{"mistyrose1", 224},
	{"mistyrose3", 181},
	{"navajowhite1", 223},
	{"navajowhite3", 144},
	{"navy", 4},
	{"navyblue", 17},
	{"olive", 3},
	{"orange1", 214},
	{"orange3", 172},
	{"orange4", 58},
	// {"orange4", 94}, // duplicate
	{"orangered1", 202},
	{"orchid", 170},
	{"orchid1", 213},
	{"orchid2", 212},
	{"palegreen1", 121},
	// {"palegreen1", 156}, // duplicate
	{"palegreen3", 114},
	// {"palegreen3", 77}, // duplicate
	{"paleturquoise1", 159},
	{"paleturquoise4", 66},
	{"palevioletred1", 211},
	{"pink1", 218},
	{"pink3", 175},
	{"plum1", 219},
	{"plum2", 183},
	{"plum3", 176},
	{"plum4", 96},
	{"purple", 129},
	//{"purple", 5}, // duplicate
	//{"purple", 93}, // duplicate
	{"purple3", 56},
	{"purple4", 54},
	//{"purple4", 55}, // duplicate
	{"red", 9},
	{"red1", 196},
	{"red3", 124},
	//{"red3", 160}, // duplicate
	{"rosybrown", 138},
	{"royalblue1", 63},
	{"salmon1", 209},
	{"sandybrown", 215},
	{"seagreen1", 84},
	//{"seagreen1", 85}, // duplicate
	{"seagreen2", 83},
	{"seagreen3", 78},
	{"silver", 7},
	{"skyblue1", 117},
	{"skyblue2", 111},
	{"skyblue3", 74},
	{"slateblue1", 99},
	{"slateblue3", 61},
	//{"slateblue3", 62}, // duplicate
	{"springgreen1", 48},
	{"springgreen2", 42},
	//{"springgreen2", 47}, // duplicate
	{"springgreen3", 35},
	//{"springgreen3", 41}, // duplicate
	{"springgreen4", 29},
	{"steelblue", 67},
	{"steelblue1", 75},
	//{"steelblue1", 81}, // duplicate
	{"steelblue3", 68},
	{"tan", 180},
	{"teal", 6},
	{"thistle1", 225},
	{"thistle3", 182},
	{"turquoise2", 45},
	{"turquoise4", 30},
	{"violet", 177},
	{"wheat1", 229},
	{"wheat4", 101},
	{"white", 15},
	{"yellow", 11},
	{"yellow1", 226},
	{"yellow2", 190},
	{"yellow3", 148},
	//{"yellow3", 184}, // duplicate
	{"yellow4", 100},
	//{"yellow4", 106}, // duplicate
}
