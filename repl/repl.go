package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ebiiim/monkey/evaluator"
	"github.com/ebiiim/monkey/lexer"
	"github.com/ebiiim/monkey/object"
	"github.com/ebiiim/monkey/parser"
)

// PROMPT is the prompt text used in the REPL.
const PROMPT = ">> "

// Start starts a REPL.
func Start(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Fprint(out, PROMPT)
		if ok := sc.Scan(); !ok {
			return
		}
		p := parser.New(lexer.New(catchREPLCommands(out, sc.Text())))
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		ev := evaluator.Eval(program, env)
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

func catchREPLCommands(out io.Writer, input string) string {
	if len(input) == 0 || input[0] != '.' {
		return input
	}
	ss := strings.Split(input, " ")
	switch ss[0] {
	case ".exit":
		return exit(out)
	case ".load":
		if len(ss) < 2 {
			return help(out)
		}
		return load(out, ss[1])
	default:
		return help(out)
	}
}

func help(out io.Writer) string {
	fmt.Fprint(out, "\tREPL Commands: [ .exit | .load FILE ]\n")
	return ""
}

func exit(out io.Writer) string {
	fmt.Fprint(out, "\tbye\n")
	os.Exit(0)
	return "" // unreachable but fulfill the interface
}

func load(out io.Writer, fileName string) string {
	p, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(out, "\tFailed to load %s: %v\n", fileName, err)
	}
	return string(p)
}
