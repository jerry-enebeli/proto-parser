package parser

import (
	"io/ioutil"
	"log"
	"regexp"
)

type tokenType int

func (ty tokenType) String() string {
	return [...]string{"name", "equal", "semicolon", "opencurly", "closecurly", "openparen", "closeparen"}[ty]
}

const (
	TOKENTYPENAME tokenType = iota
	TOKENTYPEQUATOR
	TOKENTYPESEMICOLON
	TOKENTYPEOPENCURLY
	TOKENTYPECLOSECURLY
	TOKENTYPEOPENPAREN
	TOKENTYPECLOSEPAREN
)

type Token struct {
	TokenType string
	Value     string
}

type Parser struct {
	File        string
	current     int
	fileAttr    []rune
	currentChar string
	totalChars  int
	Tokens      []Token
}

//move to the next char
func (p *Parser) next() {

	if p.current+1 == p.totalChars {
		p.currentChar = "!" //TODO find the better way to map end of File
	} else {
		p.current++
		p.currentChar = string(p.fileAttr[p.current])
	}

}

//get the next char
func (p Parser) getNext() string {
	return string(p.fileAttr[p.current+1])
}

func NewParser(protoFile string) *Parser {
	p := &Parser{
		File:        "",
		current:     0,
		fileAttr:    nil,
		currentChar: "",
		Tokens:      nil,
	}

	p.readFile(protoFile)

	//read each char till the end of File
	for p.current < p.totalChars {
		switch {
		case p.CheckName():
			continue

		case p.CheckWhiteSpace():
			continue

		case p.CheckEquator():
			continue

		case p.CheckString():
			continue

		case p.CheckSemiColon():
			continue

		case p.CheckOpenCurlyBraces():
			continue

		case p.CheckCloseCurlyBraces():
			continue
		case p.CheckOpenParenBraces():
			continue
		case p.CheckCloseParenBraces():
			continue

		case p.CheckForComments():
			continue

		}

		break
	}

	return p
}

//read proto File
func (p *Parser) readFile(file string) {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	p.File = string(data) + "\n"
	p.fileAttr = []rune(p.File)
	p.currentChar = string(p.fileAttr[p.current])
	p.totalChars = len(p.File)
}

func (p *Parser) pushToToken(tokenType, value string) {
	p.Tokens = append(p.Tokens, Token{TokenType: tokenType, Value: value})
}

//CheckName checks for keywords name
func (p *Parser) CheckName() bool {
	//Check for name Token
	LETTER := regexp.MustCompile(`^[\w&.-]+$`)
	if LETTER.Match([]byte(p.currentChar)) {
		value := ""

		for LETTER.Match([]byte(p.currentChar)) {
			value += p.currentChar
			p.next()
		}
		p.pushToToken(TOKENTYPENAME.String(), value)

		return true
	}

	return false
}

//CheckWhiteSpace checks for white spaces and ignores them
func (p *Parser) CheckWhiteSpace() bool {
	WHITESPACE := regexp.MustCompile(`\s`)
	if WHITESPACE.Match([]byte(p.currentChar)) {
		p.next()
		return true
	}

	return false
}

//CheckEquator checks for equator sign '='
func (p *Parser) CheckEquator() bool {
	if p.currentChar == "=" {
		p.pushToToken(TOKENTYPEQUATOR.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}

//CheckSemiColon checks for semicolon ';'
func (p *Parser) CheckSemiColon() bool {
	if p.currentChar == ";" {
		p.pushToToken(TOKENTYPESEMICOLON.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}

//CheckOpenCurlyBraces checks for open curly braces '}'
func (p *Parser) CheckOpenCurlyBraces() bool {
	if p.currentChar == "{" {
		p.pushToToken(TOKENTYPEOPENCURLY.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}

//CheckCloseCurlyBraces checks for closing curly braces '}'
func (p *Parser) CheckCloseCurlyBraces() bool {
	if p.currentChar == "}" {
		p.pushToToken(TOKENTYPECLOSECURLY.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}

//CheckOpenParenBraces checks for open parentheses '('
func (p *Parser) CheckOpenParenBraces() bool {
	if p.currentChar == "(" {
		p.pushToToken(TOKENTYPEOPENPAREN.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}

//CheckCloseParenBraces checks for closing parentheses ')'
func (p *Parser) CheckCloseParenBraces() bool {
	if p.currentChar == ")" {
		p.pushToToken(TOKENTYPECLOSEPAREN.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}

func (p *Parser) CheckForComments() bool{
	if p.currentChar == "/" {
		if p.getNext() == "/" {
			for p.currentChar != "\n"{
				p.next()
			}
			return true
		}
	}

	return false
}


//CheckString checks for string starting with a quote(") and ignore it
func (p *Parser) CheckString() bool {
	if p.currentChar == `"` {
		p.next()
		return true
	}
	return false
}