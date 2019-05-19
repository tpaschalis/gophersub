package main

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"
)

func TestToSRTFile(t *testing.T) {
	type testpair struct {
		inputSubfile SubtitleFile
		inputFn      string
		expectedFn   string
		expectedErr  error
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
					{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`, "", ""},
				},
				"",
			},
			"samples/exportFile-01-tmp.srt",
			"samples/exportFile-01.srt",
			nil,
		},
	}

	for _, pair := range tests {
		actualErr := ToSRTFile(pair.inputSubfile, pair.inputFn)
		if actualErr != nil && pair.expectedErr.Error() != actualErr.Error() {
			t.Errorf("Testing ToSRTFile using %v. Expected error %v but got %v instead!", pair.inputFn, pair.expectedErr, actualErr)
		}

		f1, err := ioutil.ReadFile(pair.expectedFn)
		if err != nil {
			t.Errorf("Testing ToSRTFile.\nCould not open file %v for comparing expected and actual results", pair.expectedFn)
		}
		f2, err := ioutil.ReadFile(pair.inputFn)
		if err != nil {
			t.Errorf("Testing ToSRTFile.\nCould not open file %v for comparing expected and actual results", pair.inputFn)
		}
		if !bytes.Equal(f1, f2) {
			t.Errorf("Testing ToSRTFIle.\nMismatch between %v and %v.", pair.inputFn, pair.expectedFn)
		}
	}
}
