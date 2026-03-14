package main

type MappingEntry struct {
	Physical string
	Value    string
	IsLayer  bool
	Layer    string
}

type SequenceEntry struct {
	Input  string
	Output string
}

type RelativeSequenceEntry struct {
	Input  string
	Output string
}

type SingleRule struct {
	Input  string
	Output string
}

type LayerKey struct {
	Input string
	Next  string
}

type SequenceRule struct {
	Input  string
	Output string
}

type GoogleRow struct {
	Input  string
	Output string
	Next   string
}

type Compiled struct {
	SingleRules   []SingleRule
	LayerKeys     []LayerKey
	SequenceRules []SequenceRule
	GoogleRows    []GoogleRow
}
