package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ebiiim/monkey/evaluator"
	"github.com/ebiiim/monkey/lexer"
	"github.com/ebiiim/monkey/parser"
)

// PROMPT is the prompt text used in the REPL.
const PROMPT = ">> "

// Start starts a REPL.
func Start(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMPT)
		if ok := sc.Scan(); !ok {
			return
		}
		p := parser.New(lexer.New(sc.Text()))
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		ev := evaluator.Eval(program)
		if ev != nil {
			fmt.Fprintf(out, "%s\n", ev.Inspect())
		}
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}
