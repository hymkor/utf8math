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

type Style interface {
	Conv(rune) rune
	String() string
}

type NoConv struct{}

func (this *NoConv) Conv(c rune) rune { return c }
func (this *NoConv) String() string   { return "" }

type AlphaNum struct {
	Cap, Dig rune
}

func (this *AlphaNum) String() string {
	return fmt.Sprintf("%c-%c %c-%c %c-%c",
		this.Cap,
		this.Cap+('Z'-'A'),
		this.Cap+('Z'-'A'+1),
		this.Cap+('Z'-'A'+1-'z'-'a'),
		this.Dig,
		this.Dig+9)
}

func (this *AlphaNum) Conv(c rune) rune {
	if unicode.IsUpper(c) {
		return this.Cap + (c - 'A')
	} else if unicode.IsLower(c) {
		return this.Cap + ('Z' - 'A' + 1) + (c - 'a')
	} else if unicode.IsDigit(c) {
		return this.Dig + (c - '0')
	} else {
		return c
	}
}

type Alpha struct {
	Cap rune
}

func (this *Alpha) String() string {
	return fmt.Sprintf("%c-%c %c-%c",
		this.Cap,
		this.Cap+('Z'-'A'),
		this.Cap+('Z'-'A'+1),
		this.Cap+('Z'-'A'+1-'z'-'a'))
}

func (this *Alpha) Conv(c rune) rune {
	if unicode.IsUpper(c) {
		return this.Cap + (c - 'A')
	} else if unicode.IsLower(c) {
		return this.Cap + ('Z' - 'A' + 1) + (c - 'a')
	} else {
		return c
	}
}

// https://ja.wikipedia.org/wiki/%E5%9B%B2%E3%81%BF%E6%96%87%E5%AD%97#%E3%80%8C%E2%91%A0%E3%83%BB%E2%92%B6%E3%80%8D%E7%AD%89%E3%81%AB%E3%82%88%E3%82%8B%E4%B8%80%E8%A6%A7

type Enclosed struct {
	Cap, Sml, Zero, One rune
}

func (this *Enclosed) String() string {
	return fmt.Sprintf("%c-%c %c-%c %c-%c",
		this.Cap,
		this.Cap+('Z'-'A'),
		this.Sml,
		this.Sml+('z'-'a'),
		this.Zero,
		this.One+('9'-'1'))
}

func (this *Enclosed) Conv(c rune) rune {
	if unicode.IsUpper(c) {
		return this.Cap + (c - 'A')
	} else if unicode.IsLower(c) {
		return this.Sml + (c - 'a')
	} else if c == '0' {
		return this.Zero
	} else if unicode.IsDigit(c) {
		return this.One + (c - '1')
	} else {
		return c
	}
}

const (
	_NegativeCircledDigitZero                 = '\u24FF' // unused
	_DingbatNegativeCircledDigitOne           = '\u2776'
	_DingbatNegativeCircledSansSerifDigitOne  = '\u278A'
	_DingbatNegativeCircledSansSerifDigitZero = '\U0001F10C'
	_NegativeCircledLatinCapitalLetterA       = '\U0001F150'
)

// unicode table
// made from http://www.asahi-net.or.jp/~ax2s-kmtn/ref/unicode/u1d400.html
var styles = map[string]Style{
	"-bold":                   &AlphaNum{Cap: '\U0001D400', Dig: '\U0001D7CE'},
	"-italic":                 &Alpha{Cap: '\U0001D434'},
	"-bold-italic":            &Alpha{Cap: '\U0001D468'},
	"-script":                 &Alpha{Cap: '\U0001D49C'},
	"-bold-script":            &Alpha{Cap: '\U0001D4D0'},
	"-fraktur":                &Alpha{Cap: '\U0001D504'},
	"-double-struck":          &AlphaNum{Cap: '\U0001D538', Dig: '\U0001D7D8'},
	"-bold-fraktur":           &Alpha{Cap: '\U0001D56C'},
	"-sans-serif":             &AlphaNum{Cap: '\U0001D5A0', Dig: '\U0001D7E2'},
	"-sans-serif-bold":        &AlphaNum{Cap: '\U0001D5D4', Dig: '\U0001D7EC'},
	"-sans-serif-italic":      &Alpha{Cap: '\U0001D608'},
	"-sans-serif-bold-italic": &Alpha{Cap: '\U0001D63C'},
	"-monospace":              &AlphaNum{Cap: '\U0001D670', Dig: '\U0001D7F6'},
	"-enclosed":               &Enclosed{Cap: '\u24B6', Sml: '\u24D0', Zero: '\u24EA', One: '\u2460'},
	"-black-enclosed": &Enclosed{
		Cap:  _NegativeCircledLatinCapitalLetterA,
		Sml:  'a',
		Zero: _DingbatNegativeCircledSansSerifDigitZero,
		One:  _DingbatNegativeCircledSansSerifDigitOne,
	},
}

var rxAlphabet = regexp.MustCompile("[A-Za-z0-9]+")

func replaceString(text string, style Style) string {
	return rxAlphabet.ReplaceAllStringFunc(text, func(s string) string {
		var buffer strings.Builder
		for _, c := range s {
			buffer.WriteRune(style.Conv(c))
		}
		return buffer.String()
	})
}

func replaceFromStdin(style Style) error {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fmt.Println(replaceString(sc.Text(), style))
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
			fmt.Fprintf(os.Stderr, "    %s %s\n", key, val.String())
		}
		fmt.Fprintln(os.Stderr, "A single hyphen(`-`) means reading text from stdin.")
		return nil
	}

	var conv Style = &NoConv{}
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
			fmt.Print(replaceString(s, conv))
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
