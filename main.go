package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"lox/parser"
	"os"
)

//var average = (min+max)/2;

func main() {

	if len(os.Args) > 2 {
		fmt.Println("Usage lox [script]")
		os.Exit(64)
	}

	if len(os.Args) == 2 {
		err := runFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(65)
		}
		os.Exit(0)
	}
	runPrompt()
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		err := run(scanner.Text())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("> ")
	}

	// add a newline after Ctrl+D
	fmt.Println("Lox is done! canceled by user via Ctrl+D")

	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner error: %v\n", err)
	}
}

func run(content string) error {
	tokens, err := parser.Lex(content)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
	return nil
}

func runFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return run(string(content))
}
