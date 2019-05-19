package main

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

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
			SubtitleFile{
				[]Subtitle{
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
		actual, actualErr = AddSubtitle(pair.in1, pair.in2, pair.in3, pair.in4, "", "")

		if actualErr != nil && pair.expectedErr.Error() != actualErr.Error() {
			t.Errorf("Testing AddSubtitle using %v. Expected error %v but got %v instead!", pair.in2, pair.expectedErr, actualErr)
		}
		if !cmp.Equal(actual, pair.expected) {
			t.Errorf("Testing AddSubtitle using %v, %v. Expected %v but got %v instead!", pair.in1, pair.in2, pair.expected, actual)
		}
	}
}

func TestPrintSubfileInfo(t *testing.T) {

	in := SubtitleFile{
		[]Subtitle{
			{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `one`, "", ""},
			{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `two`, "", ""},
			{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `three.`, "", ""},
			{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `four`, "", ""},
			{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `five`, "", ""},
		},
		"sample_headers",
	}

	PrintSubfileInfo(in)
	//Headers : sample_headers
	//Number of subtitles : 5
	//Start Time : 1.602s
	//End Time : 19.751s
	//First-to-last Runtime : 18.149s
	//Subtitle Runtime : 12s
	//
	//An average human reads at a pace of about 850 Characters Per Minute (CPM)
	//Highest CPM : 131.72 on subtitle index : 5
	//Lowest CPM : 63.31 on subtitle index : 2
	//Average CPM : 39.57
}

func TestSearchSubtitleFile(t *testing.T) {
	type testpair struct {
		input       SubtitleFile
		searchterm  string
		expected    []Subtitle
		expectedErr error
	}

	var emptySubtitleSlice []Subtitle

	shortSRTFile := SubtitleFile{
		[]Subtitle{
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
			`Έχουμε`,
			[]Subtitle{
				{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
				{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
			},
			nil,
		},
		{
			shortSRTFile,
			`έ|ύ`,
			[]Subtitle{
				{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`, "", ""},
				{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`, "", ""},
				{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`, "", ""},
				{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`, "", ""},
				{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`, "", ""},
			},
			nil,
		},
		{
			shortSRTFile,
			`[`,
			emptySubtitleSlice,
			errors.New("The provided search term is invalid :`[`"),
		},
	}

	for _, pair := range tests {
		actual, actualErr := SearchSubtitleFile(pair.input, pair.searchterm)

		if actualErr != nil && pair.expectedErr.Error() != actualErr.Error() {
			t.Errorf("Testing SearchSubtitleFile using %v. Expected error %v but got %v instead!", pair.input, pair.searchterm, actualErr)
		}
		if !cmp.Equal(actual, pair.expected) {
			t.Errorf("Testing SearchSubtitleFile using %v, %v. Expected %v but got %v instead!", pair.input, pair.searchterm, pair.expected, actual)
		}
	}
}
