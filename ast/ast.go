package ast

import (
	"github.com/jerry-enebeli/proto-parser/parser"
)

type Method struct {
	Name                          string
	InputTypeName, OutPutTypeName string
	TokenPosition                 int
}

type Service struct {
	Name          string
	TokenPosition int
	Empty         bool
	Methods       []Method
}
type Ast struct {
	tokens       []parser.Token
	Package      string
	ProtoVersion string
	Service      Service
}


//NewAst creeates a new instance of the ast struct and takes in a set of proto tokens
func NewAst(tokens []parser.Token) *Ast{
	return &Ast{tokens: tokens}
}


/*
GenerateAST takes in a custom proto tokens gotten form the parser and extracts the necessary details of the proto file.
It Returns details like the proto syntax version, the package name, the service name and all methods under the service.
*/
func (a *Ast) GenerateAST() {
	tokens := a.tokens
	for index, token := range a.tokens {
		switch token.Value {
		case "package":
			a.Package = tokens[index+1].Value
		case "syntax":
			if a.tokenAtPosition(index, 1).TokenType == parser.TOKENTYPEQUATOR.String() {
				nextTokenValue := a.tokenAtPosition(index, 2).Value
				switch nextTokenValue {
				case "proto3":
					a.ProtoVersion = nextTokenValue
				case "proto2":
					a.ProtoVersion = nextTokenValue

				default:
					panic("invalid proto syntax version")
				}
			}
		case "service":
			if a.tokenAtPosition(index, 1).TokenType == parser.TOKENTYPENAME.String() {
				if a.tokenAtPosition(index, 2).TokenType == parser.TOKENTYPEOPENCURLY.String() {
					if a.tokenAtPosition(index, 3).TokenType == parser.TOKENTYPECLOSECURLY.String() {
						//End of service. Service  is empty
						a.Service = Service{
							Name:          a.tokenAtPosition(index, 1).Value,
							TokenPosition: index + 1,
							Empty:         true,
							Methods:       nil,
						}

					} else {

						if a.tokenAtPosition(index, 3).Value == "rpc" {
							//Handle for rpc methods
							a.Service = Service{
								Name:          a.tokenAtPosition(index, 1).Value,
								Empty:         false,
								TokenPosition: index + 1,
								Methods:       nil,
							}
						} else {
							panic("invalid service definition")
						}

					}
				}
			}
		case "rpc":
			tokenAfterRPCKeyword := a.tokenAtPosition(index, 1)
			if tokenAfterRPCKeyword.TokenType == parser.TOKENTYPENAME.String() {

				method := Method{
					Name:           tokenAfterRPCKeyword.Value,
					InputTypeName:  "",
					OutPutTypeName: "",
					TokenPosition:  index,
				}
				tokenAfterRPCName := a.tokenAtPosition(index, 2)
				if tokenAfterRPCName.TokenType == parser.TOKENTYPEOPENPAREN.String() {

					tokenAfterOpenParen := a.tokenAtPosition(index, 3)
					if tokenAfterOpenParen.TokenType == parser.TOKENTYPENAME.String() {

						method.InputTypeName = tokenAfterOpenParen.Value

						tokenAfterInputTypeName := a.tokenAtPosition(index, 4)

						if tokenAfterInputTypeName.TokenType == parser.TOKENTYPECLOSEPAREN.String() {

							tokenAfterCloseParen := a.tokenAtPosition(index, 5)
							if tokenAfterCloseParen.Value == "returns" {

								tokenAfterReturnValue := a.tokenAtPosition(index, 6)

								if tokenAfterReturnValue.TokenType == parser.TOKENTYPEOPENPAREN.String() {
									tokenAfterOpenParen := a.tokenAtPosition(index, 7)
									if tokenAfterOpenParen.TokenType == parser.TOKENTYPENAME.String() {
										method.OutPutTypeName = tokenAfterOpenParen.Value

										tokenAfterOutPutType := a.tokenAtPosition(index, 8)
										if tokenAfterOutPutType.TokenType == parser.TOKENTYPECLOSEPAREN.String() {

											//TODO: check for curly braces
											a.Service.Methods = append(a.Service.Methods, method)
										}
									}
								}

							}
						} else {
							panic("invalid rpc method definition")
						}
					}

				} else {
					panic("invalid rpc method definition")
				}
			}

		}
	}

}


func (a Ast) tokenAtPosition(currentPosition, move int) parser.Token {
	return a.tokens[currentPosition+move]
}

