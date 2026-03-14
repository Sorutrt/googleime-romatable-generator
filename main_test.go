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

func TestLoadRelativeLayerSequences(t *testing.T) {
	dir := t.TempDir()
	mappingPath := filepath.Join(dir, "qwerty_to_other")
	layerPath := filepath.Join(dir, "layer_q")

	if err := os.WriteFile(mappingPath, []byte("q *q\nj h\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(mapping): %v", err)
	}

	if err := os.WriteFile(layerPath, []byte("h ←\nj ↓\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(layer): %v", err)
	}

	mappings, err := ParseMappingFile(mappingPath)
	if err != nil {
		t.Fatalf("ParseMappingFile() error = %v", err)
	}

	got, err := LoadRelativeLayerSequences(mappingPath, mappings)
	if err != nil {
		t.Fatalf("LoadRelativeLayerSequences() error = %v", err)
	}

	want := []SequenceEntry{
		{Input: "qh", Output: "←"},
		{Input: "qj", Output: "↓"},
	}
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestGenerateRomajiSequences(t *testing.T) {
	sequences, err := GenerateRomajiSequences([]MappingEntry{
		{Physical: "d", Value: "a"},
		{Physical: "s", Value: "i"},
		{Physical: "e", Value: "u"},
		{Physical: "a", Value: "e"},
		{Physical: "f", Value: "o"},
		{Physical: "h", Value: "k"},
		{Physical: "o", Value: "y"},
	})
	if err != nil {
		t.Fatalf("GenerateRomajiSequences() error = %v", err)
	}

	got := make(map[string]string, len(sequences))
	for _, entry := range sequences {
		got[entry.Input] = entry.Output
	}

	for input, want := range map[string]string{
		"hd":  "か",
		"hs":  "き",
		"he":  "く",
		"ha":  "け",
		"hf":  "こ",
		"hod": "きゃ",
		"hoe": "きゅ",
		"hof": "きょ",
	} {
		if got[input] != want {
			t.Fatalf("generated[%q] = %q, want %q", input, got[input], want)
		}
	}
}

func TestAppendMissingSequencesKeepsExisting(t *testing.T) {
	got := appendMissingSequences(
		[]SequenceEntry{{Input: "hd", Output: "custom"}},
		[]SequenceEntry{
			{Input: "hd", Output: "か"},
			{Input: "hs", Output: "き"},
		},
	)

	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if got[0] != (SequenceEntry{Input: "hd", Output: "custom"}) {
		t.Fatalf("got[0] = %#v, want existing custom sequence", got[0])
	}
	if got[1] != (SequenceEntry{Input: "hs", Output: "き"}) {
		t.Fatalf("got[1] = %#v, want generated hs -> き", got[1])
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

func TestRunLoadsLayerFilesAutomatically(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "qwerty_to_other")
	layerPath := filepath.Join(dir, "layer_q")

	if err := os.WriteFile(inputPath, []byte("q *q\nj h\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(input): %v", err)
	}
	if err := os.WriteFile(layerPath, []byte("h ←\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(layer): %v", err)
	}

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	code := run([]string{
		"-input", inputPath,
		"-format", "google",
	}, &stderr, &stdout)

	if code != 0 {
		t.Fatalf("run() code = %d, stderr = %q", code, stderr.String())
	}

	got := stdout.String()
	for _, want := range []string{
		"q\t\tq\n",
		"j\th\t\n",
		"qh\t←\t\n",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("stdout missing %q in %q", want, got)
		}
	}
}

func TestRunAutoGeneratesRomajiSequences(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "qwerty_to_other")

	if err := os.WriteFile(inputPath, []byte(strings.Join([]string{
		"d a",
		"s i",
		"e u",
		"a e",
		"f o",
		"h k",
		"o y",
	}, "\n")+"\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(input): %v", err)
	}

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	code := run([]string{
		"-input", inputPath,
		"-format", "google",
	}, &stderr, &stdout)

	if code != 0 {
		t.Fatalf("run() code = %d, stderr = %q", code, stderr.String())
	}

	got := stdout.String()
	for _, want := range []string{
		"h\tk\t\n",
		"d\ta\t\n",
		"hd\tか\t\n",
		"hod\tきゃ\t\n",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("stdout missing %q in %q", want, got)
		}
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
