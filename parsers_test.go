package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestEOLSplit(t *testing.T) {
	type testpair struct {
		input    rune
		expected bool
	}
	var tests = []testpair{
		{
			'\n', true,
		},
		{
			'\r', true,
		},
		{
			'\t', false,
		},
		{
			'\\', false,
		},
		{
			'\f', false,
		},
		{
			'a', false,
		},
		{
			'1', false,
		},
		{
			'世', false,
		},
	}

	for _, pair := range tests {
		actual := EOLSplit(pair.input)
		if actual != pair.expected {
			t.Errorf("Testing EOLSplit with %v. Expected %v but got %v instead", pair.input, pair.expected, actual)
		}
	}
}

func TestParseSRTFile(t *testing.T) {

	type testpair struct {
		input          string
		expected       SubtitleFile
		expectedErrors []error
	}

	var emptySubtitleFile SubtitleFile
	var emptyTimeDuration time.Duration

	shortSRTFile := SubtitleFile{[]Subtitle{
		{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
		{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
		{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
		{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
		{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`, "", ""},
	},
		"",
	}

	sampleWrongTimestamps := SubtitleFile{[]Subtitle{
		{1, emptyTimeDuration, time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
		{2, emptyTimeDuration, emptyTimeDuration, `Έχουμε χάσει αγαπημένους μας.`, "", ""},
		{3, time.Duration(time.Second*10 + time.Millisecond*88), emptyTimeDuration, `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
		{4, emptyTimeDuration, emptyTimeDuration, `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
		{5, emptyTimeDuration, time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`, "", ""},
	}, "",
	}

	sampleWrongIndices := SubtitleFile{[]Subtitle{
		{0, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
		{0, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
		{0, emptyTimeDuration, emptyTimeDuration, `αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
		{0, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
		{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`, "", ""},
	},
		"",
	}

	var tests = []testpair{
		{
			"wrongfilename",
			emptySubtitleFile,
			[]error{errors.New("Something went wrong while trying to parse the provided file!")},
		},
		{
			"samples/sample.srt",
			shortSRTFile,
			nil,
		},
		{
			"samples/sample_short_dos_eol.srt",
			shortSRTFile,
			nil,
		},
		{
			"samples/sample_short_mixed_eol.srt",
			shortSRTFile,
			nil,
		},
		{
			"samples/sample_short_nix_eol.srt",
			shortSRTFile,
			nil,
		},
		{
			"samples/sample_wrong_timestamps.srt",
			sampleWrongTimestamps,
			[]error{
				errors.New("Unexpected parsed seconds value, should be between 0 and 60"),
				errors.New("Wrong Number of fields resulting from input timestamp"),
				errors.New("Unexpected parsed millisecond value, should be between 0 and 999"),
				errors.New("Unexpected parsed seconds value, should be between 0 and 60"),
				errors.New("Unexpected parsed minute value, should be between 0 and 60"),
				errors.New("Unexpected parsed minute value, should be between 0 and 60"),
				errors.New("Wrong Number of fields resulting from input timestamp"),
			},
		},
		{
			"samples/sample_wrong_indices.srt",
			sampleWrongIndices,
			[]error{
				errors.New(`strconv.Atoi: parsing "a": invalid syntax`),
				errors.New(`strconv.Atoi: parsing "2d": invalid syntax`),
				errors.New(`strconv.Atoi: parsing "00:00:10,088 --> 00:00:14,500": invalid syntax`),
				errors.New(`Wrong Number of fields resulting from input timestamp`),
				errors.New(`Wrong Number of fields resulting from input timestamp`),
				errors.New(`strconv.Atoi: parsing "4   ": invalid syntax`),
			},
		},
	}

	for _, pair := range tests {
		actual, actualErrors := ParseSRTFile(pair.input)
		//fmt.Println(cmp.Equal(actual, pair.expected))
		if !cmp.Equal(actual, pair.expected) {
			fmt.Println(actual, actualErrors)
			t.Errorf("Testing ParseSRTFile using %v. Expected %v but got %v instead", pair.input, pair.expected, actual)
		}

		if pair.expectedErrors != nil && !ErrorSlicesEqual(actualErrors, pair.expectedErrors) {
			t.Errorf("Testing ParseSRTFile with %v. Expected errors as %v but got %v instead!", pair.input, pair.expectedErrors, actualErrors)
		}
	}
}

func TestTimestampSplitSRT(t *testing.T) {
	type testpair struct {
		input    rune
		expected bool
	}
	var tests = []testpair{
		{
			',', true,
		},
		{
			'.', true,
		},
		{
			':', true,
		},
		{
			'a', false,
		},
		{
			'1', false,
		},
		{
			'世', false,
		},
	}

	for _, pair := range tests {
		actual := TimestampSplitSRT(pair.input)
		if actual != pair.expected {
			t.Errorf("Testing TimestampSplitSRT with %v. Expected %v but got %v instead", pair.input, pair.expected, actual)
		}
	}
}
