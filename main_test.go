package main

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDurationToTimestamp(t *testing.T) {
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
		actual := DurationToTimestamp(pair.input)
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

func TestTimeshiftSRTFile(t *testing.T) {
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
		actual := TimeshiftSRTFile(pair.input, pair.shift)
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
