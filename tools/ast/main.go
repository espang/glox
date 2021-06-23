package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Cut cuts s around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func Cut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func defineAST(outputDir string, asts ...string) error {
	filename := filepath.Join(outputDir, "expressions.go")
	_, err := os.Stat(filename)
	var targetFile *os.File
	if err != nil {
		if os.IsNotExist(err) {
			targetFile, err = os.Create(filename)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		filename := filepath.Join(outputDir, "expression_tmp.go")
		targetFile, err = os.Create(filename)
		if err != nil {
			return err
		}
	}

	targetFile.WriteString("package parser")
	targetFile.WriteString("\n")
	targetFile.WriteString("import (")
	targetFile.WriteString("\t\"fmt\"")
	targetFile.WriteString(")")
	targetFile.WriteString("\n")

	for _, ast := range asts {
		b, a, ok := Cut(ast, ":")
		if !ok {
			return errors.New("expect everey ast to contain one colon")
		}
		typeName := strings.TrimSpace(b)
		fields := strings.TrimSpace(a)

		targetFile.WriteString("type " + typeName + " struct {")
		for _, field := range strings.Split(fields, ", ") {
			targetFile.WriteString("\t" + field)
		}
		targetFile.WriteString("}")
	}

	return nil
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
