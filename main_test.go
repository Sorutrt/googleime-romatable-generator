package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseMappings(t *testing.T) {
	path := writeTempFile(t, "mappings.txt", "# comment\nq *q\nj h\n\n")

	entries, err := ParseMappingFile(path)
	if err != nil {
		t.Fatalf("ParseMappingFile() error = %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2", len(entries))
	}

	if !entries[0].IsLayer || entries[0].Layer != "q" {
		t.Fatalf("entries[0] = %#v, want layer q", entries[0])
	}

	if entries[1].Value != "h" {
		t.Fatalf("entries[1].Value = %q, want h", entries[1].Value)
	}
}

func TestCompileGeneratesRows(t *testing.T) {
	compiled, err := Compile(
		[]MappingEntry{
			{Physical: "q", Value: "*q", IsLayer: true, Layer: "q"},
			{Physical: "j", Value: "h"},
			{Physical: "g", Value: "f"},
		},
		[]SequenceEntry{
			{Input: "qj", Output: "↓"},
			{Input: "qg", Output: "…"},
		},
	)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	gotGoogle := EmitGoogle(compiled)
	wantGoogle := "" +
		"q\t\tq\n" +
		"j\th\t\n" +
		"g\tf\t\n" +
		"qj\t↓\t\n" +
		"qg\t…\t\n"

	if gotGoogle != wantGoogle {
		t.Fatalf("EmitGoogle() = %q, want %q", gotGoogle, wantGoogle)
	}
}

func TestCompileRejectsUnknownLayer(t *testing.T) {
	_, err := Compile(
		[]MappingEntry{{Physical: "q", Value: "*x", IsLayer: true, Layer: "x"}},
		nil,
	)
	if err == nil || !strings.Contains(err.Error(), "undefined layer reference") {
		t.Fatalf("Compile() error = %v, want undefined layer reference", err)
	}
}

func TestCompileRejectsConflictingSequence(t *testing.T) {
	_, err := Compile(
		[]MappingEntry{{Physical: "q", Value: "*q", IsLayer: true, Layer: "q"}},
		[]SequenceEntry{{Input: "q", Output: "x"}},
	)
	if err == nil || !strings.Contains(err.Error(), "conflicting rule") {
		t.Fatalf("Compile() error = %v, want conflicting rule", err)
	}
}

func TestEmitYAML(t *testing.T) {
	compiled, err := Compile(
		[]MappingEntry{
			{Physical: "q", Value: "*q", IsLayer: true, Layer: "q"},
			{Physical: "j", Value: "h"},
		},
		[]SequenceEntry{{Input: "qj", Output: "↓"}},
	)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	got := EmitYAML(compiled)
	for _, want := range []string{
		"single_rules:",
		`input: "j"`,
		`output: "h"`,
		"layer_keys:",
		`next: "q"`,
		"sequence_rules:",
		`output: "↓"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("EmitYAML() missing %q in %q", want, got)
		}
	}
}

func TestRunWritesGoogleToFile(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "mappings.txt")
	outputPath := filepath.Join(dir, "out.txt")

	if err := os.WriteFile(inputPath, []byte("q *q\nj h\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(input): %v", err)
	}

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	code := run([]string{
		"-input", inputPath,
		"-format", "google",
		"-output", outputPath,
	}, &stderr, &stdout)

	if code != 0 {
		t.Fatalf("run() code = %d, stderr = %q", code, stderr.String())
	}

	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("os.ReadFile(output): %v", err)
	}

	want := "q\t\tq\nj\th\t\n"
	if string(got) != want {
		t.Fatalf("output file = %q, want %q", string(got), want)
	}

	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}

func writeTempFile(t *testing.T, name, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("os.WriteFile(%s): %v", name, err)
	}
	return path
}
