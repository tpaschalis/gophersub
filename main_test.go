package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDurationToTimestampSRT(t *testing.T) {
	type testpair struct {
		input    time.Duration
		expected string
	}
	var tests = []testpair{
		{time.Duration(time.Hour*2 + time.Minute*10 + time.Second*20 + time.Millisecond*183), "02:10:20.183"},
		{time.Duration(time.Hour*2 + time.Minute*80 + time.Second*90 + time.Millisecond*1800), "03:21:31.800"},
		{time.Duration(time.Hour*1 + time.Minute*0 + time.Second*10 + time.Millisecond*0), "01:00:10.000"},
		{time.Duration(time.Hour*20 + time.Minute*5 + time.Second*0 + time.Millisecond*70), "20:05:0.070"},
		{time.Duration(time.Hour*30 + time.Minute*6 + time.Second*0 + time.Millisecond*8), "30:06:0.008"},
		{time.Duration(time.Hour*0 + time.Minute*0 + time.Second*0 + time.Millisecond*8), "00:00:0.008"},
		{time.Duration(time.Hour*0 + time.Minute*0 + time.Second*0 + time.Millisecond*181), "00:00:0.181"},
		{time.Duration(time.Hour*0 + time.Minute*0 + time.Second*3 + time.Millisecond*977), "00:00:3.977"},
		{time.Duration(time.Hour*0 + time.Minute*6 + time.Second*3 + time.Millisecond*977), "00:06:3.977"},
		{time.Duration(time.Hour*0 + time.Minute*7 + time.Second*0 + time.Millisecond*500), "00:07:0.500"},
	}

	for _, pair := range tests {
		actual := DurationToTimestampSRT(pair.input)
		//fmt.Println(pair.input, pair.expected, actual)
		if actual != pair.expected {
			t.Errorf("Expected the duration-to-timestamp conversion to produce \n\n%v from \n\n%v but instead got \n\n%v!", pair.expected, pair.input, actual)
		}
	}

}

func TestStrToDuration(t *testing.T) {
	type testpair struct {
		input       string
		expectedDur time.Duration
		expectedErr error
	}
	var emptyTimeDuration time.Duration
	var tests = []testpair{
		{"1.4s", time.Duration(1*time.Second + 400*time.Millisecond), nil},
		{"5s", time.Duration(5 * time.Second), nil},
		{"0s", emptyTimeDuration, nil},
		{"0.00000s", emptyTimeDuration, nil},
		{"s", emptyTimeDuration, errors.New("time: invalid duration s")},
		{"-1.2121", emptyTimeDuration, errors.New("time: missing unit in duration -1.2121")},
		{"1m15.29s", time.Duration(time.Minute*1 + time.Second*15 + time.Millisecond*290), nil},
		{"-1h5.401s", time.Duration(time.Hour*-1 + time.Second*-5 + time.Millisecond*-401), nil},
	}

	for _, pair := range tests {
		actual, err := StrToDuration(pair.input)
		if actual != pair.expectedDur {
			t.Errorf("Testing StrToDuration with %v. Expected time.Duration as %v but got %v", pair.input, pair.expectedDur, actual)
		}
		if pair.expectedErr != nil && pair.expectedErr.Error() != err.Error() {
			t.Errorf("Testing StrToDuration with %v. Expected errors as %v but got %v instead!", pair.input, pair.expectedErr, err)
		}
	}
}

func TestTimeshiftSubtitleFile(t *testing.T) {
	originalText := `1
00:00:01,602 --> 00:00:03,314
Έχουμε όλοι υποφέρει.

2
00:00:04,536 --> 00:00:07,379
Έχουμε χάσει αγαπημένους μας.

3
00:00:10,088 --> 00:00:14,500
Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.

4
00:00:14,611 --> 00:00:16,568
Κι εγώ σκοπεύω να ζήσω.

5
00:00:17,929 --> 00:00:19,751
Σας προσφέρω την επιλογή...
`
	parsedSRTFile := []Subtitle{
		{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`},
		{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`},
		{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`},
		{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`},
		{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`},
	}
	_, _ = originalText, parsedSRTFile

	type testpair struct {
		input    SubtitleFile
		expected SubtitleFile
		shift    time.Duration
	}
	var tests = []testpair{
		{
			[]Subtitle{
				{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`},
				{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`},
				{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`},
				{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`},
				{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`},
			},
			[]Subtitle{
				{1, time.Duration(time.Second*3 + time.Millisecond*602), time.Duration(time.Second*5 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`},
				{2, time.Duration(time.Second*6 + time.Millisecond*536), time.Duration(time.Second*9 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`},
				{3, time.Duration(time.Second*12 + time.Millisecond*88), time.Duration(time.Second*16 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`},
				{4, time.Duration(time.Second*16 + time.Millisecond*611), time.Duration(time.Second*18 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`},
				{5, time.Duration(time.Second*19 + time.Millisecond*929), time.Duration(time.Second*21 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`},
			},
			time.Duration(time.Second * 2),
		},
	}
	for _, pair := range tests {
		actual := TimeshiftSubtitleFile(pair.input, pair.shift)
		if len(actual) != len(pair.expected) {
			t.Errorf("The length of the returned SubtitleFile (%v) is not the same as the lenght of the input SubtitleFile (%v) as expected", len(actual), len(pair.input))
		}
		for i, _ := range pair.input {
			if !cmp.Equal(actual[i], pair.expected[i]) {
				t.Errorf("There was an error while timeshifting a test-case subtitle. With input (%v), expected (%v) but got (%v)", pair.input[i], pair.expected[i], actual[i])
			}
		}
	}
}

func TestPaceSubtitleFile(t *testing.T) {
	type testpair struct {
		input       SubtitleFile
		expected    SubtitleFile
		rate        float64
		expectedErr error
	}

	var emptySubtitleFile SubtitleFile
	var tests = []testpair{
		{
			[]Subtitle{
				{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`},
				{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`},
				{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`},
				{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`},
				{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`},
			},
			[]Subtitle{
				{1, time.Duration(time.Millisecond * 801), time.Duration(time.Second*1 + time.Millisecond*657), `Έχουμε όλοι υποφέρει.`},
				{2, time.Duration(time.Second*2 + time.Millisecond*268), time.Duration(time.Second*3 + time.Millisecond*689 + time.Microsecond*500), `Έχουμε χάσει αγαπημένους μας.`},
				{3, time.Duration(time.Second*5 + time.Millisecond*44), time.Duration(time.Second*7 + time.Millisecond*250), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`},
				{4, time.Duration(time.Second*7 + time.Millisecond*305 + time.Microsecond*500), time.Duration(time.Second*8 + time.Millisecond*284), `Κι εγώ σκοπεύω να ζήσω.`},
				{5, time.Duration(time.Second*8 + time.Millisecond*964 + time.Microsecond*500), time.Duration(time.Second*9 + time.Millisecond*875 + time.Microsecond*500), `Σας προσφέρω την επιλογή..`},
			},
			2.,
			nil,
		},
		{
			[]Subtitle{
				{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`},
				{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`},
				{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`},
				{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`},
				{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`},
			},
			emptySubtitleFile,
			-1.2,
			errors.New("Input rate should be a positive, floating-point number"),
		},
	}

	for _, pair := range tests {
		actual, err := PaceSubtitleFile(pair.input, pair.rate)
		if !cmp.Equal(actual, pair.expected) {
			t.Errorf("Testing PaceSubtitleFile with %v. Expected SubtitleFile as %v but got %v instead", pair.input, pair.expected, actual)
		}
		if pair.expectedErr != nil && pair.expectedErr.Error() != err.Error() {
			t.Errorf("Testing PaceSubtitleFile with %v. Expected errors as %v but got %v instead!", pair.input, pair.expectedErr, err)
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

	shortSRTFile := SubtitleFile{
		{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`},
		{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`},
		{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`},
		{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`},
		{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`},
	}

	sampleWrongTimestamps := SubtitleFile{
		{1, emptyTimeDuration, time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`},
		{2, emptyTimeDuration, emptyTimeDuration, `Έχουμε χάσει αγαπημένους μας.`},
		{3, time.Duration(time.Second*10 + time.Millisecond*88), emptyTimeDuration, `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`},
		{4, emptyTimeDuration, emptyTimeDuration, `Κι εγώ σκοπεύω να ζήσω.`},
		{5, emptyTimeDuration, time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`},
	}

	sampleWrongIndices := SubtitleFile{
		{0, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`},
		{0, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`},
		{0, emptyTimeDuration, emptyTimeDuration, `αλλά τους ζωντανούς και τους νεκρούς.`},
		{0, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`},
		{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`},
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
		fmt.Println(actual, actualErrors)
		fmt.Println(cmp.Equal(actual, pair.expected))
		if !cmp.Equal(actual, pair.expected) {
			t.Errorf("Testing ParseSRTFile using %v. Expected %v but got %v instead", pair.input, pair.expected, actual)
		}

		if pair.expectedErrors != nil && !ErrorSlicesEqual(actualErrors, pair.expectedErrors) {
			t.Errorf("Testing ParseSRTFile with %v. Expected errors as %v but got %v instead!", pair.input, pair.expectedErrors, actualErrors)
		}
	}
}

func TestErrorSlicesEqual(t *testing.T) {
	type testpair struct {
		in1      []error
		in2      []error
		expected bool
	}

	var tests = []testpair{
		{
			[]error{},
			[]error{},
			true,
		},
		{
			[]error{errors.New("Error1"), errors.New("Error2")},
			[]error{errors.New("Error1"), errors.New("Error2"), errors.New("Error3")},
			false,
		},
		{
			[]error{errors.New("Error1"), errors.New("Error2")},
			[]error{errors.New("Error1"), errors.New("Error2")},
			true,
		},
		{
			[]error{errors.New("Error0"), errors.New("Error ONE")},
			[]error{errors.New("Error0"), errors.New("Error TWO")},
			false,
		},
	}

	for _, pair := range tests {
		actual := ErrorSlicesEqual(pair.in1, pair.in2)
		if actual != pair.expected {
			t.Errorf("Testing ErrorSlicesEqual with %v, %v. Expected errors as %v but got %v instead!", pair.in1, pair.in2, pair.expected, actual)
		}
	}
}
