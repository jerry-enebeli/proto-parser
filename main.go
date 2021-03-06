package main

import (
	"encoding/json"
	"fmt"
	"github.com/jerry-enebeli/proto-parser/ast"
	"github.com/jerry-enebeli/proto-parser/parser"
)

func main() {
	p := parser.NewParser("test.proto")
	tokens := p.Tokens

	a := ast.NewAst(tokens)

	a.GenerateAST()

	jsonRep, _ := json.Marshal(a)
	fmt.Println(string(jsonRep))

}
