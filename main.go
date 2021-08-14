package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// unicode table
// made from http://www.asahi-net.or.jp/~ax2s-kmtn/ref/unicode/u1d400.html
var styles = map[string][2]rune{
	"-bold":                   {'\U0001D400', '\U0001D7CE'},
	"-italic":                 {'\U0001D434', '0'},
	"-bold-italic":            {'\U0001D468', '0'},
	"-script":                 {'\U0001D49C', '0'},
	"-bold-script":            {'\U0001D4D0', '0'},
	"-fraktur":                {'\U0001D504', '0'},
	"-double-struck":          {'\U0001D538', '\U0001D7D8'},
	"-bold-fraktur":           {'\U0001D56C', '0'},
	"-sans-serif":             {'\U0001D5A0', '\U0001D7E2'},
	"-sans-serif-bold":        {'\U0001D5D4', '\U0001D7EC'},
	"-sans-serif-italic":      {'\U0001D608', '0'},
	"-sans-serif-bold-italic": {'\U0001D63C', '0'},
	"-monospace":              {'\U0001D670', '\U0001D7F6'},
}

var rxAlphabet = regexp.MustCompile("[A-Za-z0-9]+")

func replaceAlphabet(text string, capital, small, digit rune) string {
	return rxAlphabet.ReplaceAllStringFunc(text, func(s string) string {
		var buffer strings.Builder
		for _, c := range s {
			if unicode.IsUpper(c) {
				buffer.WriteRune(capital + (c - 'A'))
			} else if unicode.IsLower(c) {
				buffer.WriteRune(small + (c - 'a'))
			} else {
				buffer.WriteRune(digit + (c - '0'))
			}
		}
		return buffer.String()
	})
}

func replaceFromStdin(capital, small, digit rune) error {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fmt.Println(replaceAlphabet(sc.Text(), capital, small, digit))
	}
	if err := sc.Err(); err != io.EOF {
		return err
	}
	return nil
}

func mains(args []string) error {
	if len(args) <= 0 {
		fmt.Fprintln(os.Stderr, "Convert Alphabets and digits to Mathematical Symbols")
		fmt.Fprintf(os.Stderr, "Usage: %s {-OPTIONS TEXT}\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "(FONT-OPTIONS)")
		for key, val := range styles {
			fmt.Fprintf(os.Stderr, "    %s %c-%c %c-%c %c-%c\n",
				key,
				val[0],
				val[0]+('Z'-'A'),
				val[0]+('Z'-'A'+1),
				val[0]+('Z'-'A'+1+'z'-'a'),
				val[1],
				val[1]+9)
		}
		fmt.Fprintln(os.Stderr, "A single hyphen(`-`) means reading text from stdin.")
		return nil
	}

	capital := rune('A')
	small := rune('a')
	digit := rune('0')

	delimiter := ""
	for _, s := range args {
		if len(s) > 0 && s[0] == '-' {
			if start, ok := styles[s]; ok {
				capital = start[0]
				small = start[0] + ('Z' - 'A' + 1)
				digit = start[1]
			} else if s == "-" {
				if err := replaceFromStdin(capital, small, digit); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("%s: unknown option", s)
			}
		} else {
			fmt.Print(delimiter)
			fmt.Print(replaceAlphabet(s, capital, small, digit))
			delimiter = " "
		}
	}
	fmt.Println()
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
