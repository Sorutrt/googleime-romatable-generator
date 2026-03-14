package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ParseMappingFile(path string) ([]MappingEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseMappings(file)
}

func ParseSequenceFile(path string) ([]SequenceEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseSequences(file)
}

func ParseMappings(r io.Reader) ([]MappingEntry, error) {
	return parseMappings(r)
}

func ParseSequences(r io.Reader) ([]SequenceEntry, error) {
	return parseSequences(r)
}

func parseMappings(scannerSource io.Reader) ([]MappingEntry, error) {
	scanner := bufio.NewScanner(scannerSource)
	var entries []MappingEntry

	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields, got %d", lineNo, len(fields))
		}

		entry := MappingEntry{
			Physical: fields[0],
			Value:    fields[1],
		}

		if strings.HasPrefix(entry.Value, "*") {
			entry.IsLayer = true
			entry.Layer = strings.TrimPrefix(entry.Value, "*")
			if entry.Layer == "" {
				return nil, fmt.Errorf("line %d: empty layer reference", lineNo)
			}
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func parseSequences(scannerSource io.Reader) ([]SequenceEntry, error) {
	scanner := bufio.NewScanner(scannerSource)
	var entries []SequenceEntry

	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields, got %d", lineNo, len(fields))
		}

		entries = append(entries, SequenceEntry{
			Input:  fields[0],
			Output: fields[1],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
