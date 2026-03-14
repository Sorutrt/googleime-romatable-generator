package main

import (
	"fmt"
	"regexp"
)

var singleRomajiPattern = regexp.MustCompile(`^[a-z]$`)

type romajiKanaEntry struct {
	Romaji string
	Kana   string
}

var defaultRomajiKanaTable = []romajiKanaEntry{
	{Romaji: "ka", Kana: "か"},
	{Romaji: "ki", Kana: "き"},
	{Romaji: "ku", Kana: "く"},
	{Romaji: "ke", Kana: "け"},
	{Romaji: "ko", Kana: "こ"},
	{Romaji: "sa", Kana: "さ"},
	{Romaji: "shi", Kana: "し"},
	{Romaji: "si", Kana: "し"},
	{Romaji: "su", Kana: "す"},
	{Romaji: "se", Kana: "せ"},
	{Romaji: "so", Kana: "そ"},
	{Romaji: "ta", Kana: "た"},
	{Romaji: "chi", Kana: "ち"},
	{Romaji: "ti", Kana: "ち"},
	{Romaji: "tsu", Kana: "つ"},
	{Romaji: "tu", Kana: "つ"},
	{Romaji: "te", Kana: "て"},
	{Romaji: "to", Kana: "と"},
	{Romaji: "na", Kana: "な"},
	{Romaji: "ni", Kana: "に"},
	{Romaji: "nu", Kana: "ぬ"},
	{Romaji: "ne", Kana: "ね"},
	{Romaji: "no", Kana: "の"},
	{Romaji: "ha", Kana: "は"},
	{Romaji: "hi", Kana: "ひ"},
	{Romaji: "fu", Kana: "ふ"},
	{Romaji: "hu", Kana: "ふ"},
	{Romaji: "he", Kana: "へ"},
	{Romaji: "ho", Kana: "ほ"},
	{Romaji: "ma", Kana: "ま"},
	{Romaji: "mi", Kana: "み"},
	{Romaji: "mu", Kana: "む"},
	{Romaji: "me", Kana: "め"},
	{Romaji: "mo", Kana: "も"},
	{Romaji: "ya", Kana: "や"},
	{Romaji: "yu", Kana: "ゆ"},
	{Romaji: "yo", Kana: "よ"},
	{Romaji: "ra", Kana: "ら"},
	{Romaji: "ri", Kana: "り"},
	{Romaji: "ru", Kana: "る"},
	{Romaji: "re", Kana: "れ"},
	{Romaji: "ro", Kana: "ろ"},
	{Romaji: "wa", Kana: "わ"},
	{Romaji: "wo", Kana: "を"},
	{Romaji: "nn", Kana: "ん"},
	{Romaji: "ga", Kana: "が"},
	{Romaji: "gi", Kana: "ぎ"},
	{Romaji: "gu", Kana: "ぐ"},
	{Romaji: "ge", Kana: "げ"},
	{Romaji: "go", Kana: "ご"},
	{Romaji: "za", Kana: "ざ"},
	{Romaji: "ji", Kana: "じ"},
	{Romaji: "zi", Kana: "じ"},
	{Romaji: "zu", Kana: "ず"},
	{Romaji: "ze", Kana: "ぜ"},
	{Romaji: "zo", Kana: "ぞ"},
	{Romaji: "da", Kana: "だ"},
	{Romaji: "di", Kana: "ぢ"},
	{Romaji: "du", Kana: "づ"},
	{Romaji: "de", Kana: "で"},
	{Romaji: "do", Kana: "ど"},
	{Romaji: "ba", Kana: "ば"},
	{Romaji: "bi", Kana: "び"},
	{Romaji: "bu", Kana: "ぶ"},
	{Romaji: "be", Kana: "べ"},
	{Romaji: "bo", Kana: "ぼ"},
	{Romaji: "pa", Kana: "ぱ"},
	{Romaji: "pi", Kana: "ぴ"},
	{Romaji: "pu", Kana: "ぷ"},
	{Romaji: "pe", Kana: "ぺ"},
	{Romaji: "po", Kana: "ぽ"},
	{Romaji: "kya", Kana: "きゃ"},
	{Romaji: "kyu", Kana: "きゅ"},
	{Romaji: "kyo", Kana: "きょ"},
	{Romaji: "sha", Kana: "しゃ"},
	{Romaji: "shu", Kana: "しゅ"},
	{Romaji: "sho", Kana: "しょ"},
	{Romaji: "sya", Kana: "しゃ"},
	{Romaji: "syu", Kana: "しゅ"},
	{Romaji: "syo", Kana: "しょ"},
	{Romaji: "cha", Kana: "ちゃ"},
	{Romaji: "chu", Kana: "ちゅ"},
	{Romaji: "cho", Kana: "ちょ"},
	{Romaji: "tya", Kana: "ちゃ"},
	{Romaji: "tyu", Kana: "ちゅ"},
	{Romaji: "tyo", Kana: "ちょ"},
	{Romaji: "nya", Kana: "にゃ"},
	{Romaji: "nyu", Kana: "にゅ"},
	{Romaji: "nyo", Kana: "にょ"},
	{Romaji: "hya", Kana: "ひゃ"},
	{Romaji: "hyu", Kana: "ひゅ"},
	{Romaji: "hyo", Kana: "ひょ"},
	{Romaji: "mya", Kana: "みゃ"},
	{Romaji: "myu", Kana: "みゅ"},
	{Romaji: "myo", Kana: "みょ"},
	{Romaji: "rya", Kana: "りゃ"},
	{Romaji: "ryu", Kana: "りゅ"},
	{Romaji: "ryo", Kana: "りょ"},
	{Romaji: "gya", Kana: "ぎゃ"},
	{Romaji: "gyu", Kana: "ぎゅ"},
	{Romaji: "gyo", Kana: "ぎょ"},
	{Romaji: "ja", Kana: "じゃ"},
	{Romaji: "ju", Kana: "じゅ"},
	{Romaji: "jo", Kana: "じょ"},
	{Romaji: "jya", Kana: "じゃ"},
	{Romaji: "jyu", Kana: "じゅ"},
	{Romaji: "jyo", Kana: "じょ"},
	{Romaji: "zya", Kana: "じゃ"},
	{Romaji: "zyu", Kana: "じゅ"},
	{Romaji: "zyo", Kana: "じょ"},
	{Romaji: "bya", Kana: "びゃ"},
	{Romaji: "byu", Kana: "びゅ"},
	{Romaji: "byo", Kana: "びょ"},
	{Romaji: "pya", Kana: "ぴゃ"},
	{Romaji: "pyu", Kana: "ぴゅ"},
	{Romaji: "pyo", Kana: "ぴょ"},
}

var requiredRomajiLetters = buildRequiredRomajiLetters(defaultRomajiKanaTable)

func GenerateRomajiSequences(mappings []MappingEntry) ([]SequenceEntry, error) {
	romajiToPhysical := make(map[string]string)
	for _, entry := range mappings {
		if entry.IsLayer || !singleRomajiPattern.MatchString(entry.Value) {
			continue
		}
		if _, needed := requiredRomajiLetters[entry.Value]; !needed {
			continue
		}

		if physical, exists := romajiToPhysical[entry.Value]; exists {
			return nil, fmt.Errorf("duplicate romaji key %q for %q and %q", entry.Value, physical, entry.Physical)
		}
		romajiToPhysical[entry.Value] = entry.Physical
	}

	sequences := make([]SequenceEntry, 0, len(defaultRomajiKanaTable))
	seenInputs := make(map[string]struct{}, len(defaultRomajiKanaTable))
	for _, entry := range defaultRomajiKanaTable {
		physicalInput, ok := translateRomajiToPhysical(entry.Romaji, romajiToPhysical)
		if !ok {
			continue
		}

		if _, exists := seenInputs[physicalInput]; exists {
			continue
		}
		seenInputs[physicalInput] = struct{}{}
		sequences = append(sequences, SequenceEntry{
			Input:  physicalInput,
			Output: entry.Kana,
		})
	}

	return sequences, nil
}

func appendMissingSequences(existing []SequenceEntry, generated []SequenceEntry) []SequenceEntry {
	seenInputs := make(map[string]struct{}, len(existing)+len(generated))
	for _, entry := range existing {
		seenInputs[entry.Input] = struct{}{}
	}

	merged := append([]SequenceEntry{}, existing...)
	for _, entry := range generated {
		if _, exists := seenInputs[entry.Input]; exists {
			continue
		}
		seenInputs[entry.Input] = struct{}{}
		merged = append(merged, entry)
	}

	return merged
}

func buildRequiredRomajiLetters(entries []romajiKanaEntry) map[string]struct{} {
	letters := make(map[string]struct{})
	for _, entry := range entries {
		for _, r := range entry.Romaji {
			letters[string(r)] = struct{}{}
		}
	}
	return letters
}

func translateRomajiToPhysical(romaji string, romajiToPhysical map[string]string) (string, bool) {
	physical := make([]byte, 0, len(romaji))
	for _, r := range romaji {
		key, ok := romajiToPhysical[string(r)]
		if !ok {
			return "", false
		}
		physical = append(physical, key...)
	}
	return string(physical), true
}
