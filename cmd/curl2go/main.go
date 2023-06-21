/*
Copyright © 2023 lonjaju
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/lonjaju/curl2go"
)

var (
	input  string
	output string
	help   bool
)

func init() {
	flag.StringVar(&input, "input", "", "the path to the input file; if not specified stdin is used")
	flag.StringVar(&input, "i", "", "the short flag for -input")
	flag.StringVar(&output, "output", "", "path to the output file; if not specified stdout is used")
	flag.StringVar(&output, "o", "", "the short flag for -output")
	flag.BoolVar(&help, "help", false, "json2go help")
	flag.BoolVar(&help, "h", false, "the short flag for -help")
}

func main() {
	flag.Parse()
	args := flag.Args()
	// the only arg we care about is help.  This is in case the user uses
	// just help instead of -help or -h
	for _, arg := range args {
		if arg == "help" {
			help = true
			break
		}
	}
	if help {
		ShowHelp()
		os.Exit(0)
	}

	var in io.Reader
	var out io.Writer
	var err error

	if input == "" {
		in = os.Stdin
	} else {
		in, err = os.Open(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open input file fail: %s", err.Error())
			os.Exit(1)
		}
	}
	if output == "" {
		out = os.Stdout
	} else {
		out, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0754)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open output file fail: %s", err.Error())
			os.Exit(1)
		}
	}

	ir := bufio.NewReader(in)
	io := bufio.NewWriter(out)

	input, err := ir.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "read input fail: %s", err.Error())
		os.Exit(1)
	}

	goCode, err := curl2go.Curl2go(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "convert fail: %s", err.Error())
		os.Exit(1)
	}

	io.WriteString(goCode)
	io.Flush()
}

func ShowHelp() {
	helpText := `
Usage: curl2go [options]

A tool Convert curl to go code.

A curl source file can be specified with either the -i or -input
flags.  If none is specified, the JSON is expected to come from
stdin.

The output file of the generated Go source code is specified
with either the -o or -output flags.  If none is specified, the
output will be written to stdout.

Errors are written to stderr.

Simple examples:

    $ echo "curl -X POST https://reqbin.com/echo/post/json 
	-H "Content-Type: application/json"
	-d '{"productId": 123456, "quantity": 100}' " | ./curl2go -o example.go

Options:

flag              default   description
---------------   -------   ------------------------------------------
-i  -input        stdin     The JSON input source.
-o  -output       stdout    The Go srouce code output destination.
-h  -help         false     Print the help text; 'help' is also valid.

`
	fmt.Println(helpText)
}
