package css

import (
	"fmt"

	v4 "github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot/helper/antlr"
)

func ReadCSS(sourceCode, path string) (map[string]struct{}, error) {
	build := NewCSSBuild(sourceCode)
	rel := build.RelativePath(path)

	lex := Newcss3Lexer(nil)
	parse := Newcss3Parser(nil)
	el := antlr.AntlrSetupLexParse(rel, lex.BaseLexer, parse.BaseParser)
	el.CurrentFile = rel
	sheet := parse.Stylesheet()
	v4.ParseTreeWalkerDefault.Walk(build, sheet)
	if el.Failed() {
		return nil, fmt.Errorf("failed reading CSS file '%s'", path)
	}
	// if builder.Failed() {
	// 	return nil, b.CurrentFile, false
	// }
	v4.ParseTreeWalkerDefault.Walk(build, sheet)
	return build.ClassName, nil
}

// func xxreadInput(sourceCode, path string) (css3Listener, *css3Parser) {
// 	build := NewCSSBuild(sourceCode)
// 	rel := build.RelativePath(path)

// 	fp, err := os.Open(rel)
// 	if err != nil {
// 		wd, _ := os.Getwd()
// 		log.Fatalf("%v (wd is %s), %v", flag.Arg(0), wd, err)
// 	}
// 	buffer, err := io.ReadAll(fp)
// 	if err != nil {
// 		log.Fatalf("reading %s: %v", flag.Arg(0), err)
// 	}
// 	fp.Close()
// 	//el := errorListener{0}
// 	input := antlr.NewInputStream(string(buffer))
// 	lexer := Newcss3Lexer(input)
// 	//lexer.RemoveErrorListeners()
// 	stream := antlr.NewCommonTokenStream(lexer, 0)
// 	p := Newcss3Parser(stream)
// 	p.RemoveErrorListeners()
// 	// the diagnostic listener is good for debugging (displays good error msgs)
// 	p.AddErrorListener(&antlr.DiagnosticErrorListener{
// 		DefaultErrorListener: &antlr.DefaultErrorListener{},
// 	})
// 	//p.AddErrorListener(&el)
// 	return build, p
// }
