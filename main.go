package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"terraform-plan-summary/parser"
	"terraform-plan-summary/reader"
	"terraform-plan-summary/writer"
)

var version = "development"

func main() {
	printVersion := flag.Bool("v", false, "print version")
	tree := flag.Bool("tree", false, "[Optional] print changes in tree format")
	separateTree := flag.Bool("separate-tree", false, "[Optional] print changes in tree format for add/delete/change/recreate changes")
	drawable := flag.Bool("draw", false, "[Optional, used only with -tree or -separate-tree] draw trees instead of plain tree")
	outputFileName := flag.String("out", "", "[Optional] write output to file")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s [args] [tf-plan.json|tfplan]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *printVersion {
		_, _ = fmt.Fprintf(os.Stdout, fmt.Sprintf("Version: %s\n", version))
		os.Exit(0)
	}

	args := flag.Args()
	err := validateFlags(*tree, *separateTree, *drawable, args)
	logIfErrorAndExit("invalid input flags: %s\n", err, flag.Usage)

	newReader, err := reader.CreateReader(os.Stdin, args)
	logIfErrorAndExit("error creating input reader: %s\n", err, flag.Usage)

	input, err := newReader.Read()
	logIfErrorAndExit("error reading from input: %s", err, func() {})

	newParser, err := parser.CreateParser(input, newReader.Name())
	logIfErrorAndExit("error creating parser: %s", err, func() {})

	terraformState, err := newParser.Parse()
	logIfErrorAndExit("%s", err, func() {})

	terraformState.FilterNoOpResources()

	newWriter := writer.CreateWriter(*tree, *separateTree, *drawable, terraformState)

	var outputFile io.Writer = os.Stdout

	if *outputFileName != "" {
		file, err := os.OpenFile(*outputFileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
		logIfErrorAndExit("error opening file: %v", err, func() {})
		defer func() {
			if err := file.Close(); err != nil {
				logIfErrorAndExit("Error closing file: %s\n", err, func() {})
			}
		}()
		outputFile = file
	}

	err = newWriter.Write(outputFile)
	logIfErrorAndExit("error writing: %s", err, func() {})

	if err == nil && *outputFileName != "" {
		_, _ = fmt.Fprintf(os.Stderr, "Written plan summary to %s\n", *outputFileName)
	}
}

func logIfErrorAndExit(format string, err error, callback func()) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", format), err.Error())
		callback()
		os.Exit(1)
	}
}

func validateFlags(tree, separateTree, drawable bool, args []string) error {
	if tree && separateTree {
		return fmt.Errorf("both -tree and -seperate-tree should not be provided")
	}
	if !tree && !separateTree && drawable {
		return fmt.Errorf("drawable should be provided with -tree or -seperate-tree")
	}
	if len(args) > 1 {
		return fmt.Errorf("only one argument is allowed which is filename, but got %v", args)
	}
	return nil
}
