package curl2go

import (
	"strings"
)

var defParserOptions = ParserOptions{
	Alias: OptionAliases,
	Bool:  BoolOptions,
}

var defParser = NewParser()

type Parser struct {
	cursor int
	input  string

	flags   ParsedFlags
	options ParserOptions
}

type ParserOptions struct {
	Alias map[string]string
	Bool  map[string]bool
}

type ParsedFlags struct {
	UnFlags      []string
	BoolFlags    map[string]bool
	StringsFlags map[string][]string
}

type RelevantData struct {
	Ascii string
	Files []string
}

type BasicAuth struct {
	User string
	Pass string
}

func NewParser() *Parser {
	return NewParserWithOptions(defParserOptions)
}

func NewParserWithOptions(options ParserOptions) *Parser {
	return &Parser{
		cursor:  0,
		options: options,
	}
}

func FlagParse(input string) ParsedFlags {
	return defParser.FlagParse(input)
}

func (p *Parser) FlagParse(input string) ParsedFlags {
	p.input = input

	p.flags.UnFlags = []string{}
	p.cursor = 0
	p.flags = ParsedFlags{UnFlags: []string{}, BoolFlags: map[string]bool{}, StringsFlags: map[string][]string{}}

	p.input = strings.TrimSpace(p.input)
	// Ignore the # $ at the beginning
	if len(p.input) > 2 && (p.input[0] == '$' || p.input[0] == '#') && isWhitespace(p.input[1]) {
		p.input = p.input[1:]
		p.input = strings.TrimSpace(p.input)
	}

	for p.cursor < len(p.input) {
		p.skipWhitespace()

		if p.input[p.cursor] == '-' {
			p.flagSet()
		} else {
			p.unflagged()
		}
	}

	return p.flags
}

func (p *Parser) flagSet() {
	if p.cursor < len(p.input) && p.input[p.cursor+1] == '-' {
		p.longFlag()
		return
	}

	p.cursor++
	for p.cursor < len(p.input) && !isWhitespace(p.input[p.cursor]) {
		flagName := p.fullName(p.input[p.cursor : p.cursor+1])
		if p.flags.StringsFlags[flagName] == nil {
			p.flags.StringsFlags[flagName] = []string{}
		}
		p.cursor++
		if p.boolFlag(flagName) {
			p.flags.BoolFlags[flagName] = toBool(flagName)
		} else {
			p.flags.StringsFlags[flagName] = append(p.flags.StringsFlags[flagName], p.nextWord(" "))
		}
	}
}

// longFlag consumes a "--long-flag" sequence and
// stores it in flags.
func (p *Parser) longFlag() {
	p.cursor += 2
	flagName := p.nextWord("=")
	if p.boolFlag(flagName) {
		p.flags.BoolFlags[flagName] = toBool(flagName)
	} else {
		if p.flags.StringsFlags[flagName] == nil {
			p.flags.StringsFlags[flagName] = []string{}
		}
		p.flags.StringsFlags[flagName] = append(p.flags.StringsFlags[flagName], p.nextWord(" "))
	}
}

// unflagged consumes the next string as an unflagged value,
// storing it in the flags.
func (p *Parser) unflagged() {
	p.flags.UnFlags = append(p.flags.UnFlags, p.nextWord(" "))
}

func (p *Parser) fullName(flag string) string {
	alias, ok := p.options.Alias[flag]
	if ok {
		return alias
	}

	return flag
}

func (p *Parser) boolFlag(flag string) bool {
	return p.options.Bool[flag]
}

func toBool(flag string) bool {
	return !(strings.HasPrefix(flag, "no-") || strings.HasPrefix(flag, "disable-"))
}

// nextWord skips any leading whitespace and consumes the next
// space-delimited string value and returns it. If endChar is set,
// it will be used to determine the end of the string. Normally just
// unescaped whitespace is the end of the string, but endChar can
// be used to specify another end-of-string. This function honors \
// as an escape character and does not include it in the value, except
// in the special case of the \$ sequence, the backslash is retained
// so other code can decide whether to treat as an env var or not.
func (p *Parser) nextWord(endChar string) string {
	p.skipWhitespace()

	var str strings.Builder
	quoted := false
	quoteCh := byte(0)
	escaped := false
	quoteDS := false

	for p.cursor < len(p.input) {
		if quoted {
			if p.input[p.cursor] == quoteCh && !escaped && p.input[p.cursor-1] != '\\' {
				quoted = false
				p.cursor++
				continue
			}
		}
		if !quoted {
			if !escaped {
				if isWhitespace(p.input[p.cursor]) {
					return str.String()
				}
				if p.input[p.cursor] == '"' || p.input[p.cursor] == '\'' {
					quoted = true
					quoteCh = p.input[p.cursor]
					if str.String()+string(quoteCh) == "$'" {
						quoteDS = true
						str.Reset()
					}
					p.cursor++
				}
				if endChar != "" && string(p.input[p.cursor]) == endChar {
					p.cursor++
					return str.String()
				}
			}
		}
		if !escaped && !quoteDS && p.input[p.cursor] == '\\' {
			escaped = true
			if p.cursor < len(p.input)-1 && p.input[p.cursor+1] == '$' {
				// skip the backslash unless the next character is $
				p.cursor++
			}
			p.cursor++
			continue
		}
		str.WriteByte(p.input[p.cursor])
		escaped = false
		p.cursor++
	}

	return str.String()
}

func (p *Parser) skipWhitespace() {
	for p.cursor < len(p.input) && isWhitespace(p.input[p.cursor]) {
		p.cursor++
	}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
