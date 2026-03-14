package main

import "fmt"

func Compile(mappings []MappingEntry, sequences []SequenceEntry) (Compiled, error) {
	knownPhysical := make(map[string]struct{}, len(mappings))
	seenPhysical := make(map[string]struct{}, len(mappings))

	for _, entry := range mappings {
		if _, exists := seenPhysical[entry.Physical]; exists {
			return Compiled{}, fmt.Errorf("duplicate physical key %q", entry.Physical)
		}
		seenPhysical[entry.Physical] = struct{}{}
		knownPhysical[entry.Physical] = struct{}{}
	}

	var compiled Compiled
	seenInputs := make(map[string]struct{}, len(mappings)+len(sequences))

	for _, entry := range mappings {
		if entry.IsLayer {
			if _, exists := knownPhysical[entry.Layer]; !exists {
				return Compiled{}, fmt.Errorf("undefined layer reference %q for %q", entry.Layer, entry.Physical)
			}

			if err := reserveInput(seenInputs, entry.Physical); err != nil {
				return Compiled{}, err
			}

			layer := LayerKey{
				Input: entry.Physical,
				Next:  entry.Layer,
			}
			compiled.LayerKeys = append(compiled.LayerKeys, layer)
			compiled.GoogleRows = append(compiled.GoogleRows, GoogleRow{
				Input: entry.Physical,
				Next:  entry.Layer,
			})
			continue
		}

		if err := reserveInput(seenInputs, entry.Physical); err != nil {
			return Compiled{}, err
		}

		rule := SingleRule{
			Input:  entry.Physical,
			Output: entry.Value,
		}
		compiled.SingleRules = append(compiled.SingleRules, rule)
		compiled.GoogleRows = append(compiled.GoogleRows, GoogleRow{
			Input:  entry.Physical,
			Output: entry.Value,
		})
	}

	seenSequences := make(map[string]struct{}, len(sequences))
	for _, entry := range sequences {
		if _, exists := seenSequences[entry.Input]; exists {
			return Compiled{}, fmt.Errorf("duplicate sequence %q", entry.Input)
		}
		seenSequences[entry.Input] = struct{}{}

		if err := reserveInput(seenInputs, entry.Input); err != nil {
			return Compiled{}, err
		}

		rule := SequenceRule{
			Input:  entry.Input,
			Output: entry.Output,
		}
		compiled.SequenceRules = append(compiled.SequenceRules, rule)
		compiled.GoogleRows = append(compiled.GoogleRows, GoogleRow{
			Input:  entry.Input,
			Output: entry.Output,
		})
	}

	return compiled, nil
}

func reserveInput(seen map[string]struct{}, input string) error {
	if _, exists := seen[input]; exists {
		return fmt.Errorf("conflicting rule for input %q", input)
	}
	seen[input] = struct{}{}
	return nil
}
