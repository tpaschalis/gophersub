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
		{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
		{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
		{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
		{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
		{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
	}
	_, _ = originalText, parsedSRTFile

	type testpair struct {
		input    SubtitleFile
		expected SubtitleFile
		shift    time.Duration
	}
	var tests = []testpair{
		{
			SubtitleFile{
				[]Subtitle{
					{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
					{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
					{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
					{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
					{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
				},
				"",
			},
			SubtitleFile{[]Subtitle{
				{1, time.Duration(time.Second*3 + time.Millisecond*602), time.Duration(time.Second*5 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
				{2, time.Duration(time.Second*6 + time.Millisecond*536), time.Duration(time.Second*9 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
				{3, time.Duration(time.Second*12 + time.Millisecond*88), time.Duration(time.Second*16 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
				{4, time.Duration(time.Second*16 + time.Millisecond*611), time.Duration(time.Second*18 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
				{5, time.Duration(time.Second*19 + time.Millisecond*929), time.Duration(time.Second*21 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
			},
				"",
			},

			time.Duration(time.Second * 2),
		},
	}
	for _, pair := range tests {
		actual := TimeshiftSubtitleFile(pair.input, pair.shift)
		if len(actual.Subtitles) != len(pair.expected.Subtitles) {
			t.Errorf("The length of the returned SubtitleFile (%v) is not the same as the lenght of the input SubtitleFile (%v) as expected", len(actual.Subtitles), len(pair.input.Subtitles))
		}
		for i, _ := range pair.input.Subtitles {
			if !cmp.Equal(actual.Subtitles[i], pair.expected.Subtitles[i]) {
				t.Errorf("There was an error while timeshifting a test-case subtitle. With input (%v), expected (%v) but got (%v)", pair.input.Subtitles[i], pair.expected.Subtitles[i], actual.Subtitles[i])
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
			SubtitleFile{[]Subtitle{
				{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
				{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
				{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
				{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
				{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
			},
				"",
			},
			SubtitleFile{
				[]Subtitle{
					{1, time.Duration(time.Millisecond * 801), time.Duration(time.Second*1 + time.Millisecond*657), `Έχουμε όλοι υποφέρει.`, "", ""},
					{2, time.Duration(time.Second*2 + time.Millisecond*268), time.Duration(time.Second*3 + time.Millisecond*689 + time.Microsecond*500), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
					{3, time.Duration(time.Second*5 + time.Millisecond*44), time.Duration(time.Second*7 + time.Millisecond*250), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
					{4, time.Duration(time.Second*7 + time.Millisecond*305 + time.Microsecond*500), time.Duration(time.Second*8 + time.Millisecond*284), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
					{5, time.Duration(time.Second*8 + time.Millisecond*964 + time.Microsecond*500), time.Duration(time.Second*9 + time.Millisecond*875 + time.Microsecond*500), `Σας προσφέρω την επιλογή..`, "", ""},
				},
				"",
			},
			2.,
			nil,
		},
		{
			SubtitleFile{
				[]Subtitle{
					{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
					{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
					{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
					{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
					{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
				},
				"",
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

func TestDetectOverlaps(t *testing.T) {
	type testpair struct {
		input    []Subtitle
		expected []Subtitle
	}
	var emptyOverlaps []Subtitle
	shortSRTFile := []Subtitle{
		{1, time.Duration(time.Second*1 + time.Millisecond*500), time.Duration(time.Second*3 + time.Millisecond*300), `one`, "", ""},
		{2, time.Duration(time.Second*4 + time.Millisecond*520), time.Duration(time.Second*7 + time.Millisecond*300), `two`, "", ""},
		{3, time.Duration(time.Second*10 + time.Millisecond*80), time.Duration(time.Second*14 + time.Millisecond*500), `three`, "", ""},
		{4, time.Duration(time.Second*14 + time.Millisecond*600), time.Duration(time.Second*16 + time.Millisecond*200), `four`, "", ""},
		{5, time.Duration(time.Second*17 + time.Millisecond*900), time.Duration(time.Second*19 + time.Millisecond*800), `five`, "", ""},
	}
	overlap1 := []Subtitle{
		{1, time.Duration(time.Second*1 + time.Millisecond*500), time.Duration(time.Second*3 + time.Millisecond*300), `one`, "", ""},
		{2, time.Duration(time.Second*2 + time.Millisecond*520), time.Duration(time.Second*7 + time.Millisecond*300), `two`, "", ""},
		{3, time.Duration(time.Second*10 + time.Millisecond*80), time.Duration(time.Second*14 + time.Millisecond*500), `three`, "", ""},
		{4, time.Duration(time.Second*14 + time.Millisecond*600), time.Duration(time.Second*16 + time.Millisecond*200), `four`, "", ""},
		{5, time.Duration(time.Second*17 + time.Millisecond*900), time.Duration(time.Second*19 + time.Millisecond*800), `five`, "", ""},
	}
	overlap2 := []Subtitle{
		{1, time.Duration(time.Second*1 + time.Millisecond*500), time.Duration(time.Second*3 + time.Millisecond*300), `one`, "", ""},
		{2, time.Duration(time.Second*4 + time.Millisecond*520), time.Duration(time.Second*12 + time.Millisecond*300), `two`, "", ""},
		{3, time.Duration(time.Second*10 + time.Millisecond*80), time.Duration(time.Second*14 + time.Millisecond*500), `three`, "", ""},
		{4, time.Duration(time.Second*14 + time.Millisecond*600), time.Duration(time.Second*16 + time.Millisecond*200), `four`, "", ""},
		{5, time.Duration(time.Second*17 + time.Millisecond*900), time.Duration(time.Second*19 + time.Millisecond*800), `five`, "", ""},
	}
	overlap3 := []Subtitle{
		{1, time.Duration(time.Second*1 + time.Millisecond*500), time.Duration(time.Second*3 + time.Millisecond*300), `one`, "", ""},
		{2, time.Duration(time.Second*4 + time.Millisecond*520), time.Duration(time.Second*7 + time.Millisecond*300), `two`, "", ""},
		{3, time.Duration(time.Second*10 + time.Millisecond*80), time.Duration(time.Second*14 + time.Millisecond*500), `three`, "", ""},
		{4, time.Duration(time.Second*14 + time.Millisecond*600), time.Duration(time.Second*16 + time.Millisecond*200), `four`, "", ""},
		{5, time.Duration(time.Second*6 + time.Millisecond*900), time.Duration(time.Second*19 + time.Millisecond*800), `five`, "", ""},
	}

	var tests = []testpair{
		{
			shortSRTFile,
			emptyOverlaps,
		},
		{
			overlap1,
			[]Subtitle{
				{1, time.Duration(time.Second*1 + time.Millisecond*500), time.Duration(time.Second*3 + time.Millisecond*300), `one`, "", ""},
				{2, time.Duration(time.Second*2 + time.Millisecond*520), time.Duration(time.Second*7 + time.Millisecond*300), `two`, "", ""},
			},
		},
		{
			overlap2,
			[]Subtitle{
				{2, time.Duration(time.Second*4 + time.Millisecond*520), time.Duration(time.Second*12 + time.Millisecond*300), `two`, "", ""},
				{3, time.Duration(time.Second*10 + time.Millisecond*80), time.Duration(time.Second*14 + time.Millisecond*500), `three`, "", ""},
			},
		},
		{
			overlap3,
			[]Subtitle{
				{4, time.Duration(time.Second*14 + time.Millisecond*600), time.Duration(time.Second*16 + time.Millisecond*200), `four`, "", ""},
				{5, time.Duration(time.Second*6 + time.Millisecond*900), time.Duration(time.Second*19 + time.Millisecond*800), `five`, "", ""},
			},
		},
		{
			[]Subtitle{},
			[]Subtitle{},
		},
	}

	for _, pair := range tests {
		actual := DetectOverlaps(SubtitleFile{pair.input, ""})
		if len(actual) == 0 && len(pair.expected) != 0 {
			t.Errorf("Testing DetectOverlaps with empty input %v. Expected %v but got %v instead!", pair.input, pair.expected, actual)
		}
		if len(actual) != 0 && !cmp.Equal(actual, pair.expected) {
			t.Errorf("Testing DetectOverlaps with %v. Expected %v but got %v instead!", pair.input, pair.expected, actual)
		}

	}

}

func TestSerializeSubtitles(t *testing.T) {
	type testpair struct {
		input    SubtitleFile
		expected SubtitleFile
	}

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

	var tests = []testpair{
		{
			SubtitleFile{[]Subtitle{
				{0, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
				{0, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
				{0, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
				{0, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
				{0, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`, "", ""},
			},
				"",
			},
			shortSRTFile,
		},
		{
			SubtitleFile{[]Subtitle{
				{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
				{1, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
				{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
				{5, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
				{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`, "", ""},
			},
				"",
			},
			shortSRTFile,
		},
		{
			SubtitleFile{[]Subtitle{
				{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
				{3, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
				{2, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
				{5, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
				{4, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`, "", ""},
			},
				"",
			},
			shortSRTFile,
		},
		{
			SubtitleFile{[]Subtitle{
				{-1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
				{-2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
				{-3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
				{100, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
				{234325, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`, "", ""},
			},
				"",
			},
			shortSRTFile,
		},
		{
			SubtitleFile{[]Subtitle{}, ""},
			SubtitleFile{[]Subtitle{}, ""},
		},
		//{},
		//{},
	}

	for _, pair := range tests {
		actual := SerializeSubtitles(pair.input)
		if !cmp.Equal(actual, pair.expected) {
			t.Errorf("Testing SerializeSubtitles using %v. Expected %v but got %v instead!", pair.input, pair.expected, actual)
		}
	}
}

func TestRemoveSubtitles(t *testing.T) {
	type testpair struct {
		in1         SubtitleFile
		in2         int
		expected    SubtitleFile
		expectedErr error
	}

	shortSRTFile := SubtitleFile{[]Subtitle{
		{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
		{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
		{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
		{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
		{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
	},
		"",
	}
	var tests = []testpair{
		{
			shortSRTFile,
			-1,
			shortSRTFile,
			errors.New("The index marked for removal is invalid :-1"),
		},
		{
			shortSRTFile,
			-9,
			shortSRTFile,
			errors.New("The index marked for removal is invalid :-9"),
		},
		{
			shortSRTFile,
			0,
			shortSRTFile,
			errors.New("The index marked for removal is invalid :0"),
		},
		{
			shortSRTFile,
			99,
			shortSRTFile,
			errors.New("The index marked for removal is invalid :99"),
		},
		{
			shortSRTFile,
			2,
			SubtitleFile{[]Subtitle{
				{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
				{2, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
				{3, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
				{4, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
			},
				"",
			},
			nil,
		},
		{
			shortSRTFile,
			4,
			SubtitleFile{
				[]Subtitle{
					{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
					{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
					{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
					{4, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
				},
				"",
			},
			nil,
		},
		{
			shortSRTFile,
			1,
			SubtitleFile{
				[]Subtitle{
					{1, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
					{2, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
					{3, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
					{4, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
				},
				"",
			},
			nil,
		},
		{
			shortSRTFile,
			5,
			SubtitleFile{
				[]Subtitle{
					{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
					{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
					{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
					{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
				},
				"",
			},
			nil,
		},
	}

	for _, pair := range tests {
		var actual SubtitleFile
		var actualErr error
		actual, actualErr = RemoveSubtitle(pair.in1, pair.in2)
		if actualErr != nil && pair.expectedErr.Error() != actualErr.Error() {
			t.Errorf("Testing RemoveSubtitle using %v. Expected error %v but got %v instead!", pair.in2, pair.expectedErr, actualErr)
		}
		if !cmp.Equal(actual, pair.expected) {
			t.Errorf("Testing RemoveSubtitle using %v, %v. Expected %v but got %v instead!", pair.in1, pair.in2, pair.expected, actual)
		}
	}
}

func TestAddSubtitle(t *testing.T) {
	type testpair struct {
		in1         SubtitleFile
		in2         string
		in3         string
		in4         string
		expected    SubtitleFile
		expectedErr error
	}

	shortSRTFile := SubtitleFile{[]Subtitle{
		{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
		{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
		{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
		{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
		{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
	},
		"",
	}
	var tests = []testpair{
		{
			shortSRTFile,
			"3.4s",
			"3.9s",
			`PEW`,
			SubtitleFile{
				[]Subtitle{
					{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
					{2, time.Duration(time.Second*3 + time.Millisecond*400), time.Duration(time.Second*3 + time.Millisecond*900), `PEW`, "", ""},
					{3, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
					{4, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
					{5, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
					{6, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
				},
				"",
			},
			nil,
		},
		{
			shortSRTFile,
			"16.57s",
			"17.801s",
			`PEW PEW
PEW`,
			SubtitleFile{
				[]Subtitle{
					{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
					{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
					{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
					{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
					{5, time.Duration(time.Second*16 + time.Millisecond*570), time.Duration(time.Second*17 + time.Millisecond*801), `PEW PEW
PEW`, "", ""},

					{6, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
				},
				"",
			},
			nil,
		},
		{
			shortSRTFile,
			"1.480s",
			"1.790s",
			`test`,
			shortSRTFile,
			errors.New("New subtitle would overlap with existing ones, ignoring it...1.480s - 1.790s"),
		},
		{
			shortSRTFile,
			"19.605s",
			"20.8s",
			`test`,
			shortSRTFile,
			errors.New("New subtitle would overlap with existing ones, ignoring it...19.605s - 20.8s"),
		},
		{
			shortSRTFile,
			"15s",
			"17s",
			`test`,
			shortSRTFile,
			errors.New("New subtitle would overlap with existing ones, ignoring it...15s - 17s"),
		},
		{
			shortSRTFile,
			"-1.2345s",
			"0.01s",
			`test`,
			shortSRTFile,
			errors.New("Start and End times should be positive and ordered, ignoring input... -1.2345s - 0.01s"),
		},
		{
			shortSRTFile,
			"0.6s",
			"0.4s",
			`test`,
			shortSRTFile,
			errors.New("Start and End times should be positive and ordered, ignoring input... 0.6s - 0.4s"),
		},
		{
			shortSRTFile,
			"1.5s",
			"-1.0s",
			`test`,
			shortSRTFile,
			errors.New("Start and End times should be positive and ordered, ignoring input... 1.5s - -1.0s"),
		},
		{
			shortSRTFile,
			"0.004s",
			"0.100s",
			`start of file`,
			SubtitleFile{
				[]Subtitle{
					{1, time.Duration(time.Millisecond * 4), time.Duration(time.Millisecond * 100), `start of file`, "", ""},
					{2, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
					{3, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
					{4, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
					{5, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
					{6, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
				},
				"",
			},
			nil,
		},
		{
			shortSRTFile,
			"100.4s",
			"200.900s",
			`end of file`,
			SubtitleFile{
				[]Subtitle{
					{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
					{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
					{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
					{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
					{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
					{6, time.Duration(time.Second*100 + time.Millisecond*400), time.Duration(time.Second*200 + time.Millisecond*900), `end of file`, "", ""},
				},
				"",
			},
			nil,
		},

		// Test for adding at start and end of file
		// Test for adding overlapping subtitle, and it being skipped
	}

	for _, pair := range tests {
		var actual SubtitleFile
		var actualErr error
		actual, actualErr = AddSubtitle(pair.in1, pair.in2, pair.in3, pair.in4)

		if actualErr != nil && pair.expectedErr.Error() != actualErr.Error() {
			t.Errorf("Testing AddSubtitle using %v. Expected error %v but got %v instead!", pair.in2, pair.expectedErr, actualErr)
		}
		if !cmp.Equal(actual, pair.expected) {
			t.Errorf("Testing AddSubtitle using %v, %v. Expected %v but got %v instead!", pair.in1, pair.in2, pair.expected, actual)
		}
	}
}
