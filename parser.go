package main

import (
	"fmt"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
)

type Exp struct {
	Properties []*Property `"{" @@* "}"`
}

type Property struct {
	Key   string `@Ident ":"`
	Value *Value `@@`
}

type Value struct {
	String   *string   `  @String`
	Number   *float64  `| @Number`
	Property *Property `| "{" @@ "}"`
}

func main() {
	basicLexer := lexer.Must(ebnf.New(`
	Ident = (alpha | "$") { "_" | alpha | digit } .
	String = "\"" { "\u0000"…"\uffff"-"\""-"\\" | "\\" any } "\"" .
	Number = [ "-" | "+" ] ("." | digit) { "." | digit } .
	Punct = "!"…"/" | ":"…"@" | "["…` + "\"`\"" + ` | "{"…"~" .
	EOL = ( "\n" | "\r" ) { "\n" | "\r" }.
	Whitespace = ( " " | "\t" ) { " " | "\t" } .

	alpha = "a"…"z" | "A"…"Z" .
	digit = "0"…"9" .
	any = "\u0000"…"\uffff" .
`))

	parser := participle.MustBuild(&Exp{},
		participle.Lexer(basicLexer),
		participle.CaseInsensitive("Ident"),
		participle.Unquote("String"),
		participle.UseLookahead(2),
		participle.Elide("Whitespace"),
	)

	// parser, err := participle.Build(&Exp{})
	// if err != nil {
	// 	fmt.Printf("%s\n", err.Error())
	// }
	exp := &Exp{}
	err := parser.ParseString(`{ R_STAT: 10 }`, exp)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Printf("%#v\n", exp)

	exp2 := &Exp{}
	err = parser.ParseString(`{ ERR_S: { $gte: 1 } }`, exp2)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Printf("%#v\n", exp2)

}
