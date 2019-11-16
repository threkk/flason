package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/threkk/flason"
)

var version string = "0.0.0"
var output string
var leader string

func getInput(input string) (*os.File, error) {
	// Empty input makes read from STDIN.
	if input == "" {
		return os.Stdin, nil
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	// If there is any type of pipe, we read from it.
	if fi.Mode()&os.ModeNamedPipe != 0 {
		return os.Stdin, nil
	}

	path := filepath.Clean(input)
	if p, err := os.Stat(path); os.IsNotExist(err) || os.IsPermission(err) || p.Mode().IsDir() {
		return nil, errors.New("file cannot be read.")
	}

	return os.Open(path)

}

func init() {
	flag.StringVar(&output, "output", "ini", "format of the output [ini,json,csv]")
	flag.StringVar(&leader, "leader", "", "starter element of each path")
	flag.Usage = func() {
		fmt.Printf(`flason v%s - https://github.com/threkk/flason  

Displays JSON objects read from FILE, or standard 
input, to standard output as path - value pairs.

Usage: flason [flags] <FILE>

Flags:
`, version)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	input := ""
	if flag.NArg() > 0 {
		input = flag.Arg(0)
	}

	file, err := getInput(input)
	if err != nil {
		fmt.Printf("Error opening the file: %s\n", err.Error())
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	content := ""
	for scanner.Scan() {
		content += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading the file: %s\n", err.Error())
		os.Exit(1)
	}

	pairs, err := flason.FlattenJson(content, leader)
	if err != nil {
		fmt.Printf("Error processing file: %s\n", err.Error())
		os.Exit(1)
	}

	switch output {
	case "ini":
		err = pairs.PrintAsIni(os.Stdout)
	case "json":
		err = pairs.PrintAsJson(os.Stdout)
	case "csv":
		err = pairs.PrintAsCsv(os.Stdout)
	default:
		err = fmt.Errorf("Unkown output provided: %s", output)
	}

	if err != nil {
		fmt.Printf("Error writing output: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
