package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/ebiiim/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\nFeel free to type in commands\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
