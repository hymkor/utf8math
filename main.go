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

func alphanum(capitalStart, digitStart rune) func(rune) rune {
	smallStart := capitalStart + 26
	return func(c rune) rune {
		if unicode.IsUpper(c) {
			return capitalStart + (c - 'A')
		} else if unicode.IsLower(c) {
			return smallStart + (c - 'a')
		} else if unicode.IsDigit(c) {
			return digitStart + (c - '0')
		} else {
			return c
		}
	}
}

func alpha(capitalStart rune) func(rune) rune {
	smallStart := capitalStart + 26
	return func(c rune) rune {
		if unicode.IsUpper(c) {
			return capitalStart + (c - 'A')
		} else if unicode.IsLower(c) {
			return smallStart + (c - 'a')
		} else {
			return c
		}
	}
}

// https://ja.wikipedia.org/wiki/%E5%9B%B2%E3%81%BF%E6%96%87%E5%AD%97#%E3%80%8C%E2%91%A0%E3%83%BB%E2%92%B6%E3%80%8D%E7%AD%89%E3%81%AB%E3%82%88%E3%82%8B%E4%B8%80%E8%A6%A7

func enclosed(c rune) rune {
	if unicode.IsUpper(c) {
		return '\u24B6' + (c - 'A')
	} else if unicode.IsLower(c) {
		return '\u24D0' + (c - 'a')
	} else if c == '0' {
		return '\u24EA'
	} else if unicode.IsDigit(c) {
		return '\u245F' + (c - '0')
	} else {
		return c
	}
}

// unicode table
// made from http://www.asahi-net.or.jp/~ax2s-kmtn/ref/unicode/u1d400.html
var styles = map[string](func(rune) rune){
	"-bold":                   alphanum('\U0001D400', '\U0001D7CE'),
	"-italic":                 alpha('\U0001D434'),
	"-bold-italic":            alpha('\U0001D468'),
	"-script":                 alpha('\U0001D49C'),
	"-bold-script":            alpha('\U0001D4D0'),
	"-fraktur":                alpha('\U0001D504'),
	"-double-struck":          alphanum('\U0001D538', '\U0001D7D8'),
	"-bold-fraktur":           alpha('\U0001D56C'),
	"-sans-serif":             alphanum('\U0001D5A0', '\U0001D7E2'),
	"-sans-serif-bold":        alphanum('\U0001D5D4', '\U0001D7EC'),
	"-sans-serif-italic":      alpha('\U0001D608'),
	"-sans-serif-bold-italic": alpha('\U0001D63C'),
	"-monospace":              alphanum('\U0001D670', '\U0001D7F6'),
	"-enclosed":               enclosed,
}

var rxAlphabet = regexp.MustCompile("[A-Za-z0-9]+")

func replaceAlphabet(text string, conv func(rune) rune) string {
	return rxAlphabet.ReplaceAllStringFunc(text, func(s string) string {
		var buffer strings.Builder
		for _, c := range s {
			buffer.WriteRune(conv(c))
		}
		return buffer.String()
	})
}

func replaceFromStdin(conv func(rune) rune) error {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fmt.Println(replaceAlphabet(sc.Text(), conv))
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
		for key, _ := range styles {
			fmt.Fprintf(os.Stderr, "    %s\n", key)
		}
		fmt.Fprintln(os.Stderr, "A single hyphen(`-`) means reading text from stdin.")
		return nil
	}

	conv := func(c rune) rune { return c }
	delimiter := ""
	for _, s := range args {
		if len(s) > 0 && s[0] == '-' {
			if _conv, ok := styles[s]; ok {
				conv = _conv
			} else if s == "-" {
				if err := replaceFromStdin(conv); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("%s: unknown option", s)
			}
		} else {
			fmt.Print(delimiter)
			fmt.Print(replaceAlphabet(s, conv))
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
