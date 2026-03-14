package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stderr, os.Stdout))
}

func run(args []string, stderr io.Writer, stdout io.Writer) int {
	fs := flag.NewFlagSet("googleime-romatable-generator", flag.ContinueOnError)
	fs.SetOutput(stderr)

	inputPath := fs.String("input", "", "path to key mapping definition")
	sequencesPath := fs.String("sequences", "", "path to sequence definition")
	format := fs.String("format", "google", "output format: google or yaml")
	outputPath := fs.String("output", "", "output path (defaults to stdout)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *inputPath == "" {
		fmt.Fprintln(stderr, "missing required -input")
		return 2
	}

	mappings, err := ParseMappingFile(*inputPath)
	if err != nil {
		fmt.Fprintf(stderr, "parse mappings: %v\n", err)
		return 1
	}

	sequences, err := LoadRelativeLayerSequences(*inputPath, mappings)
	if err != nil {
		fmt.Fprintf(stderr, "load layer sequences: %v\n", err)
		return 1
	}

	if *sequencesPath != "" {
		manualSequences, err := ParseSequenceFile(*sequencesPath)
		if err != nil {
			fmt.Fprintf(stderr, "parse sequences: %v\n", err)
			return 1
		}
		sequences = append(sequences, manualSequences...)
	}

	compiled, err := Compile(mappings, sequences)
	if err != nil {
		fmt.Fprintf(stderr, "compile rules: %v\n", err)
		return 1
	}

	var out []byte
	switch *format {
	case "google":
		out = []byte(EmitGoogle(compiled))
	case "yaml":
		out = []byte(EmitYAML(compiled))
	default:
		fmt.Fprintf(stderr, "unsupported -format %q\n", *format)
		return 2
	}

	if *outputPath == "" {
		if _, err := stdout.Write(out); err != nil {
			fmt.Fprintf(stderr, "write output: %v\n", err)
			return 1
		}
		return 0
	}

	if err := os.WriteFile(*outputPath, out, 0o644); err != nil {
		fmt.Fprintf(stderr, "write file: %v\n", err)
		return 1
	}

	return 0
}
