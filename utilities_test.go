package main

import (
	"errors"
	"testing"
	"time"
)

func TestDurationToTimestampSRT(t *testing.T) {
	type testpair struct {
		input    time.Duration
		expected string
	}
	var tests = []testpair{
		{time.Duration(time.Hour*2 + time.Minute*10 + time.Second*20 + time.Millisecond*183), "02:10:20,183"},
		{time.Duration(time.Hour*2 + time.Minute*80 + time.Second*90 + time.Millisecond*1800), "03:21:31,800"},
		{time.Duration(time.Hour*1 + time.Minute*0 + time.Second*10 + time.Millisecond*0), "01:00:10,000"},
		{time.Duration(time.Hour*20 + time.Minute*5 + time.Second*0 + time.Millisecond*70), "20:05:00,070"},
		{time.Duration(time.Hour*30 + time.Minute*6 + time.Second*0 + time.Millisecond*8), "30:06:00,008"},
		{time.Duration(time.Hour*0 + time.Minute*0 + time.Second*0 + time.Millisecond*8), "00:00:00,008"},
		{time.Duration(time.Hour*0 + time.Minute*0 + time.Second*0 + time.Millisecond*181), "00:00:00,181"},
		{time.Duration(time.Hour*0 + time.Minute*0 + time.Second*3 + time.Millisecond*977), "00:00:03,977"},
		{time.Duration(time.Hour*0 + time.Minute*6 + time.Second*3 + time.Millisecond*977), "00:06:03,977"},
		{time.Duration(time.Hour*0 + time.Minute*7 + time.Second*0 + time.Millisecond*500), "00:07:00,500"},
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
