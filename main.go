package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)


/**
TODO: create an index of keywords
TODO: compare keywords index with TOKENTYPENAME before pushing to token
**/


type tokenType int

func (ty tokenType) String() string {
	return [...]string{"name", "equal", "semicolon","opencurly", "closecurly",  "openparen", "closeparen"}[ty]
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

type protoDetails struct {
	Package string
	ProtoVersion string
	Service string
}

type token struct {
	tokenType string
	value     string
}

type parser struct {
	file        string
	current     int
	fileAttr    []rune
	currentChar string
	totalChars int
	tokens      []token
}

//move to the next char
func (p *parser) next() {

	if p.current+1 == p.totalChars{
		p.currentChar = "!" //TODO find the better way to map end of file
	} else {
		p.current++
		p.currentChar = string(p.fileAttr[p.current])
	}

}

//get the next char
func (p parser) getNext() string {
	return string(p.fileAttr[p.current+1])
}

func newParser(protoFile string) *parser {
	p := &parser{
		file:        "",
		current:     0,
		fileAttr:    nil,
		currentChar: "",
		tokens:      nil,
	}

	p.readFile(protoFile)

	//read each char till the end of file
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

		}

		break
	}

	return p
}

//read proto file
func (p *parser) readFile(file string) {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	p.file = string(data) + "\n"
	p.fileAttr = []rune(p.file)
	p.currentChar = string(p.fileAttr[p.current])
	p.totalChars = len(p.file)
}

func (p *parser) pushToToken(tokenType, value string) {
	p.tokens = append(p.tokens, token{tokenType: tokenType, value: value})
}


//CheckName checks for keywords name
func (p *parser) CheckName() bool {
	//Check for name token
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
func (p *parser) CheckWhiteSpace() bool {
	WHITESPACE := regexp.MustCompile(`\s`)
	if WHITESPACE.Match([]byte(p.currentChar)) {
		p.next()
		return true
	}

	return false
}

//CheckEquator checks for equator sign '='
func (p *parser) CheckEquator() bool {
	if p.currentChar == "=" {
		p.pushToToken(TOKENTYPEQUATOR.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}

//CheckSemiColon checks for semicolon ';'
func (p *parser) CheckSemiColon() bool {
	if p.currentChar == ";" {
		p.pushToToken(TOKENTYPESEMICOLON.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}


//CheckOpenCurlyBraces checks for open curly braces '}'
func (p *parser) CheckOpenCurlyBraces() bool {
	if p.currentChar == "{" {
		p.pushToToken(TOKENTYPEOPENCURLY.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}


//CheckCloseCurlyBraces checks for closing curly braces '}'
func (p *parser) CheckCloseCurlyBraces() bool {
	if p.currentChar == "}" {
		p.pushToToken(TOKENTYPECLOSECURLY.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}


//CheckOpenParenBraces checks for open parentheses '('
func (p *parser) CheckOpenParenBraces() bool {
	if p.currentChar == "(" {
		p.pushToToken(TOKENTYPEOPENPAREN.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}


//CheckCloseParenBraces checks for closing parentheses ')'
func (p *parser) CheckCloseParenBraces() bool {
	if p.currentChar == ")" {
		p.pushToToken(TOKENTYPECLOSEPAREN.String(), p.currentChar)
		p.next()
		return true
	}
	return false
}


//CheckString checks for string starting with a quote(") and ignore it
func (p *parser) CheckString() bool {
	if p.currentChar == `"` {
		p.next()
		return true
	}
	return false
}



/*
GenerateProtoDetails takes in a custom proto tokens gotten form the parser and extracts the necessary details of the proto file.
It Returns details like the proto syntax version, the package name, the service name and all methods under the service.
*/
func (d *protoDetails) GenerateProtoDetails(tokens []token, ) {
	for index, token := range tokens {

		//check for package token and move to the next node
		if token.value == "package" {
			d.Package = tokens[index+1].value
		}

		//checks for proto syntax version
		// syntax = "proto3";
		if token.value == "syntax"  && tokens[index+1].tokenType == TOKENTYPEQUATOR.String() {
			d.ProtoVersion = tokens[index+2].value
		}

		//check for service name
		//service Greeter {
		if token.value == "service" && tokens[index+1].tokenType == TOKENTYPENAME.String() && tokens[index+2].tokenType==TOKENTYPEOPENCURLY.String() {
			d.Service = tokens[index+1].value
		}
	}

}
func main() {
	p := newParser("test.proto")
	tokens := p.tokens

	//Generate
	d := &protoDetails{}
	d.GenerateProtoDetails(tokens)

	jsonRep ,_:= json.Marshal(d)
	fmt.Println(string(jsonRep), tokens)

}
