package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ebiiim/monkey/lexer"
	"github.com/ebiiim/monkey/token"
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
		ln := sc.Text()
		l := lexer.New(ln)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
