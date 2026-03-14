package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func ParseRelativeSequenceFile(path string) ([]RelativeSequenceEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseRelativeSequences(file)
}

func ParseMappings(r io.Reader) ([]MappingEntry, error) {
	return parseMappings(r)
}

func ParseSequences(r io.Reader) ([]SequenceEntry, error) {
	return parseSequences(r)
}

func ParseRelativeSequences(r io.Reader) ([]RelativeSequenceEntry, error) {
	return parseRelativeSequences(r)
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
	relativeEntries, err := parseRelativeSequences(scannerSource)
	if err != nil {
		return nil, err
	}

	entries := make([]SequenceEntry, 0, len(relativeEntries))
	for _, entry := range relativeEntries {
		entries = append(entries, SequenceEntry{
			Input:  entry.Input,
			Output: entry.Output,
		})
	}

	return entries, nil
}

func parseRelativeSequences(scannerSource io.Reader) ([]RelativeSequenceEntry, error) {
	scanner := bufio.NewScanner(scannerSource)
	var entries []RelativeSequenceEntry

	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields, got %d", lineNo, len(fields))
		}

		entries = append(entries, RelativeSequenceEntry{
			Input:  fields[0],
			Output: fields[1],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func LoadRelativeLayerSequences(mappingPath string, mappings []MappingEntry) ([]SequenceEntry, error) {
	baseDir := filepath.Dir(mappingPath)
	loadedLayers := make(map[string][]RelativeSequenceEntry)
	var sequences []SequenceEntry

	for _, entry := range mappings {
		if !entry.IsLayer {
			continue
		}

		relativeEntries, ok := loadedLayers[entry.Layer]
		if !ok {
			path := filepath.Join(baseDir, "layer_"+entry.Layer)
			if _, err := os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					loadedLayers[entry.Layer] = nil
					continue
				}
				return nil, err
			}

			parsedEntries, err := ParseRelativeSequenceFile(path)
			if err != nil {
				return nil, fmt.Errorf("parse %s: %w", path, err)
			}
			relativeEntries = parsedEntries
			loadedLayers[entry.Layer] = relativeEntries
		}

		for _, relativeEntry := range relativeEntries {
			sequences = append(sequences, SequenceEntry{
				Input:  entry.Physical + relativeEntry.Input,
				Output: relativeEntry.Output,
			})
		}
	}

	return sequences, nil
}
