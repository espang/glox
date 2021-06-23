package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func defineAST(outputDir string, asts ...string) error {
	filename := filepath.Join(outputDir, "expressions.go")
	_, err := os.Stat(filename)
	var targetFile os.File
	if err != nil {
		if os.IsNotExist(err) {
			// create
		} else {
			return err
		}
	} else {

	}
}

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("Usage: ast <output directory>")
		os.Exit(64)
	}

	outputFolder := args[1]

	folderInfo, err := os.Stat(outputFolder)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(outputFolder, 0755)
			if err != nil {
				fmt.Printf("couldn't make folder '%s': %v\n", outputFolder, err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("unexpected error: %v\n", err)
			os.Exit(1)
		}
	}
	if !folderInfo.IsDir() {
		fmt.Printf("'%s' is not a folder\n", outputFolder)
		os.Exit(1)
	}

	defineAST(
		outputFolder,
		"Binary        : Expr left, Token operator, Expr right",
		"Grouping      : Expr expression",
		"StringLiteral : string value",
		"NumberLiteral : float64 value",
		"BoolLiteral   : bool value",
		"Unary         : Token operator, Expr right",
	)
}
